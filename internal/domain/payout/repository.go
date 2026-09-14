package payout

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, p Payout) (Payout, error)
	GetByID(ctx context.Context, id uuid.UUID) (Payout, error)
	GetByCourierPeriod(ctx context.Context, courierID uuid.UUID, period time.Time) (Payout, error)
	GetByCourierPeriodForUpdate(ctx context.Context, courierID uuid.UUID, period time.Time) (Payout, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status Status) (Payout, error)
	List(ctx context.Context, p ListParams) ([]Payout, int64, error)
}

type AdjustmentRepository interface {
	Create(ctx context.Context, a Adjustment) (Adjustment, error)
	GetByID(ctx context.Context, id uuid.UUID) (Adjustment, error)
	ListByPayoutID(ctx context.Context, payoutID uuid.UUID) ([]Adjustment, error)
	List(ctx context.Context, p AdjustmentListParams) ([]Adjustment, int64, error)
	SumByPayoutID(ctx context.Context, payoutID uuid.UUID) (AdjustmentSum, error)
}
