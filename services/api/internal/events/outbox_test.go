package events

import (
	"context"
	"encoding/json"
	"testing"
)

func TestOutboxRelayer(t *testing.T) {
	relayer := NewOutboxRelayer()

	var receivedEvents []Event
	relayer.Subscribe("TestExpenseCreated", func(ctx context.Context, evt Event) error {
		receivedEvents = append(receivedEvents, evt)
		return nil
	})

	rec := relayer.PublishSync(Event{
		AggregateType: "Ledger",
		AggregateID:   "tx-test-01",
		EventType:     "TestExpenseCreated",
		CorrelationID: "corr-123",
		CausationID:   "caus-123",
		HouseholdID:   "hh-demo-silva-001",
		Payload:       json.RawMessage(`{"amountMinor":15000}`),
	})

	if rec.Status != StatusPublished {
		t.Fatalf("expected status PUBLISHED, got %s", rec.Status)
	}

	if len(receivedEvents) != 1 {
		t.Fatalf("expected subscriber to receive 1 event, got %d", len(receivedEvents))
	}

	events := relayer.ListEvents("hh-demo-silva-001")
	if len(events) < 5 { // 4 seed events + 1 newly published
		t.Fatalf("expected at least 5 events in outbox list, got %d", len(events))
	}
}
