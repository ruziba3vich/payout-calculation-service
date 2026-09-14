package payout

import (
	"testing"

	"github.com/shopspring/decimal"
)

func TestRateFor(t *testing.T) {
	c := NewCalculator(nil)

	cases := []struct {
		count int64
		want  string
	}{
		{0, "0.1"},
		{1, "0.1"},
		{50, "0.1"},
		{51, "0.12"},
		{150, "0.12"},
		{151, "0.15"},
		{1000, "0.15"},
	}

	for _, tc := range cases {
		got := c.RateFor(tc.count)
		if got.String() != tc.want {
			t.Errorf("RateFor(%d) = %s, want %s", tc.count, got, tc.want)
		}
	}
}

func TestCalculate(t *testing.T) {
	c := NewCalculator(nil)

	got := c.Calculate(60, decimal.RequireFromString("1000000"))

	if got.CommissionRate.String() != "0.12" {
		t.Errorf("rate = %s, want 0.12", got.CommissionRate)
	}
	if got.CommissionAmount.String() != "120000" {
		t.Errorf("commission = %s, want 120000", got.CommissionAmount)
	}
	if !got.NetAmount.Equal(got.CommissionAmount) {
		t.Errorf("net = %s, want %s", got.NetAmount, got.CommissionAmount)
	}
}

func TestCalculateRounds(t *testing.T) {
	c := NewCalculator(nil)

	got := c.Calculate(1, decimal.RequireFromString("333.33"))

	if got.CommissionAmount.String() != "33.33" {
		t.Errorf("commission = %s, want 33.33", got.CommissionAmount)
	}
}

func TestParsePeriod(t *testing.T) {
	p, err := ParsePeriod("2026-08")
	if err != nil {
		t.Fatal(err)
	}
	if p.Format("2006-01-02") != "2026-08-01" {
		t.Errorf("got %s", p)
	}

	if _, err := ParsePeriod("2026-8"); err == nil {
		t.Error("expected error for bad month")
	}
	if _, err := ParsePeriod("08-2026"); err == nil {
		t.Error("expected error for bad format")
	}
}

func TestPeriodBounds(t *testing.T) {
	p, _ := ParsePeriod("2026-12")
	start, end := PeriodBounds(p)

	if start.Format("2006-01-02") != "2026-12-01" {
		t.Errorf("start = %s", start)
	}
	if end.Format("2006-01-02") != "2027-01-01" {
		t.Errorf("end = %s", end)
	}
}
