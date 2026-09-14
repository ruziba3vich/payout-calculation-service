package order

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Status string

const (
	StatusPending   Status = "pending"
	StatusDelivered Status = "delivered"
	StatusCancelled Status = "cancelled"
	StatusReturned  Status = "returned"
)

type Order struct {
	ID          uuid.UUID
	CourierID   uuid.UUID
	Amount      decimal.Decimal
	Status      Status
	DeliveredAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ListParams struct {
	CourierID     *uuid.UUID
	Status        *Status
	DeliveredFrom *time.Time
	DeliveredTo   *time.Time
	CreatedFrom   *time.Time
	CreatedTo     *time.Time
	AmountMin     *decimal.Decimal
	AmountMax     *decimal.Decimal
	SortBy        string
	SortDir       string
	Limit         int32
	Offset        int32
}
