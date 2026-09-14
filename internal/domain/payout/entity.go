package payout

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Status string

const (
	StatusCalculated Status = "calculated"
	StatusPaid       Status = "paid"
	StatusCancelled  Status = "cancelled"
)

type Payout struct {
	ID               uuid.UUID
	CourierID        uuid.UUID
	Period           time.Time
	DeliveredCount   int32
	GrossAmount      decimal.Decimal
	CommissionRate   decimal.Decimal
	CommissionAmount decimal.Decimal
	NetAmount        decimal.Decimal
	Status           Status
	CalculatedAt     time.Time
	PaidAt           *time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type ListParams struct {
	CourierID  *uuid.UUID
	Status     *Status
	PeriodFrom *time.Time
	PeriodTo   *time.Time
	SortBy     string
	SortDir    string
	Limit      int32
	Offset     int32
}
