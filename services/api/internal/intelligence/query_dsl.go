package intelligence

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/MauricioTAzevedo/finora-gemini-flash-student/services/api/internal/domain"
)

// QueryFilter defines safe, validated query criteria.
type QueryFilter struct {
	Category       string       `json:"category,omitempty"`
	Merchant       string       `json:"merchant,omitempty"`
	StartDate      *time.Time   `json:"startDate,omitempty"`
	EndDate        *time.Time   `json:"endDate,omitempty"`
	MinAmountMinor int64        `json:"minAmountMinor,omitempty"`
	MaxAmountMinor int64        `json:"maxAmountMinor,omitempty"`
	NaturalQuery   string       `json:"naturalQuery"`
	Explanation    string       `json:"explanation"`
}

// QueryResult represents the execution result of the parsed DSL query.
type QueryResult struct {
	Filter       QueryFilter         `json:"filter"`
	TotalCount   int                 `json:"totalCount"`
	TotalAmount  domain.Money        `json:"totalAmount"`
	SummaryText  string              `json:"summaryText"`
	Transactions []TransactionRecord `json:"transactions"`
}

// NaturalLanguageInterpreter parses human queries into safe AST filters.
type NaturalLanguageInterpreter struct {
	now time.Time
}

// NewNaturalLanguageInterpreter returns a new interpreter.
func NewNaturalLanguageInterpreter(now time.Time) *NaturalLanguageInterpreter {
	return &NaturalLanguageInterpreter{now: now}
}

// Parse translates a Portuguese natural language string into a safe QueryFilter.
func (n *NaturalLanguageInterpreter) Parse(query string) QueryFilter {
	q := strings.ToLower(strings.TrimSpace(query))
	filter := QueryFilter{
		NaturalQuery: query,
	}

	var explanationParts []string

	// 1. Detect Timeframe
	if strings.Contains(q, "últimos 3 meses") || strings.Contains(q, "ultimos 3 meses") || strings.Contains(q, "3 meses") {
		start := n.now.AddDate(0, -3, 0)
		filter.StartDate = &start
		filter.EndDate = &n.now
		explanationParts = append(explanationParts, "período: últimos 3 meses")
	} else if strings.Contains(q, "últimos 6 meses") || strings.Contains(q, "ultimos 6 meses") || strings.Contains(q, "6 meses") {
		start := n.now.AddDate(0, -6, 0)
		filter.StartDate = &start
		filter.EndDate = &n.now
		explanationParts = append(explanationParts, "período: últimos 6 meses")
	} else if strings.Contains(q, "este mês") || strings.Contains(q, "esse mês") || strings.Contains(q, "neste mês") {
		start := time.Date(n.now.Year(), n.now.Month(), 1, 0, 0, 0, 0, n.now.Location())
		filter.StartDate = &start
		filter.EndDate = &n.now
		explanationParts = append(explanationParts, "período: mês atual")
	} else if strings.Contains(q, "mês passado") || strings.Contains(q, "mes passado") {
		firstThisMonth := time.Date(n.now.Year(), n.now.Month(), 1, 0, 0, 0, 0, n.now.Location())
		start := firstThisMonth.AddDate(0, -1, 0)
		end := firstThisMonth.AddDate(0, 0, -1)
		filter.StartDate = &start
		filter.EndDate = &end
		explanationParts = append(explanationParts, "período: mês anterior")
	} else if strings.Contains(q, "última semana") || strings.Contains(q, "semana passada") {
		start := n.now.AddDate(0, 0, -7)
		filter.StartDate = &start
		filter.EndDate = &n.now
		explanationParts = append(explanationParts, "período: últimos 7 dias")
	}

	// 2. Detect Categories & Merchants
	if strings.Contains(q, "mercado") || strings.Contains(q, "supermercado") || strings.Contains(q, "alimentação") {
		filter.Category = "Alimentação & Mercado"
		explanationParts = append(explanationParts, "categoria: Alimentação & Mercado")
	} else if strings.Contains(q, "transporte") || strings.Contains(q, "posto") || strings.Contains(q, "combustível") || strings.Contains(q, "uber") {
		filter.Category = "Transporte & Mobilidade"
		explanationParts = append(explanationParts, "categoria: Transporte & Mobilidade")
	} else if strings.Contains(q, "luz") || strings.Contains(q, "energia") || strings.Contains(q, "cpfl") {
		filter.Merchant = "CPFL"
		explanationParts = append(explanationParts, "estabelecimento/serviço: CPFL")
	} else if strings.Contains(q, "farmácia") || strings.Contains(q, "farmacia") || strings.Contains(q, "saúde") {
		filter.Category = "Saúde & Farmácia"
		explanationParts = append(explanationParts, "categoria: Saúde & Farmácia")
	} else if strings.Contains(q, "lazer") || strings.Contains(q, "restaurante") || strings.Contains(q, "ifood") {
		filter.Category = "Lazer & Restaurantes"
		explanationParts = append(explanationParts, "categoria: Lazer & Restaurantes")
	} else if strings.Contains(q, "educação") || strings.Contains(q, "escola") || strings.Contains(q, "faculdade") {
		filter.Category = "Educação & Cursos"
		explanationParts = append(explanationParts, "categoria: Educação & Cursos")
	}

	// 3. Detect Amount Thresholds (e.g. "maiores que 500", "acima de R$ 200", "menores que 50")
	greaterRegex := regexp.MustCompile(`(?:maior(?:es)?\s+que|acima\s+de|mais\s+de)\s*(?:r\$\s*)?(\d+(?:[.,]\d{2})?)`)
	if matches := greaterRegex.FindStringSubmatch(q); len(matches) > 1 {
		valStr := strings.ReplaceAll(matches[1], ",", ".")
		if val, err := strconv.ParseFloat(valStr, 64); err == nil {
			filter.MinAmountMinor = int64(val * 100)
			explanationParts = append(explanationParts, fmt.Sprintf("valor mínimo: R$ %.2f", val))
		}
	}

	lessRegex := regexp.MustCompile(`(?:menor(?:es)?\s+que|abaixo\s+de|menos\s+de)\s*(?:r\$\s*)?(\d+(?:[.,]\d{2})?)`)
	if matches := lessRegex.FindStringSubmatch(q); len(matches) > 1 {
		valStr := strings.ReplaceAll(matches[1], ",", ".")
		if val, err := strconv.ParseFloat(valStr, 64); err == nil {
			filter.MaxAmountMinor = int64(val * 100)
			explanationParts = append(explanationParts, fmt.Sprintf("valor máximo: R$ %.2f", val))
		}
	}

	if len(explanationParts) == 0 {
		filter.Explanation = "Busca ampla em todas as transações"
	} else {
		filter.Explanation = "Filtro aplicado: " + strings.Join(explanationParts, " • ")
	}

	return filter
}

