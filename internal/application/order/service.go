package order

import (
	"context"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/ruziba3vich/payout-calculation-service/internal/domain/courier"
	"github.com/ruziba3vich/payout-calculation-service/internal/domain/errs"
	"github.com/ruziba3vich/payout-calculation-service/internal/domain/order"
	"github.com/ruziba3vich/payout-calculation-service/internal/domain/payout"
)

// Reconciler is implemented by the payout service. It is called inside the
// status-change transaction so the adjustment and the status update commit together.
type Reconciler interface {
	Reconcile(ctx context.Context, r payout.Repos, o order.Order, typ payout.AdjustmentType) error
}

type Service struct {
	repo       order.Repository
	couriers   courier.Repository
	tx         payout.Transactor
	reconciler Reconciler
}

func NewService(repo order.Repository, couriers courier.Repository, tx payout.Transactor, reconciler Reconciler) *Service {
	return &Service{
		repo:       repo,
		couriers:   couriers,
		tx:         tx,
		reconciler: reconciler,
	}
}

type CreateInput struct {
	CourierID uuid.UUID
	Amount    decimal.Decimal
}

func (s *Service) Create(ctx context.Context, in CreateInput) (order.Order, error) {
	if _, err := s.couriers.GetByID(ctx, in.CourierID); err != nil {
		return order.Order{}, errs.Wrap(err, "order service: create")
	}

	return s.repo.Create(ctx, order.Order{
		ID:        uuid.New(),
		CourierID: in.CourierID,
		Amount:    in.Amount,
		Status:    order.StatusPending,
	})
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (order.Order, error) {
	return s.repo.GetByID(ctx, id)
}

// UpdateStatus changes the status and, when the order's month already has a
// payout, writes an adjustment in the same transaction.
func (s *Service) UpdateStatus(ctx context.Context, id uuid.UUID, status order.Status) (order.Order, error) {
	var updated order.Order

	err := s.tx.InTx(ctx, func(r payout.Repos) error {
		// lock the row so two parallel status changes can't both reconcile
		if _, err := r.Orders.GetByIDForUpdate(ctx, id); err != nil {
			return err
		}

		o, err := r.Orders.UpdateStatus(ctx, id, status)
		if err != nil {
			return err
		}
		updated = o

		return s.reconciler.Reconcile(ctx, r, o, adjustmentTypeFor(status))
	})
	if err != nil {
		return order.Order{}, errs.Wrap(err, "order service: update status")
	}

	return updated, nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *Service) List(ctx context.Context, in order.ListParams) ([]order.Order, int64, error) {
	return s.repo.List(ctx, in)
}

func adjustmentTypeFor(status order.Status) payout.AdjustmentType {
	switch status {
	case order.StatusDelivered:
		return payout.AdjustmentOrderAdded
	case order.StatusReturned:
		return payout.AdjustmentOrderReturned
	case order.StatusCancelled:
		return payout.AdjustmentOrderCancelled
	default:
		return payout.AdjustmentManual
	}
}
