package events

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/google/uuid"
)

// EventStatus defines the processing lifecycle of an outbox event.
type EventStatus string

const (
	StatusPending   EventStatus = "PENDING"
	StatusPublished EventStatus = "PUBLISHED"
	StatusFailed    EventStatus = "FAILED"
)

// Event represents an immutable domain event with tracing context.
type Event struct {
	ID            string          `json:"id"`
	AggregateType string          `json:"aggregateType"`
	AggregateID   string          `json:"aggregateId"`
	EventType     string          `json:"eventType"`
	CorrelationID string          `json:"correlationId"`
	CausationID   string          `json:"causationId"`
	HouseholdID   string          `json:"householdId"`
	Payload       json.RawMessage `json:"payload"`
	CreatedAt     time.Time       `json:"createdAt"`
}

// OutboxRecord represents an event persisted in the outbox store awaiting relay.
type OutboxRecord struct {
	Event
	Status       EventStatus `json:"status"`
	RetryCount   int         `json:"retryCount"`
	ErrorMessage string      `json:"errorMessage,omitempty"`
	PublishedAt  *time.Time  `json:"publishedAt,omitempty"`
}

// EventHandler processes dispatched domain events.
type EventHandler func(ctx context.Context, evt Event) error

// OutboxRelayer coordinates transactional outbox publishing and event subscription.
type OutboxRelayer struct {
	mu          sync.RWMutex
	records     []*OutboxRecord
	subscribers map[string][]EventHandler
	stopChan    chan struct{}
}

// NewOutboxRelayer initializes a relayer instance with default synthetic seed events.
func NewOutboxRelayer() *OutboxRelayer {
	relayer := &OutboxRelayer{
		records:     make([]*OutboxRecord, 0),
		subscribers: make(map[string][]EventHandler),
		stopChan:    make(chan struct{}),
	}

	// Populate realistic domain events for demonstration and audit exploration
	now := time.Now().UTC()
	hID := "b0000000-0000-0000-0000-000000000001"

	relayer.PublishSync(Event{
		ID:            uuid.New().String(),
		AggregateType: "Household",
		AggregateID:   hID,
		EventType:     "HouseholdCreated",
		CorrelationID: "corr-init-001",
		CausationID:   "cmd-bootstrap-001",
		HouseholdID:   hID,
		Payload:       json.RawMessage(`{"name":"Família Silva","currency":"BRL","locale":"pt-BR"}`),
		CreatedAt:     now.Add(-48 * time.Hour),
	})

	relayer.PublishSync(Event{
		ID:            uuid.New().String(),
		AggregateType: "Account",
		AggregateID:   "acc-nubank-silva",
		EventType:     "AccountOpened",
		CorrelationID: "corr-init-002",
		CausationID:   "corr-init-001",
		HouseholdID:   hID,
		Payload:       json.RawMessage(`{"accountName":"Nubank Principal","accountType":"CHECKING"}`),
		CreatedAt:     now.Add(-47 * time.Hour),
	})

	relayer.PublishSync(Event{
		ID:            uuid.New().String(),
		AggregateType: "Ledger",
		AggregateID:   "tx-salary-aug",
		EventType:     "TransactionPosted",
		CorrelationID: "corr-tx-001",
		CausationID:   "corr-init-002",
		HouseholdID:   hID,
		Payload:       json.RawMessage(`{"description":"Salário Empresa Tech","amountMinor":1250000,"currency":"BRL"}`),
		CreatedAt:     now.Add(-24 * time.Hour),
	})

	relayer.PublishSync(Event{
		ID:            uuid.New().String(),
		AggregateType: "Intelligence",
		AggregateID:   "anom-cpfl-01",
		EventType:     "AnomalyDetected",
		CorrelationID: "corr-ai-001",
		CausationID:   "tx-cpfl-aug",
		HouseholdID:   hID,
		Payload:       json.RawMessage(`{"anomalyType":"SPIKE_UTILITY","merchant":"CPFL Paulista","zScore":2.4,"spikePercentage":62}`),
		CreatedAt:     now.Add(-2 * time.Hour),
	})

	return relayer
}

// Stage adds an event to the outbox queue in PENDING status.
func (r *OutboxRelayer) Stage(evt Event) *OutboxRecord {
	r.mu.Lock()
	defer r.mu.Unlock()

	if evt.ID == "" {
		evt.ID = uuid.New().String()
	}
	if evt.CreatedAt.IsZero() {
		evt.CreatedAt = time.Now().UTC()
	}

	record := &OutboxRecord{
		Event:  evt,
		Status: StatusPending,
	}
	r.records = append(r.records, record)
	return record
}

// PublishSync stages and immediately publishes an event through registered subscribers.
func (r *OutboxRelayer) PublishSync(evt Event) *OutboxRecord {
	rec := r.Stage(evt)
	r.dispatch(rec)
	return rec
}

// Subscribe registers a listener for a specific event type (or "*" for all events).
func (r *OutboxRelayer) Subscribe(eventType string, handler EventHandler) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.subscribers[eventType] = append(r.subscribers[eventType], handler)
}

func (r *OutboxRelayer) dispatch(rec *OutboxRecord) {
	ctx := context.Background()
	r.mu.RLock()
	handlers := append([]EventHandler{}, r.subscribers[rec.EventType]...)
	handlers = append(handlers, r.subscribers["*"]...)
	r.mu.RUnlock()

	for _, h := range handlers {
		if err := h(ctx, rec.Event); err != nil {
			rec.Status = StatusFailed
			rec.ErrorMessage = err.Error()
			rec.RetryCount++
			return
		}
	}

	now := time.Now().UTC()
	rec.Status = StatusPublished
	rec.PublishedAt = &now
}

// RelayPending iterates through pending events and publishes them (used by background worker).
func (r *OutboxRelayer) RelayPending(ctx context.Context) int {
	r.mu.Lock()
	var toProcess []*OutboxRecord
	for _, rec := range r.records {
		if rec.Status == StatusPending || (rec.Status == StatusFailed && rec.RetryCount < 3) {
			toProcess = append(toProcess, rec)
		}
	}
	r.mu.Unlock()

	count := 0
	for _, rec := range toProcess {
		r.dispatch(rec)
		count++
	}
	return count
}

// ListEvents returns recent events for the Event Explorer, optionally filtered by household.
func (r *OutboxRelayer) ListEvents(householdID string) []OutboxRecord {
	r.mu.RLock()
	defer r.mu.RUnlock()

	results := make([]OutboxRecord, 0)
	for i := len(r.records) - 1; i >= 0; i-- {
		rec := r.records[i]
		if householdID == "" || rec.HouseholdID == householdID {
			results = append(results, *rec)
		}
	}
	return results
}
