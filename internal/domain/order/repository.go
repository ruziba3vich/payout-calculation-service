package order

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Repository interface {
	Create(ctx context.Context, o Order) (Order, error)
	GetByID(ctx context.Context, id uuid.UUID) (Order, error)
	GetByIDForUpdate(ctx context.Context, id uuid.UUID) (Order, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status Status) (Order, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, p ListParams) ([]Order, int64, error)
	// DeliveredStats returns count and sum of delivered orders in [start, end).
	DeliveredStats(ctx context.Context, courierID uuid.UUID, start, end time.Time) (int64, decimal.Decimal, error)
}
