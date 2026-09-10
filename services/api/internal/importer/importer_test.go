package importer_test

import (
	"strings"
	"testing"
	"time"

	"github.com/MauricioTAzevedo/finora-gemini-flash-student/services/api/internal/domain"
	"github.com/MauricioTAzevedo/finora-gemini-flash-student/services/api/internal/importer"
	"github.com/google/uuid"
)

const sampleOFX = `OFXHEADER:100
DATA:OFXSGML
VERSION:102
SECURITY:NONE
ENCODING:USASCII
CHARSET:1252
COMPRESSION:NONE
OLDFILEUID:NONE
NEWFILEUID:NONE

<OFX>
<BANKMSGSRSV1>
<STMTTRNRS>
<STMTRS>
<BANKTRANLIST>
<DTSTART>20260801
<DTEND>20260831
<STMTTRN>
<TRNTYPE>DEBIT
<DTPOSTED>20260810120000
<TRNAMT>-129.90
<FITID>20260810-001
<MEMO>AMZN MKTP BR
</STMTTRN>
<STMTTRN>
<TRNTYPE>DEBIT
<DTPOSTED>20260812120000
<TRNAMT>-340.50
<FITID>20260812-002
<MEMO>SUPERMERCADO HORIZONTE
</STMTTRN>
</BANKTRANLIST>
</STMTRS>
</STMTTRNRS>
</BANKMSGSRSV1>
</OFX>`

const sampleCSV = `Data;Descrição;Valor;Categoria
10/08/2026;AMZN MKTP BR;-129,90;Assinaturas & Lazer
12/08/2026;SUPERMERCADO HORIZONTE;-340,50;Alimentação & Supermercado
15/08/2026;SALARIO EMPRESA;8.500,00;Salário & Renda
`

func TestParseOFX(t *testing.T) {
	r := strings.NewReader(sampleOFX)
	records, err := importer.ParseOFX(r)
	if err != nil {
		t.Fatalf("unexpected error parsing OFX: %v", err)
	}

	if len(records) != 2 {
		t.Fatalf("expected 2 records, got %d", len(records))
	}

	r1 := records[0]
	if r1.Amount.AmountMinor != -12990 {
		t.Errorf("expected -12990 cents, got %d", r1.Amount.AmountMinor)
	}
	if r1.Description != "AMZN MKTP BR" {
		t.Errorf("expected 'AMZN MKTP BR', got %q", r1.Description)
	}
	if r1.ExternalID != "20260810-001" {
		t.Errorf("expected FITID '20260810-001', got %q", r1.ExternalID)
	}

	r2 := records[1]
	if r2.Amount.AmountMinor != -34050 {
		t.Errorf("expected -34050 cents, got %d", r2.Amount.AmountMinor)
	}
}

func TestParseCSV(t *testing.T) {
	r := strings.NewReader(sampleCSV)
	records, err := importer.ParseCSV(r)
	if err != nil {
		t.Fatalf("unexpected error parsing CSV: %v", err)
	}

	if len(records) != 3 {
		t.Fatalf("expected 3 records, got %d", len(records))
	}

	// Record 1: -R$ 129,90
	if records[0].Amount.AmountMinor != -12990 {
		t.Errorf("expected -12990 cents, got %d", records[0].Amount.AmountMinor)
	}

	// Record 3: +R$ 8.500,00
	if records[2].Amount.AmountMinor != 850000 {
		t.Errorf("expected 850000 cents, got %d", records[2].Amount.AmountMinor)
	}
}

func TestReconciliationAndDuplicatePrevention(t *testing.T) {
	householdID := uuid.New()
	checkingID := uuid.New()
	categoryID := uuid.New()

	// Existing manual ledger transaction: "Amazon" R$ 129,90 on 2026-08-08
	occurredAt, _ := time.Parse("2006-01-02", "2026-08-08")
	existingTx, _ := domain.NewExpense(
		householdID,
		checkingID,
		categoryID,
		domain.NewBRL(12990),
		"Amazon Marketplace",
		occurredAt,
		domain.SourceManual,
	)

	existingLedger := []domain.LedgerTransaction{*existingTx}

	// Case 1: Ingested statement transaction with slight description variation and +1 day distance
	// e.g. "AMZN MKTP BR" on 2026-08-09
	importDate, _ := time.Parse("2006-01-02", "2026-08-09")
	ingestedDuplicate := importer.CanonicalIngestionRecord{
		ExternalID:  "fit-001",
		Source:      "ofx",
		OccurredAt:  importDate,
		Description: "AMZN MKTP BR",
		Amount:      domain.NewBRL(12990),
	}

	result := importer.ReconcileTransaction(ingestedDuplicate, existingLedger)

	// INVARIANT: Must be detected as duplicate or candidate match to prevent balance inflation!
	if result.Status != importer.MatchStatusDuplicate && result.Status != importer.MatchStatusCandidate {
		t.Fatalf("CRITICAL FAILURE: Duplicate transaction was classified as %s (score: %.2f); expected duplicate/candidate", result.Status, result.Score)
	}

	if result.Score < 0.80 {
		t.Errorf("expected high match score (>= 0.80), got %.2f", result.Score)
	}

	// Case 2: Ingested brand new transaction: "Padaria Primavera" R$ 45,00
	newRecord := importer.CanonicalIngestionRecord{
		ExternalID:  "fit-002",
		Source:      "ofx",
		OccurredAt:  importDate,
		Description: "Padaria Primavera",
		Amount:      domain.NewBRL(4500),
	}

	resNew := importer.ReconcileTransaction(newRecord, existingLedger)
	if resNew.Status != importer.MatchStatusNew {
		t.Errorf("expected new transaction to be classified as 'new', got %s", resNew.Status)
	}
}
