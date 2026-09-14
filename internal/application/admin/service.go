package admin

import (
	"context"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/ruziba3vich/payout-calculation-service/internal/domain/admin"
)

type Service struct {
	repo admin.Repository
}

func NewService(repo admin.Repository) *Service {
	return &Service{repo: repo}
}

type CreateInput struct {
	FullName string
	Username string
	Password string
}

type UpdateInput struct {
	ID       uuid.UUID
	FullName string
	Username string
	Password string
}

type ListInput struct {
	Limit  int32
	Offset int32
}

func (s *Service) Create(ctx context.Context, in CreateInput) (admin.Admin, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return admin.Admin{}, err
	}

	return s.repo.Create(ctx, admin.Admin{
		ID:           uuid.New(),
		FullName:     in.FullName,
		Username:     in.Username,
		PasswordHash: string(hash),
	})
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (admin.Admin, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) GetByUsername(ctx context.Context, username string) (admin.Admin, error) {
	return s.repo.GetByUsername(ctx, username)
}

func (s *Service) Update(ctx context.Context, in UpdateInput) (admin.Admin, error) {
	var hash string
	if in.Password != "" {
		b, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
		if err != nil {
			return admin.Admin{}, err
		}
		hash = string(b)
	}

	return s.repo.Update(ctx, admin.UpdateParams{
		ID:           in.ID,
		FullName:     in.FullName,
		Username:     in.Username,
		PasswordHash: hash,
	})
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.SoftDelete(ctx, id)
}

func (s *Service) List(ctx context.Context, in ListInput) ([]admin.Admin, int64, error) {
	return s.repo.List(ctx, admin.ListParams{
		Limit:  in.Limit,
		Offset: in.Offset,
	})
}

// EnsureExists creates the admin if no admin with that username exists yet.
func (s *Service) EnsureExists(ctx context.Context, username, password string) error {
	if _, err := s.repo.GetByUsername(ctx, username); err == nil {
		return nil
	}

	_, err := s.Create(ctx, CreateInput{
		FullName: username,
		Username: username,
		Password: password,
	})
	return err
}
