package order

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, o Order) (Order, error)
	GetByID(ctx context.Context, id uuid.UUID) (Order, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status Status) (Order, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, p ListParams) ([]Order, int64, error)
}
