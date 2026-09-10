package importer

import (
	"encoding/csv"
	"errors"
	"io"
	"strings"
	"time"

	"github.com/MauricioTAzevedo/finora-gemini-flash-student/services/api/internal/domain"
)

// ParseCSV parses Brazilian financial CSV exports with auto-detection of delimiters (comma or semicolon).
func ParseCSV(r io.Reader) ([]CanonicalIngestionRecord, error) {
	content, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	contentStr := string(content)
	lines := strings.Split(contentStr, "\n")
	if len(lines) == 0 {
		return nil, errors.New("empty CSV file")
	}

	firstLine := lines[0]
	delimiter := ','
	if strings.Count(firstLine, ";") > strings.Count(firstLine, ",") {
		delimiter = ';'
	}

	reader := csv.NewReader(strings.NewReader(contentStr))
	reader.Comma = delimiter
	reader.LazyQuotes = true
	reader.TrimLeadingSpace = true
	reader.FieldsPerRecord = -1

	headers, err := reader.Read()
	if err != nil {
		return nil, err
	}

	headerMap := make(map[string]int)
	for i, col := range headers {
		cleaned := strings.ToLower(strings.TrimSpace(col))
		headerMap[cleaned] = i
	}

	// Identify relevant columns
	dateIdx := findColumnIndex(headerMap, []string{"data", "dt.", "dia", "date"})
	descIdx := findColumnIndex(headerMap, []string{"descrição", "descricao", "historico", "histórico", "estabelecimento", "description", "memo"})
	amtIdx := findColumnIndex(headerMap, []string{"valor", "valor pg.", "quantia", "total", "amount"})
	catIdx := findColumnIndex(headerMap, []string{"categoria", "category", "tipo"})

	if dateIdx == -1 || descIdx == -1 || amtIdx == -1 {
		return nil, errors.New("CSV missing required columns (expected date, description, and amount headers)")
	}

	var records []CanonicalIngestionRecord

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue // skip malformed row
		}

		if len(row) <= dateIdx || len(row) <= descIdx || len(row) <= amtIdx {
			continue
		}

		// Parse date (DD/MM/YYYY or YYYY-MM-DD)
		dateStr := strings.TrimSpace(row[dateIdx])
		occurredAt, err := parseDate(dateStr)
		if err != nil {
			continue
		}

		// Parse amount via domain.ParseBRL
		amtStr := strings.TrimSpace(row[amtIdx])
		money, err := domain.ParseBRL(amtStr)
		if err != nil {
			continue
		}

		desc := strings.TrimSpace(row[descIdx])
		catHint := ""
		if catIdx != -1 && len(row) > catIdx {
			catHint = strings.TrimSpace(row[catIdx])
		}

		records = append(records, CanonicalIngestionRecord{
			Source:       "csv",
			OccurredAt:   occurredAt,
			Description:  desc,
			Amount:       money,
			CategoryHint: catHint,
		})
	}

	return records, nil
}

func findColumnIndex(headers map[string]int, candidates []string) int {
	for _, c := range candidates {
		for h, idx := range headers {
			if strings.Contains(h, c) {
				return idx
			}
		}
	}
	return -1
}

func parseDate(raw string) (time.Time, error) {
	formats := []string{
		"02/01/2006",
		"02/01/06",
		"2006-01-02",
		"2006/01/02",
	}

	for _, f := range formats {
		if t, err := time.Parse(f, raw); err == nil {
			return t, nil
		}
	}
	return time.Time{}, errors.New("unsupported date format")
}
