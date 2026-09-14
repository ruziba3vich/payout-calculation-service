package payout

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/ruziba3vich/payout-calculation-service/internal/domain/payout"
)

type Service struct {
	payouts     payout.Repository
	adjustments payout.AdjustmentRepository
}

func NewService(payouts payout.Repository, adjustments payout.AdjustmentRepository) *Service {
	return &Service{
		payouts:     payouts,
		adjustments: adjustments,
	}
}

type CreateInput struct {
	CourierID        uuid.UUID
	Period           time.Time
	DeliveredCount   int32
	GrossAmount      decimal.Decimal
	CommissionRate   decimal.Decimal
	CommissionAmount decimal.Decimal
	NetAmount        decimal.Decimal
}

type CreateAdjustmentInput struct {
	PayoutID        uuid.UUID
	OrderID         *uuid.UUID
	Type            payout.AdjustmentType
	GrossDelta      decimal.Decimal
	CommissionDelta decimal.Decimal
	NetDelta        decimal.Decimal
	Reason          *string
}

type PayoutWithAdjustments struct {
	Payout      payout.Payout
	Adjustments []payout.Adjustment
}

func (s *Service) Create(ctx context.Context, in CreateInput) (payout.Payout, error) {
	return s.payouts.Create(ctx, payout.Payout{
		ID:               uuid.New(),
		CourierID:        in.CourierID,
		Period:           in.Period,
		DeliveredCount:   in.DeliveredCount,
		GrossAmount:      in.GrossAmount,
		CommissionRate:   in.CommissionRate,
		CommissionAmount: in.CommissionAmount,
		NetAmount:        in.NetAmount,
		Status:           payout.StatusCalculated,
	})
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (payout.Payout, error) {
	return s.payouts.GetByID(ctx, id)
}

func (s *Service) GetWithAdjustments(ctx context.Context, id uuid.UUID) (PayoutWithAdjustments, error) {
	p, err := s.payouts.GetByID(ctx, id)
	if err != nil {
		return PayoutWithAdjustments{}, err
	}

	adjs, err := s.adjustments.ListByPayoutID(ctx, id)
	if err != nil {
		return PayoutWithAdjustments{}, err
	}

	return PayoutWithAdjustments{Payout: p, Adjustments: adjs}, nil
}

func (s *Service) GetByCourierPeriod(ctx context.Context, courierID uuid.UUID, period time.Time) (payout.Payout, error) {
	return s.payouts.GetByCourierPeriod(ctx, courierID, period)
}

func (s *Service) UpdateStatus(ctx context.Context, id uuid.UUID, status payout.Status) (payout.Payout, error) {
	return s.payouts.UpdateStatus(ctx, id, status)
}

func (s *Service) List(ctx context.Context, in payout.ListParams) ([]payout.Payout, int64, error) {
	return s.payouts.List(ctx, in)
}

func (s *Service) CreateAdjustment(ctx context.Context, in CreateAdjustmentInput) (payout.Adjustment, error) {
	return s.adjustments.Create(ctx, payout.Adjustment{
		ID:              uuid.New(),
		PayoutID:        in.PayoutID,
		OrderID:         in.OrderID,
		Type:            in.Type,
		GrossDelta:      in.GrossDelta,
		CommissionDelta: in.CommissionDelta,
		NetDelta:        in.NetDelta,
		Reason:          in.Reason,
	})
}

func (s *Service) GetAdjustmentByID(ctx context.Context, id uuid.UUID) (payout.Adjustment, error) {
	return s.adjustments.GetByID(ctx, id)
}

func (s *Service) ListAdjustments(ctx context.Context, in payout.AdjustmentListParams) ([]payout.Adjustment, int64, error) {
	return s.adjustments.List(ctx, in)
}
