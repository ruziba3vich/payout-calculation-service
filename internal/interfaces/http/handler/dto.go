package handler

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/ruziba3vich/payout-calculation-service/internal/domain/courier"
	"github.com/ruziba3vich/payout-calculation-service/internal/domain/order"
	"github.com/ruziba3vich/payout-calculation-service/internal/domain/payout"
)

type CourierResponse struct {
	ID        uuid.UUID `json:"id"`
	FullName  string    `json:"full_name"`
	Phone     string    `json:"phone"`
	HiredAt   string    `json:"hired_at"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func toCourierResponse(c courier.Courier) CourierResponse {
	return CourierResponse{
		ID:        c.ID,
		FullName:  c.FullName,
		Phone:     c.Phone,
		HiredAt:   c.HiredAt.Format("2006-01-02"),
		IsActive:  c.IsActive,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

func toCourierResponses(cs []courier.Courier) []CourierResponse {
	out := make([]CourierResponse, 0, len(cs))
	for _, c := range cs {
		out = append(out, toCourierResponse(c))
	}
	return out
}

type OrderResponse struct {
	ID          uuid.UUID       `json:"id"`
	CourierID   uuid.UUID       `json:"courier_id"`
	Amount      decimal.Decimal `json:"amount"`
	Status      order.Status    `json:"status"`
	DeliveredAt *time.Time      `json:"delivered_at"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

func toOrderResponse(o order.Order) OrderResponse {
	return OrderResponse{
		ID:          o.ID,
		CourierID:   o.CourierID,
		Amount:      o.Amount,
		Status:      o.Status,
		DeliveredAt: o.DeliveredAt,
		CreatedAt:   o.CreatedAt,
		UpdatedAt:   o.UpdatedAt,
	}
}

func toOrderResponses(os []order.Order) []OrderResponse {
	out := make([]OrderResponse, 0, len(os))
	for _, o := range os {
		out = append(out, toOrderResponse(o))
	}
	return out
}

type PayoutResponse struct {
	ID               uuid.UUID       `json:"id"`
	CourierID        uuid.UUID       `json:"courier_id"`
	Period           string          `json:"period"`
	DeliveredCount   int32           `json:"delivered_count"`
	GrossAmount      decimal.Decimal `json:"gross_amount"`
	CommissionRate   decimal.Decimal `json:"commission_rate"`
	CommissionAmount decimal.Decimal `json:"commission_amount"`
	NetAmount        decimal.Decimal `json:"net_amount"`
	Status           payout.Status   `json:"status"`
	CalculatedAt     time.Time       `json:"calculated_at"`
	PaidAt           *time.Time      `json:"paid_at"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
}

func toPayoutResponse(p payout.Payout) PayoutResponse {
	return PayoutResponse{
		ID:               p.ID,
		CourierID:        p.CourierID,
		Period:           p.Period.Format("2006-01"),
		DeliveredCount:   p.DeliveredCount,
		GrossAmount:      p.GrossAmount,
		CommissionRate:   p.CommissionRate,
		CommissionAmount: p.CommissionAmount,
		NetAmount:        p.NetAmount,
		Status:           p.Status,
		CalculatedAt:     p.CalculatedAt,
		PaidAt:           p.PaidAt,
		CreatedAt:        p.CreatedAt,
		UpdatedAt:        p.UpdatedAt,
	}
}

func toPayoutResponses(ps []payout.Payout) []PayoutResponse {
	out := make([]PayoutResponse, 0, len(ps))
	for _, p := range ps {
		out = append(out, toPayoutResponse(p))
	}
	return out
}

type AdjustmentResponse struct {
	ID              uuid.UUID             `json:"id"`
	PayoutID        uuid.UUID             `json:"payout_id"`
	OrderID         *uuid.UUID            `json:"order_id"`
	Type            payout.AdjustmentType `json:"type"`
	GrossDelta      decimal.Decimal       `json:"gross_delta"`
	CommissionDelta decimal.Decimal       `json:"commission_delta"`
	NetDelta        decimal.Decimal       `json:"net_delta"`
	Reason          *string               `json:"reason"`
	CreatedAt       time.Time             `json:"created_at"`
}

func toAdjustmentResponse(a payout.Adjustment) AdjustmentResponse {
	return AdjustmentResponse{
		ID:              a.ID,
		PayoutID:        a.PayoutID,
		OrderID:         a.OrderID,
		Type:            a.Type,
		GrossDelta:      a.GrossDelta,
		CommissionDelta: a.CommissionDelta,
		NetDelta:        a.NetDelta,
		Reason:          a.Reason,
		CreatedAt:       a.CreatedAt,
	}
}

func toAdjustmentResponses(as []payout.Adjustment) []AdjustmentResponse {
	out := make([]AdjustmentResponse, 0, len(as))
	for _, a := range as {
		out = append(out, toAdjustmentResponse(a))
	}
	return out
}

type PayoutWithAdjustmentsResponse struct {
	PayoutResponse
	Adjustments []AdjustmentResponse `json:"adjustments"`
}
