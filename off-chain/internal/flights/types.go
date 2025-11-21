package flights

import "errors"

// Status describes the lifecycle state of a flight tracked by the mock API.
type Status string

const (
	StatusScheduled Status = "SCHEDULED"
	StatusDelayed   Status = "DELAYED"
	StatusDeparted  Status = "DEPARTED"
)

// Airline represents a carrier that can have flights scheduled on-chain.
type Airline struct {
	AirlineID string `json:"airlineId"`
	Name      string `json:"name"`
	Code      string `json:"code"`
}

// Flight is a single tracked flight instance owned by an airline.
type Flight struct {
	AirlineID          string `json:"airlineId"`
	FlightID           string `json:"flightId"`
	DepartureTimestamp int64  `json:"departureTimestamp"`
	Status             Status `json:"status"`
	UpdatedAt          int64  `json:"updatedAt"`
}

var (
	ErrAirlineExists           = errors.New("airline already exists")
	ErrAirlineNotFound         = errors.New("airline not found")
	ErrInvalidAirline          = errors.New("invalid airline id")
	ErrFlightExists            = errors.New("flight already exists")
	ErrFlightNotFound          = errors.New("flight not found")
	ErrInvalidFlight           = errors.New("invalid flight id")
	ErrInvalidStatusTransition = errors.New("invalid flight status transition")
	ErrInvalidStatus           = errors.New("invalid flight status")
)

// validStatus reports whether the provided status is recognised by the API.
func validStatus(status Status) bool {
	switch status {
	case StatusScheduled, StatusDelayed, StatusDeparted:
		return true
	default:
		return false
	}
}
