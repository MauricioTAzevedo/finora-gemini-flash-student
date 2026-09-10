package domain_test

import (
	"testing"

	"github.com/MauricioTAzevedo/finora-gemini-flash-student/services/api/internal/domain"
)

func TestMoneyArithmetic(t *testing.T) {
	m1 := domain.NewBRL(10050) // R$ 100,50
	m2 := domain.NewBRL(2525)  // R$ 25,25

	sum, err := m1.Add(m2)
	if err != nil {
		t.Fatalf("unexpected error on Add: %v", err)
	}
	if sum.AmountMinor != 12575 {
		t.Errorf("expected 12575 minor units, got %d", sum.AmountMinor)
	}

	diff, err := m1.Sub(m2)
	if err != nil {
		t.Fatalf("unexpected error on Sub: %v", err)
	}
	if diff.AmountMinor != 7525 {
		t.Errorf("expected 7525 minor units, got %d", diff.AmountMinor)
	}

	neg := m1.Negate()
	if neg.AmountMinor != -10050 {
		t.Errorf("expected -10050 minor units, got %d", neg.AmountMinor)
	}
}

func TestMoneyCurrencyMismatch(t *testing.T) {
	brl := domain.NewMoney(1000, "BRL")
	usd := domain.NewMoney(1000, "USD")

	_, err := brl.Add(usd)
	if err != domain.ErrCurrencyMismatch {
		t.Errorf("expected ErrCurrencyMismatch, got %v", err)
	}
}

func TestMoneyFormatBRL(t *testing.T) {
	tests := []struct {
		cents    int64
		expected string
	}{
		{0, "R$ 0,00"},
		{50, "R$ 0,50"},
		{100, "R$ 1,00"},
		{125050, "R$ 1.250,50"},
		{100000000, "R$ 1.000.000,00"},
		{-50, "-R$ 0,50"},
		{-312000, "-R$ 3.120,00"},
	}

	for _, tc := range tests {
		m := domain.NewBRL(tc.cents)
		got := m.FormatBRL()
		if got != tc.expected {
			t.Errorf("FormatBRL(%d) = %q; want %q", tc.cents, got, tc.expected)
		}
	}
}

func TestParseBRL(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"R$ 1.250,50", 125050},
		{"1250,50", 125050},
		{"-R$ 3.120,00", -312000},
		{"-45,90", -4590},
		{"0,00", 0},
		{"100", 10000},
		{"R$ 15,5", 1550},
	}

	for _, tc := range tests {
		m, err := domain.ParseBRL(tc.input)
		if err != nil {
			t.Errorf("ParseBRL(%q) unexpected error: %v", tc.input, err)
			continue
		}
		if m.AmountMinor != tc.expected {
			t.Errorf("ParseBRL(%q) = %d; want %d", tc.input, m.AmountMinor, tc.expected)
		}
	}
}
