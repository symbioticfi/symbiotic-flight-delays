package flights

import (
	"sort"
	"strings"
	"sync"
	"time"
)

// Store keeps airlines and flights in memory for the mock API.
type Store struct {
	mu       sync.RWMutex
	airlines map[string]Airline
	flights  map[string]map[string]*Flight // airlineID -> flightID -> Flight
}

// NewStore creates an in-memory store seeded with the provided airlines and flights.
func NewStore(initialAirlines []Airline, initialFlights []Flight) *Store {
	s := &Store{
		airlines: make(map[string]Airline),
		flights:  make(map[string]map[string]*Flight),
	}
	for _, airline := range initialAirlines {
		_ = s.AddAirline(airline)
	}
	s.mu.Lock()
	for _, flight := range initialFlights {
		if flight.Status == "" {
			flight.Status = StatusScheduled
		}
		_ = s.createFlightLocked(flight.AirlineID, flight)
	}
	s.mu.Unlock()
	return s
}

// ListAirlines returns airlines sorted alphabetically by code then name.
func (s *Store) ListAirlines() []Airline {
	s.mu.RLock()
	defer s.mu.RUnlock()
	items := make([]Airline, 0, len(s.airlines))
	for _, airline := range s.airlines {
		items = append(items, airline)
	}
	sort.Slice(items, func(i, j int) bool {
		if strings.EqualFold(items[i].Code, items[j].Code) {
			return items[i].Name < items[j].Name
		}
		return items[i].Code < items[j].Code
	})
	return items
}

// AddAirline registers a new airline.
func (s *Store) AddAirline(airline Airline) error {
	if strings.TrimSpace(airline.AirlineID) == "" {
		return ErrInvalidAirline
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.airlines[airline.AirlineID]; ok {
		return ErrAirlineExists
	}
	s.airlines[airline.AirlineID] = airline
	return nil
}

// ListFlights returns flights for the given airline sorted by departure.
func (s *Store) ListFlights(airlineID string) ([]Flight, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if _, ok := s.airlines[airlineID]; !ok {
		return nil, ErrAirlineNotFound
	}
	flightMap := s.flights[airlineID]
	items := make([]Flight, 0, len(flightMap))
	for _, flight := range flightMap {
		items = append(items, *flight)
	}
	sort.Slice(items, func(i, j int) bool { return items[i].DepartureTimestamp < items[j].DepartureTimestamp })
	return items, nil
}

// GetFlight returns a specific flight copy.
func (s *Store) GetFlight(airlineID, flightID string) (Flight, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	flightMap, ok := s.flights[airlineID]
	if !ok {
		return Flight{}, ErrAirlineNotFound
	}
	flight, ok := flightMap[flightID]
	if !ok {
		return Flight{}, ErrFlightNotFound
	}
	return *flight, nil
}

// CreateFlight registers a new flight for an airline.
func (s *Store) CreateFlight(airlineID string, flight Flight) (Flight, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.createFlightLocked(airlineID, flight); err != nil {
		return Flight{}, err
	}
	return *s.flights[airlineID][flight.FlightID], nil
}

func (s *Store) createFlightLocked(airlineID string, flight Flight) error {
	if flight.AirlineID == "" {
		flight.AirlineID = airlineID
	}
	if flight.Status == "" {
		flight.Status = StatusScheduled
	}
	if strings.TrimSpace(flight.FlightID) == "" {
		return ErrInvalidFlight
	}
	if !validStatus(flight.Status) {
		return ErrInvalidStatus
	}
	if _, ok := s.airlines[airlineID]; !ok {
		return ErrAirlineNotFound
	}
	if s.flights[airlineID] == nil {
		s.flights[airlineID] = make(map[string]*Flight)
	}
	if _, exists := s.flights[airlineID][flight.FlightID]; exists {
		return ErrFlightExists
	}
	flight.UpdatedAt = time.Now().Unix()
	copy := flight
	s.flights[airlineID][flight.FlightID] = &copy
	return nil
}

// UpdateStatus updates the status for a flight with basic validation.
func (s *Store) UpdateStatus(airlineID, flightID string, status Status) (Flight, error) {
	if !validStatus(status) {
		return Flight{}, ErrInvalidStatus
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	flightMap, ok := s.flights[airlineID]
	if !ok {
		return Flight{}, ErrAirlineNotFound
	}
	flight, ok := flightMap[flightID]
	if !ok {
		return Flight{}, ErrFlightNotFound
	}
	if !isValidTransition(flight.Status, status) {
		return Flight{}, ErrInvalidStatusTransition
	}
	flight.Status = status
	flight.UpdatedAt = time.Now().Unix()
	return *flight, nil
}

func isValidTransition(from, to Status) bool {
	switch from {
	case StatusScheduled:
		return to == StatusDelayed || to == StatusDeparted || to == StatusScheduled
	case StatusDelayed:
		return to == StatusDelayed
	case StatusDeparted:
		return to == StatusDeparted
	default:
		return false
	}
}
