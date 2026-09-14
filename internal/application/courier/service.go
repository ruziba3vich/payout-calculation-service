package courier

import (
	"context"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/ruziba3vich/payout-calculation-service/internal/domain/courier"
	"github.com/ruziba3vich/payout-calculation-service/internal/domain/errs"
)

type Service struct {
	repo courier.Repository
}

func NewService(repo courier.Repository) *Service {
	return &Service{repo: repo}
}

type CreateInput struct {
	FullName string
	Phone    string
	Password string
	HiredAt  time.Time
}

type UpdateInput struct {
	ID       uuid.UUID
	FullName *string
	Phone    *string
	Password *string
	HiredAt  *time.Time
	IsActive *bool
}

type ListInput struct {
	Limit  int32
	Offset int32
}

func (s *Service) Create(ctx context.Context, in CreateInput) (courier.Courier, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return courier.Courier{}, errs.Wrap(err, "courier service: hash password")
	}

	return s.repo.Create(ctx, courier.Courier{
		ID:           uuid.New(),
		FullName:     in.FullName,
		Phone:        in.Phone,
		PasswordHash: string(hash),
		HiredAt:      in.HiredAt,
		IsActive:     true,
	})
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (courier.Courier, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) GetByPhone(ctx context.Context, phone string) (courier.Courier, error) {
	return s.repo.GetByPhone(ctx, phone)
}

func (s *Service) Update(ctx context.Context, in UpdateInput) (courier.Courier, error) {
	var hash *string
	if in.Password != nil {
		b, err := bcrypt.GenerateFromPassword([]byte(*in.Password), bcrypt.DefaultCost)
		if err != nil {
			return courier.Courier{}, errs.Wrap(err, "courier service: hash password")
		}
		h := string(b)
		hash = &h
	}

	return s.repo.Update(ctx, courier.UpdateParams{
		ID:           in.ID,
		FullName:     in.FullName,
		Phone:        in.Phone,
		PasswordHash: hash,
		HiredAt:      in.HiredAt,
		IsActive:     in.IsActive,
	})
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.SoftDelete(ctx, id)
}

func (s *Service) List(ctx context.Context, in ListInput) ([]courier.Courier, int64, error) {
	return s.repo.List(ctx, courier.ListParams{
		Limit:  in.Limit,
		Offset: in.Offset,
	})
}

func (s *Service) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	return s.repo.Exists(ctx, id)
}
