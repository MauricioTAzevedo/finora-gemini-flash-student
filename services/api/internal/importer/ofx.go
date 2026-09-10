package importer

import (
	"bufio"
	"errors"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/MauricioTAzevedo/finora-gemini-flash-student/services/api/internal/domain"
)

var (
	ErrInvalidOFX = errors.New("invalid or empty OFX statement")
)

var (
	reOFXTag = regexp.MustCompile(`<([A-Za-z0-9_]+)>([^<\r\n]*)`)
)

// ParseOFX parses standard Brazilian banking OFX (Open Financial Exchange) files.
func ParseOFX(r io.Reader) ([]CanonicalIngestionRecord, error) {
	scanner := bufio.NewScanner(r)
	var records []CanonicalIngestionRecord

	inTransaction := false
	var currentTrn map[string]string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if strings.Contains(line, "<STMTTRN>") {
			inTransaction = true
			currentTrn = make(map[string]string)
			continue
		}

		if strings.Contains(line, "</STMTTRN>") {
			if inTransaction && currentTrn != nil {
				record, err := buildCanonicalFromOFX(currentTrn)
				if err == nil {
					records = append(records, record)
				}
			}
			inTransaction = false
			currentTrn = nil
			continue
		}

		if inTransaction && currentTrn != nil {
			matches := reOFXTag.FindAllStringSubmatch(line, -1)
			for _, m := range matches {
				if len(m) >= 3 {
					tag := strings.ToUpper(strings.TrimSpace(m[1]))
					val := strings.TrimSpace(m[2])
					currentTrn[tag] = val
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return nil, ErrInvalidOFX
	}

	return records, nil
}

func buildCanonicalFromOFX(tags map[string]string) (CanonicalIngestionRecord, error) {
	// Parse amount
	amtStr := tags["TRNAMT"]
	if amtStr == "" {
		return CanonicalIngestionRecord{}, errors.New("missing TRNAMT tag")
	}

	// Example: -129.90 or 8500.00
	cleanAmt := strings.ReplaceAll(amtStr, ",", ".")
	valFloat, err := strconv.ParseFloat(cleanAmt, 64)
	if err != nil {
		return CanonicalIngestionRecord{}, err
	}

	// Minor units (cents)
	var cents int64
	if valFloat < 0 {
		cents = int64((valFloat * 100) - 0.5)
	} else {
		cents = int64((valFloat * 100) + 0.5)
	}

	// Parse date (OFX YYYYMMDDHHMMSS or YYYYMMDD)
	dtStr := tags["DTPOSTED"]
	var occurredAt time.Time
	if len(dtStr) >= 8 {
		datePart := dtStr[:8]
		t, err := time.Parse("20060102", datePart)
		if err == nil {
			occurredAt = t
		}
	}
	if occurredAt.IsZero() {
		occurredAt = time.Now().UTC()
	}

	// Description from MEMO or NAME
	desc := tags["MEMO"]
	if desc == "" {
		desc = tags["NAME"]
	}
	if desc == "" {
		desc = "Lançamento OFX"
	}

	fitID := tags["FITID"]
	if fitID == "" {
		fitID = tags["CHECKNUM"]
	}

	return CanonicalIngestionRecord{
		ExternalID:  fitID,
		Source:      "ofx",
		OccurredAt:  occurredAt,
		Description: desc,
		Amount:      domain.NewBRL(cents),
	}, nil
}
