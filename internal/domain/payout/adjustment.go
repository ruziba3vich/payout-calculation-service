package payout

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type AdjustmentType string

const (
	AdjustmentOrderReturned  AdjustmentType = "order_returned"
	AdjustmentOrderCancelled AdjustmentType = "order_cancelled"
	AdjustmentOrderAdded     AdjustmentType = "order_added"
	AdjustmentManual         AdjustmentType = "manual"
)

type Adjustment struct {
	ID              uuid.UUID
	PayoutID        uuid.UUID
	OrderID         *uuid.UUID
	Type            AdjustmentType
	GrossDelta      decimal.Decimal
	CommissionDelta decimal.Decimal
	NetDelta        decimal.Decimal
	Reason          *string
	CreatedAt       time.Time
}

type AdjustmentListParams struct {
	PayoutID    *uuid.UUID
	OrderID     *uuid.UUID
	Type        *AdjustmentType
	CreatedFrom *time.Time
	CreatedTo   *time.Time
	SortBy      string
	SortDir     string
	Limit       int32
	Offset      int32
}

type AdjustmentSum struct {
	GrossDelta      decimal.Decimal
	CommissionDelta decimal.Decimal
	NetDelta        decimal.Decimal
}
