package intelligence

import (
	"testing"
	"time"
)

func TestSubscriptionDetector(t *testing.T) {
	detector := NewSubscriptionDetector()

	// 3 months of Netflix with a slight price rise
	txs := []TransactionRecord{
		{Description: "Netflix Assinatura", Category: "Lazer & Streaming", AmountMinor: 3990, Currency: "BRL", Date: time.Date(2026, 6, 10, 0, 0, 0, 0, time.UTC)},
		{Description: "Netflix Mensalidade", Category: "Lazer & Streaming", AmountMinor: 4490, Currency: "BRL", Date: time.Date(2026, 7, 10, 0, 0, 0, 0, time.UTC)},
		{Description: "Netflix.com", Category: "Lazer & Streaming", AmountMinor: 4490, Currency: "BRL", Date: time.Date(2026, 8, 10, 0, 0, 0, 0, time.UTC)},
		// Spotify regular monthly
		{Description: "Spotify Brasil", Category: "Lazer & Streaming", AmountMinor: 2190, Currency: "BRL", Date: time.Date(2026, 6, 5, 0, 0, 0, 0, time.UTC)},
		{Description: "Spotify Premium", Category: "Lazer & Streaming", AmountMinor: 2190, Currency: "BRL", Date: time.Date(2026, 7, 5, 0, 0, 0, 0, time.UTC)},
		{Description: "Spotify", Category: "Lazer & Streaming", AmountMinor: 2190, Currency: "BRL", Date: time.Date(2026, 8, 5, 0, 0, 0, 0, time.UTC)},
		// One-off purchase (should NOT be detected)
		{Description: "Restaurante Fogo de Chao", Category: "Lazer", AmountMinor: 35000, Currency: "BRL", Date: time.Date(2026, 7, 15, 0, 0, 0, 0, time.UTC)},
	}

	subs := detector.Detect(txs)
	if len(subs) != 2 {
		t.Fatalf("expected 2 subscriptions, got %d", len(subs))
	}

	// Check Netflix detected
	var foundNetflix, foundSpotify bool
	for _, s := range subs {
		if s.MerchantName == "Netflix" {
			foundNetflix = true
			if s.Cadence != CadenceMonthly {
				t.Errorf("expected Netflix cadence to be MONTHLY, got %s", s.Cadence)
			}
			if s.Amount.AmountMinor != 4490 {
				t.Errorf("expected Netflix last amount 4490, got %d", s.Amount.AmountMinor)
			}
			if s.AnnualizedCost.AmountMinor != 4490*12 {
				t.Errorf("expected annualized cost %d, got %d", 4490*12, s.AnnualizedCost.AmountMinor)
			}
			if s.PriceDriftNote == "" {
				t.Errorf("expected price drift note for Netflix increase from 3990 to 4490")
			}
		}
		if s.MerchantName == "Spotify" {
			foundSpotify = true
			if s.Cadence != CadenceMonthly {
				t.Errorf("expected Spotify cadence to be MONTHLY, got %s", s.Cadence)
			}
			if s.Amount.AmountMinor != 2190 {
				t.Errorf("expected Spotify amount 2190, got %d", s.Amount.AmountMinor)
			}
		}
	}

	if !foundNetflix || !foundSpotify {
		t.Fatalf("expected both Netflix and Spotify to be detected")
	}
}

