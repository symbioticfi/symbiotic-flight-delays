package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"sort"
	"strings"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"
	v1 "github.com/symbioticfi/relay/api/client/v1"

	"sum/internal/contracts"
	"sum/internal/flights"
	"sum/internal/utils"
)

const (
	keyTag    = 15
	maxUint48 = (1 << 48) - 1
)

type config struct {
	relayAPIURL       string
	evmRPCURL         string
	contractAddress   string
	flightsAPIURL     string
	privateKeyHex     string
	pollInterval      time.Duration
	proofPollInterval time.Duration
	logLevel          string
}

var cfg config

var rootCmd = &cobra.Command{
	Use:           "flight-node",
	Short:         "Flight delay oracle node",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		switch strings.ToLower(cfg.logLevel) {
		case "debug":
			slog.SetLogLoggerLevel(slog.LevelDebug)
		case "warn":
			slog.SetLogLoggerLevel(slog.LevelWarn)
		case "error":
			slog.SetLogLoggerLevel(slog.LevelError)
		default:
			slog.SetLogLoggerLevel(slog.LevelInfo)
		}

		ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer cancel()

		conn, err := utils.GetGRPCConnection(cfg.relayAPIURL)
		if err != nil {
			return fmt.Errorf("create relay client: %w", err)
		}
		defer conn.Close()

		relayClient := v1.NewSymbioticClient(conn)

		evmClient, err := ethclient.DialContext(ctx, cfg.evmRPCURL)
		if err != nil {
			return fmt.Errorf("dial evm rpc: %w", err)
		}
		defer evmClient.Close()

		chainID, err := evmClient.ChainID(ctx)
		if err != nil {
			return fmt.Errorf("fetch chain id: %w", err)
		}

		contractAddr := common.HexToAddress(cfg.contractAddress)
		flightDelays, err := contracts.NewFlightDelays(contractAddr, evmClient)
		if err != nil {
			return fmt.Errorf("bind flight delays: %w", err)
		}

		flightsClient := newFlightsAPIClient(cfg.flightsAPIURL)

		privKey, err := crypto.HexToECDSA(strings.TrimPrefix(cfg.privateKeyHex, "0x"))
		if err != nil {
			return fmt.Errorf("parse private key: %w", err)
		}

		node := &flightNode{
			relayClient: relayClient,
			ethClient:   evmClient,
			contract:    flightDelays,
			chainID:     chainID,
			privateKey:  privKey,
			flightsAPI:  flightsClient,
			pending:     make(map[string]*pendingAction),
		}

		if err := node.syncFlights(ctx); err != nil {
			slog.Warn("initial sync failed", "error", err)
		}

		pollTicker := time.NewTicker(cfg.pollInterval)
		defer pollTicker.Stop()
		proofTicker := time.NewTicker(cfg.proofPollInterval)
		defer proofTicker.Stop()

		for {
			select {
			case <-pollTicker.C:
				if err := node.syncFlights(ctx); err != nil {
					slog.Warn("sync flights failed", "error", err)
				}
			case <-proofTicker.C:
				if err := node.fetchProofs(ctx); err != nil {
					slog.Warn("fetch proofs failed", "error", err)
				}
				if err := node.submitReadyActions(ctx); err != nil {
					slog.Warn("submit actions failed", "error", err)
				}
			case <-ctx.Done():
				slog.Info("shutting down flight node")
				return nil
			}
		}
	},
}

func main() {
	rootCmd.PersistentFlags().StringVar(&cfg.relayAPIURL, "relay-api-url", "", "Relay API URL")
	rootCmd.PersistentFlags().StringVar(&cfg.evmRPCURL, "evm-rpc-url", "", "Execution client RPC URL")
	rootCmd.PersistentFlags().StringVar(&cfg.contractAddress, "flight-delays-address", "", "FlightDelays contract address")
	rootCmd.PersistentFlags().StringVar(&cfg.flightsAPIURL, "flights-api-url", "", "Mock flights API URL")
	rootCmd.PersistentFlags().StringVar(&cfg.privateKeyHex, "private-key", "", "Flight oracle ECDSA private key")
	rootCmd.PersistentFlags().DurationVar(&cfg.pollInterval, "poll-interval", 5*time.Second, "Polling interval for flights API")
	rootCmd.PersistentFlags().DurationVar(&cfg.proofPollInterval, "proof-poll-interval", 3*time.Second, "Polling interval for settlement proofs")
	rootCmd.PersistentFlags().StringVar(&cfg.logLevel, "log-level", "info", "Log level (debug,info,warn,error)")

	_ = rootCmd.MarkPersistentFlagRequired("relay-api-url")
	_ = rootCmd.MarkPersistentFlagRequired("evm-rpc-url")
	_ = rootCmd.MarkPersistentFlagRequired("flight-delays-address")
	_ = rootCmd.MarkPersistentFlagRequired("flights-api-url")
	_ = rootCmd.MarkPersistentFlagRequired("private-key")

	if err := rootCmd.Execute(); err != nil {
		slog.Error("node failed", "error", err)
		os.Exit(1)
	}
}

