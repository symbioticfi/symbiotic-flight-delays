package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/spf13/cobra"

	"sum/internal/flights"
)

type config struct {
	listenAddr string
}

var cfg config

var rootCmd = &cobra.Command{
	Use:           "flights-api",
	Short:         "Mock flights API backing the flight delay protocol",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer cancel()

		store := flights.NewStore(seedAirlines(), seedFlights())
		generator := newFlightGenerator(store)
		generator.start(ctx)
		srv := newFlightServer(store)

		httpServer := &http.Server{
			Addr:              cfg.listenAddr,
			Handler:           srv.routes(),
			ReadHeaderTimeout: 5 * time.Second,
			ReadTimeout:       10 * time.Second,
			WriteTimeout:      10 * time.Second,
		}

		done := make(chan struct{})
		go func() {
			<-ctx.Done()
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := httpServer.Shutdown(shutdownCtx); err != nil {
				slog.Error("Failed to shut down flights API", "error", err)
			}
			close(done)
		}()

		slog.Info("Flights API listening", "addr", cfg.listenAddr)
		err := httpServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		<-done
		return nil
	},
}

func main() {
	rootCmd.PersistentFlags().StringVar(&cfg.listenAddr, "listen", ":8085", "HTTP listen address")

	if err := rootCmd.Execute(); err != nil {
		if !errors.Is(err, context.Canceled) {
			slog.Error("flights-api exited with error", "error", err)
		}
		os.Exit(1)
	}
}

type flightServer struct {
	store *flights.Store
}

func newFlightServer(store *flights.Store) *flightServer {
	return &flightServer{store: store}
}

func (s *flightServer) routes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(corsMiddleware)
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
	r.Get("/airlines", s.handleListAirlines)
	r.Post("/airlines", s.handleCreateAirline)
	r.Get("/airlines/{airlineId}/flights", s.handleListFlights)
	r.Post("/airlines/{airlineId}/flights", s.handleCreateFlight)
	r.Post("/airlines/{airlineId}/flights/{flightId}/delay", s.handleUpdateStatus(flights.StatusDelayed))
	r.Post("/airlines/{airlineId}/flights/{flightId}/depart", s.handleUpdateStatus(flights.StatusDeparted))
	return r
}

func (s *flightServer) handleListAirlines(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"airlines": s.store.ListAirlines()})
}

type createAirlineRequest struct {
	AirlineID string `json:"airlineId"`
	Name      string `json:"name"`
	Code      string `json:"code"`
}

func (s *flightServer) handleCreateAirline(w http.ResponseWriter, r *http.Request) {
	var body createAirlineRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if body.AirlineID == "" || body.Name == "" {
		respondError(w, http.StatusBadRequest, "airlineId and name are required")
		return
	}
	if err := s.store.AddAirline(flights.Airline{AirlineID: body.AirlineID, Name: body.Name, Code: body.Code}); err != nil {
		status, msg := mapStoreError(err)
		respondError(w, status, msg)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"airline": flights.Airline{AirlineID: body.AirlineID, Name: body.Name, Code: body.Code}})
}

func (s *flightServer) handleListFlights(w http.ResponseWriter, r *http.Request) {
	airlineID := chi.URLParam(r, "airlineId")
	flightsList, err := s.store.ListFlights(airlineID)
	if err != nil {
		status, msg := mapStoreError(err)
		respondError(w, status, msg)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"flights": flightsList})
}

type createFlightRequest struct {
	FlightID           string `json:"flightId"`
	DepartureTimestamp int64  `json:"departureTimestamp"`
}

func (s *flightServer) handleCreateFlight(w http.ResponseWriter, r *http.Request) {
	airlineID := chi.URLParam(r, "airlineId")
	var body createFlightRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if body.FlightID == "" || body.DepartureTimestamp <= 0 {
		respondError(w, http.StatusBadRequest, "flightId and departureTimestamp are required")
		return
	}
	flight := flights.Flight{AirlineID: airlineID, FlightID: body.FlightID, DepartureTimestamp: body.DepartureTimestamp, Status: flights.StatusScheduled}
	created, err := s.store.CreateFlight(airlineID, flight)
	if err != nil {
		status, msg := mapStoreError(err)
		respondError(w, status, msg)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"flight": created})
}

func (s *flightServer) handleUpdateStatus(status flights.Status) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		airlineID := chi.URLParam(r, "airlineId")
		flightID := chi.URLParam(r, "flightId")
		updated, err := s.store.UpdateStatus(airlineID, flightID, status)
		if err != nil {
			statusCode, msg := mapStoreError(err)
			respondError(w, statusCode, msg)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"flight": updated})
	}
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if payload == nil {
		return
	}
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		slog.Error("failed to encode response", "error", err)
	}
}

func respondError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]any{"error": message})
}

func mapStoreError(err error) (int, string) {
	switch {
	case errors.Is(err, flights.ErrAirlineNotFound):
		return http.StatusNotFound, err.Error()
	case errors.Is(err, flights.ErrInvalidAirline), errors.Is(err, flights.ErrInvalidFlight):
		return http.StatusBadRequest, err.Error()
	case errors.Is(err, flights.ErrAirlineExists):
		return http.StatusConflict, err.Error()
	case errors.Is(err, flights.ErrFlightExists):
		return http.StatusConflict, err.Error()
	case errors.Is(err, flights.ErrFlightNotFound):
		return http.StatusNotFound, err.Error()
	case errors.Is(err, flights.ErrInvalidStatus), errors.Is(err, flights.ErrInvalidStatusTransition):
		return http.StatusBadRequest, err.Error()
	default:
		return http.StatusInternalServerError, "internal server error"
	}
}

