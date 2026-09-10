package importer

import (
	"time"

	"github.com/MauricioTAzevedo/finora-gemini-flash-student/services/api/internal/domain"
)

type MatchStatus string

const (
	MatchStatusNew       MatchStatus = "new"       // No match found; ready to be posted as new transaction
	MatchStatusCandidate MatchStatus = "candidate" // Potential match; requires user review
	MatchStatusDuplicate MatchStatus = "duplicate" // Exact match found; will NOT inflate totals
)

// CanonicalIngestionRecord represents the normalized transaction structure before being committed to the ledger.
type CanonicalIngestionRecord struct {
	ExternalID     string       `json:"externalId"`
	Source         string       `json:"source"` // "ofx", "csv", "xlsx"
	OccurredAt     time.Time    `json:"occurredAt"`
	Description    string       `json:"description"`
	Amount         domain.Money `json:"amount"`
	AccountHint    string       `json:"accountHint,omitempty"`
	CategoryHint   string       `json:"categoryHint,omitempty"`
	RawPayload     string       `json:"rawPayload,omitempty"`
}

// ReconciliationResult holds match determination for an ingested row against existing ledger entries.
type ReconciliationResult struct {
	Record               CanonicalIngestionRecord `json:"record"`
	Status               MatchStatus              `json:"status"`
	Score                float64                  `json:"score"` // 0.0 to 1.0
	MatchedTransactionID string                   `json:"matchedTransactionId,omitempty"`
	Reason               string                   `json:"reason"`
}
