package courier

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, c Courier) (Courier, error)
	GetByID(ctx context.Context, id uuid.UUID) (Courier, error)
	GetByPhone(ctx context.Context, phone string) (Courier, error)
	Update(ctx context.Context, p UpdateParams) (Courier, error)
	SoftDelete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, p ListParams) ([]Courier, int64, error)
	Exists(ctx context.Context, id uuid.UUID) (bool, error)
}