var (
	bytes32Type     = mustABIType("bytes32")
	uint48Type      = mustABIType("uint48")
	outerArgs       = abi.Arguments{{Type: bytes32Type}}
	createInnerArgs = abi.Arguments{{Type: bytes32Type}, {Type: bytes32Type}, {Type: bytes32Type}, {Type: uint48Type}, {Type: bytes32Type}}
	statusInnerArgs = abi.Arguments{{Type: bytes32Type}, {Type: bytes32Type}, {Type: bytes32Type}}
	createTypehash  = crypto.Keccak256Hash([]byte("Create(bytes32 airlineId,bytes32 flightId,uint48 departure,bytes32 previousFlightId)"))
	delayTypehash   = crypto.Keccak256Hash([]byte("Delay(bytes32 airlineId,bytes32 flightId)"))
	departTypehash  = crypto.Keccak256Hash([]byte("Depart(bytes32 airlineId,bytes32 flightId)"))
)

func mustABIType(name string) abi.Type {
	t, err := abi.NewType(name, "", nil)
	if err != nil {
		panic(err)
	}
	return t
}

type actionType string

type flightStatus uint8

const (
	actionCreate actionType = "CREATE"
	actionDelay  actionType = "DELAY"
	actionDepart actionType = "DEPART"

	statusNone      flightStatus = 0
	statusScheduled flightStatus = 1
	statusDelayed   flightStatus = 2
	statusDeparted  flightStatus = 3
)

type pendingAction struct {
	Key                string
	Airline            flights.Airline
	Flight             flights.Flight
	AirlineHash        common.Hash
	FlightHash         common.Hash
	PreviousFlightHash common.Hash
	Type               actionType
	Epoch              uint64
	RequestID          string
	Proof              []byte
	Submitted          bool
	TargetStatus       flightStatus
	TxHash             common.Hash
	CreatedAt          time.Time
}

type flightNode struct {
	relayClient *v1.SymbioticClient
	ethClient   *ethclient.Client
	contract    *contracts.FlightDelays
	chainID     *big.Int
	privateKey  *ecdsa.PrivateKey
	flightsAPI  *flightsAPIClient

	pending map[string]*pendingAction
}

func (n *flightNode) syncFlights(ctx context.Context) error {
	airlines, err := n.flightsAPI.ListAirlines(ctx)
	if err != nil {
		return fmt.Errorf("list airlines: %w", err)
	}
	for _, airline := range airlines {
		flightsForAirline, err := n.flightsAPI.ListFlights(ctx, airline.AirlineID)
		if err != nil {
			slog.Warn("list flights failed", "airline", airline.AirlineID, "error", err)
			continue
		}
		sort.SliceStable(flightsForAirline, func(i, j int) bool {
			if flightsForAirline[i].DepartureTimestamp == flightsForAirline[j].DepartureTimestamp {
				return flightsForAirline[i].FlightID < flightsForAirline[j].FlightID
			}
			return flightsForAirline[i].DepartureTimestamp < flightsForAirline[j].DepartureTimestamp
		})
		if err := n.evaluateAirlineFlights(ctx, airline, flightsForAirline); err != nil {
			slog.Warn("evaluate airline failed", "airline", airline.AirlineID, "error", err)
		}
	}
	return nil
}

func (n *flightNode) evaluateAirlineFlights(ctx context.Context, airline flights.Airline, flightsForAirline []flights.Flight) error {
	airlineHash := hashIdentifier(airline.AirlineID)
	prevMap := make(map[string]common.Hash, len(flightsForAirline))
	var prev common.Hash
	for _, flight := range flightsForAirline {
		prevMap[flight.FlightID] = prev
		prev = hashIdentifier(flight.FlightID)
	}

	for _, flight := range flightsForAirline {
		prevHash := prevMap[flight.FlightID]
		if err := n.evaluateFlight(ctx, airline, flight, airlineHash, prevHash); err != nil {
			slog.Warn("evaluate flight failed", "airline", airline.AirlineID, "flight", flight.FlightID, "error", err)
		}
	}
	return nil
}