func TestAnomalyDetector(t *testing.T) {
	detector := NewAnomalyDetector()

	// CPFL bills with a big spike in August (420 vs ~250 average)
	txs := []TransactionRecord{
		{Description: "CPFL Paulista", Category: "Moradia & Utilidades", AmountMinor: 24500, Currency: "BRL", Date: time.Date(2026, 4, 15, 0, 0, 0, 0, time.UTC)},
		{Description: "CPFL Paulista", Category: "Moradia & Utilidades", AmountMinor: 26000, Currency: "BRL", Date: time.Date(2026, 5, 15, 0, 0, 0, 0, time.UTC)},
		{Description: "CPFL Paulista", Category: "Moradia & Utilidades", AmountMinor: 25000, Currency: "BRL", Date: time.Date(2026, 6, 15, 0, 0, 0, 0, time.UTC)},
		{Description: "CPFL Paulista", Category: "Moradia & Utilidades", AmountMinor: 25500, Currency: "BRL", Date: time.Date(2026, 7, 15, 0, 0, 0, 0, time.UTC)},
		{Description: "CPFL Paulista", Category: "Moradia & Utilidades", AmountMinor: 48000, Currency: "BRL", Date: time.Date(2026, 8, 15, 0, 0, 0, 0, time.UTC)}, // Spike!

		// Duplicate charge on the same day
		{Description: "Farmacia Raia", Category: "Saúde", AmountMinor: 8990, Currency: "BRL", Date: time.Date(2026, 8, 20, 0, 0, 0, 0, time.UTC)},
		{Description: "Farmacia Raia", Category: "Saúde", AmountMinor: 8990, Currency: "BRL", Date: time.Date(2026, 8, 20, 0, 0, 0, 0, time.UTC)},
	}

	anomalies := detector.Detect(txs)
	if len(anomalies) < 2 {
		t.Fatalf("expected at least 2 anomalies (utility spike and duplicate charge), got %d", len(anomalies))
	}

	var hasUtilitySpike, hasDupe bool
	for _, a := range anomalies {
		if a.Type == AnomalySpikeUtility {
			hasUtilitySpike = true
			if a.Amount.AmountMinor != 48000 {
				t.Errorf("expected spike amount 48000, got %d", a.Amount.AmountMinor)
			}
			if a.ZScore < 1.8 {
				t.Errorf("expected z-score >= 1.8, got %f", a.ZScore)
			}
		}
		if a.Type == AnomalyPotentialDupe {
			hasDupe = true
			if a.Amount.AmountMinor != 8990 {
				t.Errorf("expected duplicate amount 8990, got %d", a.Amount.AmountMinor)
			}
		}
	}

	if !hasUtilitySpike {
		t.Errorf("failed to detect CPFL utility spike")
	}
	if !hasDupe {
		t.Errorf("failed to detect same-day duplicate charge")
	}
}

func TestNaturalLanguageInterpreter(t *testing.T) {
	now := time.Date(2026, 9, 10, 12, 0, 0, 0, time.UTC)
	interp := NewNaturalLanguageInterpreter(now)

	txs := []TransactionRecord{
		{Description: "Pão de Açúcar", Category: "Alimentação & Mercado", AmountMinor: 45000, Currency: "BRL", Date: now.AddDate(0, -1, 0)},
		{Description: "Carrefour Express", Category: "Alimentação & Mercado", AmountMinor: 12000, Currency: "BRL", Date: now.AddDate(0, -2, 0)},
		{Description: "Posto Ipiranga", Category: "Transporte & Mobilidade", AmountMinor: 22000, Currency: "BRL", Date: now.AddDate(0, -1, 0)},
		{Description: "Pão de Açúcar Grande", Category: "Alimentação & Mercado", AmountMinor: 85000, Currency: "BRL", Date: now.AddDate(0, -8, 0)}, // Older than 6 months
	}

	// Test query: "Quanto gastamos com mercado nos últimos 3 meses?"
	filter := interp.Parse("Quanto gastamos com mercado nos últimos 3 meses?")
	if filter.Category != "Alimentação & Mercado" {
		t.Errorf("expected category Alimentação & Mercado, got %s", filter.Category)
	}
	if filter.StartDate == nil {
		t.Fatalf("expected StartDate to be set for últimos 3 meses")
	}

	result := interp.Execute(filter, txs)
	// Should match 2 transactions (45000 + 12000 = 57000 cents)
	if result.TotalCount != 2 {
		t.Fatalf("expected 2 matches, got %d", result.TotalCount)
	}
	if result.TotalAmount.AmountMinor != 57000 {
		t.Errorf("expected total 57000, got %d", result.TotalAmount.AmountMinor)
	}

	// Test amount filter query: "gastos acima de 200 no transporte"
	f2 := interp.Parse("gastos acima de 200 no transporte")
	if f2.MinAmountMinor != 20000 {
		t.Errorf("expected MinAmountMinor 20000, got %d", f2.MinAmountMinor)
	}
	r2 := interp.Execute(f2, txs)
	if r2.TotalCount != 1 {
		t.Fatalf("expected 1 match for transporte > 200, got %d", r2.TotalCount)
	}
}
