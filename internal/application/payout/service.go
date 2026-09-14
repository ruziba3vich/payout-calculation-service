package payout

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/ruziba3vich/payout-calculation-service/internal/domain/courier"
	"github.com/ruziba3vich/payout-calculation-service/internal/domain/errs"
	"github.com/ruziba3vich/payout-calculation-service/internal/domain/order"
	"github.com/ruziba3vich/payout-calculation-service/internal/domain/payout"
)

type Service struct {
	payouts     payout.Repository
	adjustments payout.AdjustmentRepository
	couriers    courier.Repository
	tx          payout.Transactor
	calc        *payout.Calculator
}

func NewService(
	payouts payout.Repository,
	adjustments payout.AdjustmentRepository,
	couriers courier.Repository,
	tx payout.Transactor,
	calc *payout.Calculator,
) *Service {
	return &Service{
		payouts:     payouts,
		adjustments: adjustments,
		couriers:    couriers,
		tx:          tx,
		calc:        calc,
	}
}

type PayoutWithAdjustments struct {
	Payout      payout.Payout
	Adjustments []payout.Adjustment
}

// Calculate creates the payout for a courier and month.
// If it already exists, the existing one is returned together with ErrAlreadyExists.
func (s *Service) Calculate(ctx context.Context, courierID uuid.UUID, period time.Time) (payout.Payout, error) {
	period = payout.PeriodOf(period)
	if period.After(payout.PeriodOf(time.Now())) {
		return payout.Payout{}, payout.ErrFuturePeriod
	}

	var result payout.Payout

	err := s.tx.InTx(ctx, func(r payout.Repos) error {
		if _, err := r.Couriers.GetByID(ctx, courierID); err != nil {
			return err
		}

		start, end := payout.PeriodBounds(period)
		count, gross, err := r.Orders.DeliveredStats(ctx, courierID, start, end)
		if err != nil {
			return err
		}

		c := s.calc.Calculate(count, gross)

		result, err = r.Payouts.Create(ctx, payout.Payout{
			ID:               uuid.New(),
			CourierID:        courierID,
			Period:           period,
			DeliveredCount:   int32(c.DeliveredCount),
			GrossAmount:      c.GrossAmount,
			CommissionRate:   c.CommissionRate,
			CommissionAmount: c.CommissionAmount,
			NetAmount:        c.NetAmount,
			Status:           payout.StatusCalculated,
		})
		return err
	})

	if errors.Is(err, payout.ErrAlreadyExists) {
		existing, getErr := s.payouts.GetByCourierPeriod(ctx, courierID, period)
		if getErr != nil {
			return payout.Payout{}, errs.Wrap(getErr, "payout service: calculate")
		}
		return existing, payout.ErrAlreadyExists
	}
	if err != nil {
		return payout.Payout{}, errs.Wrap(err, "payout service: calculate")
	}

	return result, nil
}

// Reconcile is called inside the order status transaction. If the order's
// delivery month already has a payout, the month is recalculated and the
// difference is stored as an adjustment. The payout row itself is never edited.
func (s *Service) Reconcile(ctx context.Context, r payout.Repos, o order.Order, typ payout.AdjustmentType) error {
	if o.DeliveredAt == nil {
		return nil
	}

	period := payout.PeriodOf(*o.DeliveredAt)

	p, err := r.Payouts.GetByCourierPeriodForUpdate(ctx, o.CourierID, period)
	if errors.Is(err, payout.ErrNotFound) {
		return nil
	}
	if err != nil {
		return err
	}

	start, end := payout.PeriodBounds(period)
	count, gross, err := r.Orders.DeliveredStats(ctx, o.CourierID, start, end)
	if err != nil {
		return err
	}

	fresh := s.calc.Calculate(count, gross)

	applied, err := r.Adjustments.SumByPayoutID(ctx, p.ID)
	if err != nil {
		return err
	}

	grossDelta := fresh.GrossAmount.Sub(p.GrossAmount.Add(applied.GrossDelta))
	commissionDelta := fresh.CommissionAmount.Sub(p.CommissionAmount.Add(applied.CommissionDelta))
	netDelta := fresh.NetAmount.Sub(p.NetAmount.Add(applied.NetDelta))

	if grossDelta.IsZero() && commissionDelta.IsZero() && netDelta.IsZero() {
		return nil
	}

	reason := "order " + o.ID.String() + " changed to " + string(o.Status) +
		", month recalculated at rate " + fresh.CommissionRate.String()

	orderID := o.ID
	_, err = r.Adjustments.Create(ctx, payout.Adjustment{
		ID:              uuid.New(),
		PayoutID:        p.ID,
		OrderID:         &orderID,
		Type:            typ,
		GrossDelta:      grossDelta,
		CommissionDelta: commissionDelta,
		NetDelta:        netDelta,
		Reason:          &reason,
	})
	return err
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (payout.Payout, error) {
	return s.payouts.GetByID(ctx, id)
}

func (s *Service) GetWithAdjustments(ctx context.Context, id uuid.UUID) (PayoutWithAdjustments, error) {
	p, err := s.payouts.GetByID(ctx, id)
	if err != nil {
		return PayoutWithAdjustments{}, errs.Wrap(err, "payout service: get")
	}

	adjs, err := s.adjustments.ListByPayoutID(ctx, id)
	if err != nil {
		return PayoutWithAdjustments{}, errs.Wrap(err, "payout service: get adjustments")
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

func (s *Service) ListAdjustments(ctx context.Context, in payout.AdjustmentListParams) ([]payout.Adjustment, int64, error) {
	return s.adjustments.List(ctx, in)
}

// Total returns the effective payout after adjustments.
func Total(p payout.Payout, adjs []payout.Adjustment) decimal.Decimal {
	total := p.NetAmount
	for _, a := range adjs {
		total = total.Add(a.NetDelta)
	}
	return total
}