func (n *flightNode) evaluateFlight(ctx context.Context, airline flights.Airline, flight flights.Flight, airlineHash common.Hash, previousFlightHash common.Hash) error {
	flightHash := hashIdentifier(flight.FlightID)

	info, err := n.contract.Flights(&bind.CallOpts{Context: ctx}, airlineHash, flightHash)
	if err != nil {
		return fmt.Errorf("read flight state: %w", err)
	}

	onChainStatus := flightStatus(info.Status)
	n.clearSatisfiedPending(airlineHash, flightHash, onChainStatus)

	nextAction, ok := determineAction(flight.Status, onChainStatus)
	if !ok {
		return nil
	}

	key := actionKey(airlineHash, flightHash, nextAction)
	if _, exists := n.pending[key]; exists {
		return nil
	}

	if err := n.enqueueAction(ctx, key, nextAction, airline, flight, airlineHash, flightHash, previousFlightHash); err != nil {
		return err
	}
	return nil
}

func determineAction(apiStatus flights.Status, onChain flightStatus) (actionType, bool) {
	switch onChain {
	case statusNone:
		if apiStatus == flights.StatusScheduled || apiStatus == flights.StatusDelayed || apiStatus == flights.StatusDeparted {
			return actionCreate, true
		}
	case statusScheduled:
		switch apiStatus {
		case flights.StatusDelayed:
			return actionDelay, true
		case flights.StatusDeparted:
			return actionDepart, true
		}
	}
	return "", false
}

func (n *flightNode) enqueueAction(ctx context.Context, key string, action actionType, airline flights.Airline, flight flights.Flight, airlineHash, flightHash, previousFlightHash common.Hash) error {
	if action == actionCreate && flight.DepartureTimestamp <= 0 {
		return fmt.Errorf("flight %s has invalid departure timestamp", flight.FlightID)
	}
	payload, err := buildMessagePayload(action, airlineHash, flightHash, previousFlightHash, uint64(flight.DepartureTimestamp))
	if err != nil {
		return err
	}

	epoch, requestID, err := n.requestSignature(ctx, payload)
	if err != nil {
		return fmt.Errorf("sign message: %w", err)
	}

	pending := &pendingAction{
		Key:                key,
		Airline:            airline,
		Flight:             flight,
		AirlineHash:        airlineHash,
		FlightHash:         flightHash,
		PreviousFlightHash: previousFlightHash,
		Type:               action,
		Epoch:              epoch,
		RequestID:          requestID,
		TargetStatus:       targetStatusFor(action),
		CreatedAt:          time.Now(),
	}
	n.pending[key] = pending

	slog.Info("scheduled flight action", "airline", airline.AirlineID, "flight", flight.FlightID, "action", string(action), "epoch", epoch)
	return nil
}

func (n *flightNode) requestSignature(ctx context.Context, payload []byte) (uint64, string, error) {
	epochInfos, err := n.relayClient.GetLastAllCommitted(ctx, &v1.GetLastAllCommittedRequest{})
	if err != nil {
		return 0, "", fmt.Errorf("last committed: %w", err)
	}
	var suggestedEpoch uint64
	for _, info := range epochInfos.EpochInfos {
		last := info.GetLastCommittedEpoch()
		if suggestedEpoch == 0 || last < suggestedEpoch {
			suggestedEpoch = last
		}
	}

	resp, err := n.relayClient.SignMessage(ctx, &v1.SignMessageRequest{
		KeyTag:        keyTag,
		Message:       payload,
		RequiredEpoch: &suggestedEpoch,
	})
	if err != nil {
		return 0, "", err
	}
	return resp.Epoch, resp.RequestId, nil
}

func (n *flightNode) fetchProofs(ctx context.Context) error {
	for _, action := range n.pending {
		if action.Proof != nil {
			continue
		}
		resp, err := n.relayClient.GetAggregationProof(ctx, &v1.GetAggregationProofRequest{RequestId: action.RequestID})
		if err != nil {
			continue
		}
		if resp.GetAggregationProof() == nil {
			continue
		}
		action.Proof = resp.AggregationProof.Proof
		slog.Info("aggregation proof ready", "airline", action.Airline.AirlineID, "flight", action.Flight.FlightID, "action", string(action.Type))
	}
	return nil
}

func (n *flightNode) submitReadyActions(ctx context.Context) error {
	for key, action := range n.pending {
		if action.Proof == nil || action.Submitted {
			continue
		}
		if action.Type == actionCreate {
			ready, err := n.canSubmitCreate(ctx, action)
			if err != nil {
				slog.Warn("check create readiness failed", "airline", action.Airline.AirlineID, "flight", action.Flight.FlightID, "error", err)
				continue
			}
			if !ready {
				continue
			}
		}
		if err := n.submitAction(ctx, action); err != nil {
			slog.Warn("submit action failed", "airline", action.Airline.AirlineID, "flight", action.Flight.FlightID, "action", string(action.Type), "error", err)
			continue
		}
		n.pending[key] = action
	}
	return nil
}