func seedAirlines() []flights.Airline {
	return []flights.Airline{
		{AirlineID: "ALPHA", Name: "Alpha Air", Code: "AA"},
		{AirlineID: "BETA", Name: "Beta Wings", Code: "BW"},
		{AirlineID: "GAMMA", Name: "Gamma Connect", Code: "GC"},
	}
}

func seedFlights() []flights.Flight {
	now := time.Now().Unix()
	toSeconds := func(d time.Duration) int64 { return int64(d.Seconds()) }
	return []flights.Flight{
		{AirlineID: "ALPHA", FlightID: "ALPHA-001", DepartureTimestamp: now + toSeconds(1*time.Minute), Status: flights.StatusScheduled},
		{AirlineID: "ALPHA", FlightID: "ALPHA-002", DepartureTimestamp: now + toSeconds(3*time.Minute), Status: flights.StatusScheduled},
		{AirlineID: "BETA", FlightID: "BETA-451", DepartureTimestamp: now + toSeconds(2*time.Minute), Status: flights.StatusScheduled},
		{AirlineID: "GAMMA", FlightID: "GAMMA-882", DepartureTimestamp: now + toSeconds(4*time.Minute), Status: flights.StatusScheduled},
	}
}

type flightGenerator struct {
	store          *flights.Store
	rand           *rand.Rand
	counters       map[string]int
	createInterval time.Duration
	updateInterval time.Duration
	mu             sync.Mutex
}

func newFlightGenerator(store *flights.Store) *flightGenerator {
	gen := &flightGenerator{
		store:          store,
		rand:           rand.New(rand.NewSource(time.Now().UnixNano())),
		counters:       make(map[string]int),
		createInterval: 45 * time.Second,
		updateInterval: 15 * time.Second,
	}
	gen.bootstrapCounters()
	return gen
}

func (g *flightGenerator) start(ctx context.Context) {
	createTicker := time.NewTicker(g.createInterval)
	updateTicker := time.NewTicker(g.updateInterval)
	go func() {
		defer createTicker.Stop()
		defer updateTicker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-createTicker.C:
				g.maybeCreateFlights()
			case <-updateTicker.C:
				g.advanceFlights()
			}
		}
	}()
}

func (g *flightGenerator) bootstrapCounters() {
	airlines := g.store.ListAirlines()
	for _, airline := range airlines {
		flightsList, err := g.store.ListFlights(airline.AirlineID)
		if err != nil {
			continue
		}
		maxNum := 0
		for _, f := range flightsList {
			if n := parseFlightNumber(f.FlightID); n > maxNum {
				maxNum = n
			}
		}
		g.counters[airline.AirlineID] = maxNum
	}
}

func (g *flightGenerator) maybeCreateFlights() {
	airlines := g.store.ListAirlines()
	now := time.Now().Unix()
	for _, airline := range airlines {
		flightsList, err := g.store.ListFlights(airline.AirlineID)
		if err != nil {
			continue
		}
		scheduled := 0
		for _, f := range flightsList {
			if f.Status == flights.StatusScheduled && f.DepartureTimestamp > now {
				scheduled++
			}
		}
		if scheduled >= 3 {
			continue
		}
		flightID := g.nextFlightID(airline.AirlineID)
		minutesAhead := time.Duration(g.rand.Intn(4)+1) * time.Minute
		departure := now + int64(minutesAhead.Seconds())
		latest := now
		for _, f := range flightsList {
			if f.DepartureTimestamp > latest {
				latest = f.DepartureTimestamp
			}
		}
		minGap := int64((1 * time.Minute).Seconds())
		if latest >= departure {
			departure = latest + minGap
		}
		_, err = g.store.CreateFlight(airline.AirlineID, flights.Flight{
			AirlineID:          airline.AirlineID,
			FlightID:           flightID,
			DepartureTimestamp: departure,
			Status:             flights.StatusScheduled,
		})
		if err != nil {
			continue
		}
		slog.Info("auto-created flight", "airline", airline.AirlineID, "flight", flightID, "departure", departure)
	}
}

func (g *flightGenerator) advanceFlights() {
	airlines := g.store.ListAirlines()
	now := time.Now().Unix()
	for _, airline := range airlines {
		flightsList, err := g.store.ListFlights(airline.AirlineID)
		if err != nil {
			continue
		}
		for _, f := range flightsList {
			switch f.Status {
			case flights.StatusScheduled:
				if now < f.DepartureTimestamp {
					continue
				}
				if g.rand.Float64() < 0.4 {
					g.updateStatus(airline.AirlineID, f.FlightID, flights.StatusDelayed)
				} else {
					g.updateStatus(airline.AirlineID, f.FlightID, flights.StatusDeparted)
				}
			case flights.StatusDelayed:
				if now >= f.DepartureTimestamp+int64((1*time.Minute).Seconds()) {
					g.updateStatus(airline.AirlineID, f.FlightID, flights.StatusDeparted)
				}
			}
		}
	}
}

func (g *flightGenerator) updateStatus(airlineID, flightID string, status flights.Status) {
	if _, err := g.store.UpdateStatus(airlineID, flightID, status); err == nil {
		slog.Info("auto-updated flight", "airline", airlineID, "flight", flightID, "status", status)
	}
}

func (g *flightGenerator) nextFlightID(airlineID string) string {
	g.mu.Lock()
	defer g.mu.Unlock()
	n := g.counters[airlineID] + 1
	g.counters[airlineID] = n
	return fmt.Sprintf("%s-%03d", strings.ToUpper(airlineID), n)
}

func parseFlightNumber(id string) int {
	idx := strings.LastIndex(id, "-")
	if idx == -1 || idx == len(id)-1 {
		return 0
	}
	n, err := strconv.Atoi(id[idx+1:])
	if err != nil {
		return 0
	}
	return n
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
