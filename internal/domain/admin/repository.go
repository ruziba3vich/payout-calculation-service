package admin

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, a Admin) (Admin, error)
	GetByID(ctx context.Context, id uuid.UUID) (Admin, error)
	GetByUsername(ctx context.Context, username string) (Admin, error)
	Update(ctx context.Context, p UpdateParams) (Admin, error)
	SoftDelete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, p ListParams) ([]Admin, int64, error)
}