// Execute applies the safe filter over a transaction dataset.
func (n *NaturalLanguageInterpreter) Execute(filter QueryFilter, txs []TransactionRecord) QueryResult {
	var matches []TransactionRecord
	var totalMinor int64

	for _, tx := range txs {
		// Category match
		if filter.Category != "" && !strings.EqualFold(tx.Category, filter.Category) {
			continue
		}

		// Merchant match
		if filter.Merchant != "" && !strings.Contains(strings.ToLower(tx.Description), strings.ToLower(filter.Merchant)) {
			continue
		}

		// Date filters
		if filter.StartDate != nil && tx.Date.Before(*filter.StartDate) {
			continue
		}
		if filter.EndDate != nil && tx.Date.After(*filter.EndDate) {
			continue
		}

		// Amount filters
		if filter.MinAmountMinor > 0 && tx.AmountMinor < filter.MinAmountMinor {
			continue
		}
		if filter.MaxAmountMinor > 0 && tx.AmountMinor > filter.MaxAmountMinor {
			continue
		}

		matches = append(matches, tx)
		totalMinor += tx.AmountMinor
	}

	totalMoney := domain.Money{AmountMinor: totalMinor, Currency: "BRL"}
	summary := fmt.Sprintf("Encontramos %d transações totalizando %s no período pesquisado.", len(matches), totalMoney.FormatBRL())
	if len(matches) == 0 {
		summary = "Nenhuma transação encontrada para os critérios solicitados."
	}

	return QueryResult{
		Filter:       filter,
		TotalCount:   len(matches),
		TotalAmount:  totalMoney,
		SummaryText:  summary,
		Transactions: matches,
	}
}