func (n *flightNode) submitAction(ctx context.Context, action *pendingAction) error {
	txOpts, err := bind.NewKeyedTransactorWithChainID(n.privateKey, n.chainID)
	if err != nil {
		return err
	}
	txOpts.Context = ctx

	epoch := big.NewInt(int64(action.Epoch))
	var txHash common.Hash

	switch action.Type {
	case actionCreate:
		scheduled := big.NewInt(action.Flight.DepartureTimestamp)
		prev := [32]byte(action.PreviousFlightHash)
		tx, err := n.contract.CreateFlight(txOpts, action.AirlineHash, action.FlightHash, scheduled, prev, epoch, action.Proof)
		if err != nil {
			return err
		}
		txHash = tx.Hash()
	case actionDelay:
		tx, err := n.contract.DelayFlight(txOpts, action.AirlineHash, action.FlightHash, epoch, action.Proof)
		if err != nil {
			return err
		}
		txHash = tx.Hash()
	case actionDepart:
		tx, err := n.contract.DepartFlight(txOpts, action.AirlineHash, action.FlightHash, epoch, action.Proof)
		if err != nil {
			return err
		}
		txHash = tx.Hash()
	default:
		return fmt.Errorf("unknown action %s", action.Type)
	}

	action.Submitted = true
	action.TxHash = txHash
	slog.Info("submitted flight action", "airline", action.Airline.AirlineID, "flight", action.Flight.FlightID, "action", string(action.Type), "tx", txHash.Hex())
	return nil
}

func (n *flightNode) canSubmitCreate(ctx context.Context, action *pendingAction) (bool, error) {
	airlineState, err := n.contract.Airlines(&bind.CallOpts{Context: ctx}, action.AirlineHash)
	if err != nil {
		return false, fmt.Errorf("read airline state: %w", err)
	}
	current := common.BytesToHash(airlineState.LastFlightId[:])
	return current == action.PreviousFlightHash, nil
}

func (n *flightNode) clearSatisfiedPending(airlineHash, flightHash common.Hash, status flightStatus) {
	for key, action := range n.pending {
		if action.AirlineHash == airlineHash && action.FlightHash == flightHash && action.TargetStatus == status {
			slog.Info("action confirmed on-chain", "airline", action.Airline.AirlineID, "flight", action.Flight.FlightID, "action", string(action.Type))
			delete(n.pending, key)
		}
	}
}

func targetStatusFor(action actionType) flightStatus {
	switch action {
	case actionCreate:
		return statusScheduled
	case actionDelay:
		return statusDelayed
	case actionDepart:
		return statusDeparted
	default:
		return statusNone
	}
}

func buildMessagePayload(action actionType, airlineHash, flightHash, previousFlightHash common.Hash, departure uint64) ([]byte, error) {
	var inner []byte
	var err error
	switch action {
	case actionCreate:
		if departure == 0 || departure > maxUint48 {
			return nil, fmt.Errorf("invalid departure timestamp %d", departure)
		}
		inner, err = createInnerArgs.Pack(createTypehash, airlineHash, flightHash, big.NewInt(int64(departure)), previousFlightHash)
	case actionDelay:
		inner, err = statusInnerArgs.Pack(delayTypehash, airlineHash, flightHash)
	case actionDepart:
		inner, err = statusInnerArgs.Pack(departTypehash, airlineHash, flightHash)
	default:
		return nil, errors.New("unsupported action")
	}
	if err != nil {
		return nil, err
	}
	return inner, nil
}

func hashIdentifier(id string) common.Hash {
	normalized := strings.ToUpper(strings.TrimSpace(id))
	return crypto.Keccak256Hash([]byte(normalized))
}

func actionKey(airlineHash, flightHash common.Hash, action actionType) string {
	return fmt.Sprintf("%s|%s|%s", airlineHash.Hex(), flightHash.Hex(), action)
}

type flightsAPIClient struct {
	baseURL    string
	httpClient *http.Client
}

func newFlightsAPIClient(baseURL string) *flightsAPIClient {
	return &flightsAPIClient{
		baseURL:    strings.TrimSuffix(baseURL, "/"),
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}
}

func (c *flightsAPIClient) ListAirlines(ctx context.Context) ([]flights.Airline, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/airlines", nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("airlines request status %d", resp.StatusCode)
	}
	var body struct {
		Airlines []flights.Airline `json:"airlines"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}
	return body.Airlines, nil
}

func (c *flightsAPIClient) ListFlights(ctx context.Context, airlineID string) ([]flights.Flight, error) {
	endpoint := fmt.Sprintf("%s/airlines/%s/flights", c.baseURL, url.PathEscape(airlineID))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("flights request status %d", resp.StatusCode)
	}
	var body struct {
		Flights []flights.Flight `json:"flights"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}
	return body.Flights, nil
}
