package payout

import "github.com/shopspring/decimal"

// Tier is a commission step. Rate applies when the courier delivered at
// least MinCount orders in the month.
type Tier struct {
	MinCount int64
	Rate     decimal.Decimal
}

// DefaultTiers: 0-50 -> 10%, 51-150 -> 12%, 151+ -> 15%.
// The tier is chosen by the total delivered count of the month and that
// single rate is applied to the whole month.
var DefaultTiers = []Tier{
	{MinCount: 0, Rate: decimal.RequireFromString("0.10")},
	{MinCount: 51, Rate: decimal.RequireFromString("0.12")},
	{MinCount: 151, Rate: decimal.RequireFromString("0.15")},
}

type Calculation struct {
	DeliveredCount   int64
	GrossAmount      decimal.Decimal
	CommissionRate   decimal.Decimal
	CommissionAmount decimal.Decimal
	NetAmount        decimal.Decimal
}

type Calculator struct {
	tiers []Tier
}

func NewCalculator(tiers []Tier) *Calculator {
	if len(tiers) == 0 {
		tiers = DefaultTiers
	}
	return &Calculator{tiers: tiers}
}

func (c *Calculator) RateFor(deliveredCount int64) decimal.Decimal {
	rate := c.tiers[0].Rate
	for _, t := range c.tiers {
		if deliveredCount >= t.MinCount {
			rate = t.Rate
		}
	}
	return rate
}

// Calculate computes the courier's payout for a month.
// The courier earns the commission: net = gross * rate.
func (c *Calculator) Calculate(deliveredCount int64, gross decimal.Decimal) Calculation {
	rate := c.RateFor(deliveredCount)
	commission := gross.Mul(rate).Round(2)

	return Calculation{
		DeliveredCount:   deliveredCount,
		GrossAmount:      gross,
		CommissionRate:   rate,
		CommissionAmount: commission,
		NetAmount:        commission,
	}
}
