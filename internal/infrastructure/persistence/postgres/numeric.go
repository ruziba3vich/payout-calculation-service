package postgres

import (
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
)

func toDecimal(n pgtype.Numeric) decimal.Decimal {
	if !n.Valid {
		return decimal.Zero
	}
	b, err := n.MarshalJSON()
	if err != nil {
		return decimal.Zero
	}
	d, err := decimal.NewFromString(string(b))
	if err != nil {
		return decimal.Zero
	}
	return d
}

func fromDecimal(d decimal.Decimal) pgtype.Numeric {
	var n pgtype.Numeric
	_ = n.Scan(d.String())
	return n
}

func fromDecimalPtr(d *decimal.Decimal) pgtype.Numeric {
	if d == nil {
		return pgtype.Numeric{}
	}
	return fromDecimal(*d)
}
