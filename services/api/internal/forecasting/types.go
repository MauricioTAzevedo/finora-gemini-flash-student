package forecasting

import (
	"time"

	"github.com/MauricioTAzevedo/finora-gemini-flash-student/services/api/internal/domain"
)

// DailyProjection represents the projected financial balance for a specific day.
type DailyProjection struct {
	Date              time.Time    `json:"date"`
	DateFormatted     string       `json:"dateFormatted"`
	StartingBalance   domain.Money `json:"startingBalance"`
	InflowsConfirmed  domain.Money `json:"inflowsConfirmed"`
	InflowsEstimated  domain.Money `json:"inflowsEstimated"`
	OutflowsConfirmed domain.Money `json:"outflowsConfirmed"`
	OutflowsEstimated domain.Money `json:"outflowsEstimated"`
	EndingBalance     domain.Money `json:"endingBalance"`
}

// ForecastSummary highlights critical balance inflection points over a time horizon.
type ForecastSummary struct {
	HorizonDays             int          `json:"horizonDays"`
	StartingBalance         domain.Money `json:"startingBalance"`
	EndingBalance           domain.Money `json:"endingBalance"`
	LowestProjectedBalance  domain.Money `json:"lowestProjectedBalance"`
	LowestProjectedDate     string       `json:"lowestProjectedDate"`
	TotalInflows            domain.Money `json:"totalInflows"`
	TotalOutflows           domain.Money `json:"totalOutflows"`
	NetChange               domain.Money `json:"netChange"`
	ReserveThresholdReached bool         `json:"reserveThresholdReached"`
}

// Scenario defines hypothetical parameters for "What-If" simulations.
type Scenario struct {
	Name                    string `json:"name"`
	TotalAmountMinor        int64  `json:"totalAmountMinor"`
	InstallmentCount        int    `json:"installmentCount"`
	InstallmentMonthlyMinor int64  `json:"installmentMonthlyMinor"`
	MonthlyIncomeDeltaMinor int64  `json:"monthlyIncomeDeltaMinor"`
}

// ScenarioComparison presents the contrast between the baseline financial twin and the hypothetical scenario.
type ScenarioComparison struct {
	ScenarioName             string          `json:"scenarioName"`
	Baseline                 ForecastSummary `json:"baseline"`
	Simulated                ForecastSummary `json:"simulated"`
	LowestCashDeltaMinor     int64           `json:"lowestCashDeltaMinor"`
	LowestCashDeltaFormatted string          `json:"lowestCashDeltaFormatted"`
	CanAfford                bool            `json:"canAfford"`
	DecisionSupportNote      string          `json:"decisionSupportNote"`
}
