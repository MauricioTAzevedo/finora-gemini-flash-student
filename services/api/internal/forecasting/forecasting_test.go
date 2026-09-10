package forecasting_test

import (
	"testing"

	"github.com/MauricioTAzevedo/finora-gemini-flash-student/services/api/internal/domain"
	"github.com/MauricioTAzevedo/finora-gemini-flash-student/services/api/internal/forecasting"
)

func TestDeterministicForecastMath(t *testing.T) {
	startingBalance := domain.NewBRL(842000) // R$ 8.420,00
	salaryMonthly := int64(850000)           // R$ 8.500,00
	salaryDay := 5
	billsMonthly := int64(200000)            // R$ 2.000,00
	billDay := 15
	cardDueDay := 10
	cardStatement := int64(312000)           // R$ 3.120,00
	horizonDays := 30
	reserveTarget := int64(500000)           // R$ 5.000,00

	projections, summary := forecasting.RunDeterministicForecast(
		startingBalance,
		salaryMonthly,
		salaryDay,
		billsMonthly,
		billDay,
		cardDueDay,
		cardStatement,
		horizonDays,
		reserveTarget,
	)

	if len(projections) != 30 {
		t.Fatalf("expected 30 daily projections, got %d", len(projections))
	}

	if summary.HorizonDays != 30 {
		t.Errorf("expected summary horizon 30, got %d", summary.HorizonDays)
	}

	// Verify all daily ending balances equal start + inflows - outflows
	for i, p := range projections {
		expectedEnd := p.StartingBalance.AmountMinor + p.InflowsConfirmed.AmountMinor - p.OutflowsConfirmed.AmountMinor
		if p.EndingBalance.AmountMinor != expectedEnd {
			t.Errorf("day %d: balance arithmetic discrepancy; expected %d, got %d", i, expectedEnd, p.EndingBalance.AmountMinor)
		}
	}
}

func TestCanWeAffordThisScenarioSimulation(t *testing.T) {
	// Scenario: Can we buy a laptop for R$ 6.000 in 12x of R$ 500?
	scenario := forecasting.Scenario{
		Name:                    "Notebook Novo para Home Office",
		TotalAmountMinor:        600000,
		InstallmentCount:        12,
		InstallmentMonthlyMinor: 50000, // R$ 500,00/mês
	}

	startingBalance := domain.NewBRL(1500000) // R$ 15.000,00
	salaryMonthly := int64(850000)            // R$ 8.500,00
	salaryDay := 5
	billsMonthly := int64(300000)             // R$ 3.000,00
	billDay := 15
	cardDueDay := 10
	cardStatement := int64(312000)            // R$ 3.120,00
	reserveTarget := int64(500000)            // R$ 5.000,00 minimum reserve

	comparison := forecasting.EvaluateAffordability(
		scenario,
		startingBalance,
		salaryMonthly,
		salaryDay,
		billsMonthly,
		billDay,
		cardDueDay,
		cardStatement,
		reserveTarget,
	)

	if comparison.ScenarioName != "Notebook Novo para Home Office" {
		t.Errorf("unexpected scenario name: %s", comparison.ScenarioName)
	}

	// The simulated lowest balance must be exactly R$ 500 (50000 cents) lower than baseline
	if comparison.LowestCashDeltaMinor != -50000 {
		t.Errorf("expected lowest cash delta -50000, got %d", comparison.LowestCashDeltaMinor)
	}

	// Because lowest balance remains above R$ 5.000,00 reserve, canAfford must be true
	if !comparison.CanAfford {
		t.Errorf("expected canAfford to be true")
	}

	if comparison.DecisionSupportNote == "" {
		t.Errorf("expected non-empty decision support note")
	}
}
