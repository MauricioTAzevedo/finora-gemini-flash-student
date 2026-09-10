package domain

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

var (
	ErrCurrencyMismatch = errors.New("cannot operate on monies with different currencies")
	ErrInvalidAmount    = errors.New("invalid monetary string representation")
)

// Money represents a monetary value in integer minor units (cents) to avoid floating point drift.
// For BRL: 12345 minor units = R$ 123,45.
type Money struct {
	AmountMinor int64  `json:"amountMinor"`
	Currency    string `json:"currency"`
}

// NewMoney creates a Money instance with integer minor units.
func NewMoney(amountMinor int64, currency string) Money {
	if currency == "" {
		currency = "BRL"
	}
	return Money{
		AmountMinor: amountMinor,
		Currency:    strings.ToUpper(currency),
	}
}

// NewBRL is a convenience constructor for Brazilian Real amounts in cents.
func NewBRL(cents int64) Money {
	return NewMoney(cents, "BRL")
}

// Add sums two Money values of the same currency.
func (m Money) Add(other Money) (Money, error) {
	if m.Currency != other.Currency {
		return Money{}, ErrCurrencyMismatch
	}
	return Money{
		AmountMinor: m.AmountMinor + other.AmountMinor,
		Currency:    m.Currency,
	}, nil
}

// Sub subtracts other from m with currency validation.
func (m Money) Sub(other Money) (Money, error) {
	if m.Currency != other.Currency {
		return Money{}, ErrCurrencyMismatch
	}
	return Money{
		AmountMinor: m.AmountMinor - other.AmountMinor,
		Currency:    m.Currency,
	}, nil
}

// Negate returns the negated Money amount.
func (m Money) Negate() Money {
	return Money{
		AmountMinor: -m.AmountMinor,
		Currency:    m.Currency,
	}
}

// IsZero returns true if amount is zero.
func (m Money) IsZero() bool {
	return m.AmountMinor == 0
}

// IsPositive returns true if amount > 0.
func (m Money) IsPositive() bool {
	return m.AmountMinor > 0
}

// IsNegative returns true if amount < 0.
func (m Money) IsNegative() bool {
	return m.AmountMinor < 0
}

// FormatBRL formats the money value in pt-BR standard notation: R$ 1.234,56 or -R$ 1.234,56.
func (m Money) FormatBRL() string {
	negative := m.AmountMinor < 0
	cents := m.AmountMinor
	if negative {
		cents = -cents
	}

	whole := cents / 100
	fraction := cents % 100

	wholeStr := strconv.FormatInt(whole, 10)
	var formattedWhole strings.Builder
	length := len(wholeStr)

	for i, ch := range wholeStr {
		if i > 0 && (length-i)%3 == 0 {
			formattedWhole.WriteRune('.')
		}
		formattedWhole.WriteRune(ch)
	}

	prefix := "R$ "
	if negative {
		prefix = "-R$ "
	}

	return fmt.Sprintf("%s%s,%02d", prefix, formattedWhole.String(), fraction)
}

// ParseBRL parses strings like "R$ 1.250,50", "-R$ 45,00", "1250,50", or "-100.50" into Money.
func ParseBRL(raw string) (Money, error) {
	cleaned := strings.TrimSpace(raw)
	if cleaned == "" {
		return Money{}, ErrInvalidAmount
	}

	isNegative := false
	if strings.HasPrefix(cleaned, "-") || strings.HasPrefix(cleaned, "-R$") || strings.Contains(cleaned, "-") {
		isNegative = true
	}

	// Remove non-digit characters except commas and dots
	var numParts strings.Builder
	for _, r := range cleaned {
		if unicode.IsDigit(r) || r == ',' || r == '.' {
			numParts.WriteRune(r)
		}
	}

	numStr := numParts.String()
	if numStr == "" {
		return Money{}, ErrInvalidAmount
	}

	var wholePart, fracPart string
	if strings.Contains(numStr, ",") {
		// Brazilian notation: 1.250,50
		parts := strings.Split(numStr, ",")
		wholePart = strings.ReplaceAll(parts[0], ".", "")
		if len(parts) > 1 {
			fracPart = parts[1]
		}
	} else if strings.Contains(numStr, ".") {
		// Standard decimal notation: 1250.50
		parts := strings.Split(numStr, ".")
		wholePart = parts[0]
		if len(parts) > 1 {
			fracPart = parts[1]
		}
	} else {
		wholePart = numStr
	}

	whole, err := strconv.ParseInt(wholePart, 10, 64)
	if err != nil && wholePart != "" {
		return Money{}, fmt.Errorf("%w: %v", ErrInvalidAmount, err)
	}

	var fraction int64
	if len(fracPart) > 0 {
		if len(fracPart) == 1 {
			fracPart += "0"
		} else if len(fracPart) > 2 {
			fracPart = fracPart[:2] // truncate beyond 2 decimals
		}
		fraction, err = strconv.ParseInt(fracPart, 10, 64)
		if err != nil {
			return Money{}, fmt.Errorf("%w: %v", ErrInvalidAmount, err)
		}
	}

	totalMinor := (whole * 100) + fraction
	if isNegative {
		totalMinor = -totalMinor
	}

	return NewBRL(totalMinor), nil
}

// String implements fmt.Stringer interface.
func (m Money) String() string {
	if m.Currency == "BRL" {
		return m.FormatBRL()
	}
	return fmt.Sprintf("%s %d.%02d", m.Currency, m.AmountMinor/100, m.AmountMinor%100)
}

// MarshalJSON marshals Money into standard JSON.
func (m Money) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		AmountMinor int64  `json:"amountMinor"`
		Currency    string `json:"currency"`
		Formatted   string `json:"formatted"`
	}{
		AmountMinor: m.AmountMinor,
		Currency:    m.Currency,
		Formatted:   m.FormatBRL(),
	})
}
