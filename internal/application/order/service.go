package order

import (
	"context"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/ruziba3vich/payout-calculation-service/internal/domain/order"
)

type Service struct {
	repo order.Repository
}

func NewService(repo order.Repository) *Service {
	return &Service{repo: repo}
}

type CreateInput struct {
	CourierID uuid.UUID
	Amount    decimal.Decimal
}

func (s *Service) Create(ctx context.Context, in CreateInput) (order.Order, error) {
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

func (s *Service) UpdateStatus(ctx context.Context, id uuid.UUID, status order.Status) (order.Order, error) {
	return s.repo.UpdateStatus(ctx, id, status)
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *Service) List(ctx context.Context, in order.ListParams) ([]order.Order, int64, error) {
	return s.repo.List(ctx, in)
}
