package auth

import (
	"context"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/ruziba3vich/payout-calculation-service/internal/domain/admin"
	"github.com/ruziba3vich/payout-calculation-service/internal/domain/auth"
	"github.com/ruziba3vich/payout-calculation-service/internal/domain/courier"
)

type TokenIssuer interface {
	Issue(subject uuid.UUID, role auth.Role) (string, time.Time, error)
}

type Service struct {
	admins   admin.Repository
	couriers courier.Repository
	tokens   TokenIssuer
}

func NewService(admins admin.Repository, couriers courier.Repository, tokens TokenIssuer) *Service {
	return &Service{
		admins:   admins,
		couriers: couriers,
		tokens:   tokens,
	}
}

type Token struct {
	AccessToken string
	ExpiresAt   time.Time
	Role        auth.Role
	SubjectID   uuid.UUID
}

func (s *Service) LoginAdmin(ctx context.Context, username, password string) (Token, error) {
	a, err := s.admins.GetByUsername(ctx, username)
	if err != nil {
		return Token{}, auth.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(a.PasswordHash), []byte(password)); err != nil {
		return Token{}, auth.ErrInvalidCredentials
	}

	return s.issue(a.ID, auth.RoleAdmin)
}

func (s *Service) LoginCourier(ctx context.Context, phone, password string) (Token, error) {
	c, err := s.couriers.GetByPhone(ctx, phone)
	if err != nil {
		return Token{}, auth.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(c.PasswordHash), []byte(password)); err != nil {
		return Token{}, auth.ErrInvalidCredentials
	}

	return s.issue(c.ID, auth.RoleCourier)
}

func (s *Service) issue(id uuid.UUID, role auth.Role) (Token, error) {
	token, expiresAt, err := s.tokens.Issue(id, role)
	if err != nil {
		return Token{}, err
	}

	return Token{
		AccessToken: token,
		ExpiresAt:   expiresAt,
		Role:        role,
		SubjectID:   id,
	}, nil
}
