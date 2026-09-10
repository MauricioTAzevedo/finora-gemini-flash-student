package forecasting

import (
	"fmt"
	"time"

	"github.com/MauricioTAzevedo/finora-gemini-flash-student/services/api/internal/domain"
)

// RunDeterministicForecast computes day-by-day cash flow balances over a specified horizon.
// It never relies on LLM approximations, ensuring exact penny-level mathematical integrity.
func RunDeterministicForecast(
	startingBalance domain.Money,
	salaryMonthlyMinor int64,
	salaryDay int,
	recurringMonthlyBillsMinor int64,
	billDay int,
	cardDueDay int,
	cardStatementMinor int64,
	horizonDays int,
	reserveTargetMinor int64,
) ([]DailyProjection, ForecastSummary) {
	now := time.Now().UTC()
	currentDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	projections := make([]DailyProjection, horizonDays)
	runningBalance := startingBalance

	lowestBalance := startingBalance
	lowestDateStr := currentDate.Format("02/01/2006")

	var totalInflows int64
	var totalOutflows int64

	for i := 0; i < horizonDays; i++ {
		date := currentDate.AddDate(0, 0, i)
		dayOfMonth := date.Day()

		var inConfirmed int64
		var outConfirmed int64

		// Monthly recurring salary
		if dayOfMonth == salaryDay {
			inConfirmed += salaryMonthlyMinor
		}

		// Monthly recurring fixed bills (e.g. rent, electricity, internet)
		if dayOfMonth == billDay {
			outConfirmed += recurringMonthlyBillsMinor
		}

		// Credit card statement due date payment
		if dayOfMonth == cardDueDay {
			outConfirmed += cardStatementMinor
		}

		startOfToday := runningBalance
		endOfTodayMinor := startOfToday.AmountMinor + inConfirmed - outConfirmed
		endOfToday := domain.NewBRL(endOfTodayMinor)

		totalInflows += inConfirmed
		totalOutflows += outConfirmed

		projections[i] = DailyProjection{
			Date:              date,
			DateFormatted:     date.Format("02/01/2006"),
			StartingBalance:   startOfToday,
			InflowsConfirmed:  domain.NewBRL(inConfirmed),
			OutflowsConfirmed: domain.NewBRL(outConfirmed),
			EndingBalance:     endOfToday,
		}

		if endOfToday.AmountMinor < lowestBalance.AmountMinor {
			lowestBalance = endOfToday
			lowestDateStr = date.Format("02/01/2006")
		}

		runningBalance = endOfToday
	}

	netChange := domain.NewBRL(runningBalance.AmountMinor - startingBalance.AmountMinor)

	summary := ForecastSummary{
		HorizonDays:             horizonDays,
		StartingBalance:         startingBalance,
		EndingBalance:           runningBalance,
		LowestProjectedBalance:  lowestBalance,
		LowestProjectedDate:     lowestDateStr,
		TotalInflows:            domain.NewBRL(totalInflows),
		TotalOutflows:           domain.NewBRL(totalOutflows),
		NetChange:               netChange,
		ReserveThresholdReached: lowestBalance.AmountMinor < reserveTargetMinor,
	}

	return projections, summary
}

// EvaluateAffordability simulates adding a new installment commitment (e.g., Laptop 12x R$ 500)
// and assesses its impact against household cash flow without mutating the actual ledger.
func EvaluateAffordability(
	scenario Scenario,
	startingBalance domain.Money,
	salaryMonthlyMinor int64,
	salaryDay int,
	recurringMonthlyBillsMinor int64,
	billDay int,
	cardDueDay int,
	cardStatementMinor int64,
	reserveTargetMinor int64,
) ScenarioComparison {
	horizonDays := 180 // 6-month simulation horizon

	// 1. Run Baseline Forecast
	_, baseline := RunDeterministicForecast(
		startingBalance,
		salaryMonthlyMinor,
		salaryDay,
		recurringMonthlyBillsMinor,
		billDay,
		cardDueDay,
		cardStatementMinor,
		horizonDays,
		reserveTargetMinor,
	)

	// 2. Run Scenario Forecast (Adding installment to credit card due day or bill day)
	simulatedBills := recurringMonthlyBillsMinor + scenario.InstallmentMonthlyMinor
	_, simulated := RunDeterministicForecast(
		startingBalance,
		salaryMonthlyMinor+scenario.MonthlyIncomeDeltaMinor,
		salaryDay,
		simulatedBills,
		billDay,
		cardDueDay,
		cardStatementMinor,
		horizonDays,
		reserveTargetMinor,
	)

	lowestDelta := simulated.LowestProjectedBalance.AmountMinor - baseline.LowestProjectedBalance.AmountMinor
	lowestDeltaFmt := domain.NewBRL(lowestDelta).FormatBRL()

	canAfford := simulated.LowestProjectedBalance.AmountMinor >= reserveTargetMinor

	var decisionNote string
	if canAfford {
		decisionNote = fmt.Sprintf(
			"A compra de '%s' adiciona %s/mês por %d parcelas. O menor saldo projetado em 6 meses passará de %s para %s, mantendo-se ACIMA da reserva mínima de segurança (%s).",
			scenario.Name,
			domain.NewBRL(scenario.InstallmentMonthlyMinor).FormatBRL(),
			scenario.InstallmentCount,
			baseline.LowestProjectedBalance.FormatBRL(),
			simulated.LowestProjectedBalance.FormatBRL(),
			domain.NewBRL(reserveTargetMinor).FormatBRL(),
		)
	} else {
		decisionNote = fmt.Sprintf(
			"Atenção: A compra de '%s' adiciona %s/mês por %d parcelas. Sob essa simulação, o menor saldo projetado atinge %s em %s, VIOLANDO a reserva mínima de segurança configurada pela família (%s).",
			scenario.Name,
			domain.NewBRL(scenario.InstallmentMonthlyMinor).FormatBRL(),
			scenario.InstallmentCount,
			simulated.LowestProjectedBalance.FormatBRL(),
			simulated.LowestProjectedDate,
			domain.NewBRL(reserveTargetMinor).FormatBRL(),
		)
	}

	return ScenarioComparison{
		ScenarioName:             scenario.Name,
		Baseline:                 baseline,
		Simulated:                simulated,
		LowestCashDeltaMinor:     lowestDelta,
		LowestCashDeltaFormatted: lowestDeltaFmt,
		CanAfford:                canAfford,
		DecisionSupportNote:      decisionNote,
	}
}
