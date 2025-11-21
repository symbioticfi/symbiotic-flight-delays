package flights

import (
	"testing"
	"time"
)

func TestUpdateStatusAllowsDelayedToDeparted(t *testing.T) {
	store := NewStore(
		[]Airline{{AirlineID: "ALPHA", Name: "Alpha Air"}},
		[]Flight{{AirlineID: "ALPHA", FlightID: "ALPHA-1", DepartureTimestamp: time.Now().Unix(), Status: StatusScheduled}},
	)

	_, err := store.UpdateStatus("ALPHA", "ALPHA-1", StatusDelayed)
	if err != nil {
		t.Fatalf("expected delayed update to succeed, got %v", err)
	}

	updated, err := store.UpdateStatus("ALPHA", "ALPHA-1", StatusDeparted)
	if err != nil {
		t.Fatalf("expected delayed->departed transition to succeed, got %v", err)
	}
	if updated.Status != StatusDeparted {
		t.Fatalf("expected status %s, got %s", StatusDeparted, updated.Status)
	}
}

func TestUpdateStatusRejectsInvalidTransitions(t *testing.T) {
	now := time.Now().Unix()
	store := NewStore(
		[]Airline{{AirlineID: "BETA", Name: "Beta Wings"}},
		[]Flight{{AirlineID: "BETA", FlightID: "BETA-1", DepartureTimestamp: now, Status: StatusDeparted}},
	)

	if _, err := store.UpdateStatus("BETA", "BETA-1", StatusDelayed); err == nil {
		t.Fatalf("expected depart->delayed transition to fail")
	}
}
