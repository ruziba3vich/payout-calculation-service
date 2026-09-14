package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/ruziba3vich/payout-calculation-service/internal/config"
	"github.com/ruziba3vich/payout-calculation-service/internal/domain/auth"
)

type Claims struct {
	jwt.RegisteredClaims
	Role auth.Role `json:"role"`
}

type Manager struct {
	secret []byte
	ttl    time.Duration
	issuer string
}

func NewManager(cfg config.JWTConfig) *Manager {
	return &Manager{
		secret: []byte(cfg.Secret),
		ttl:    cfg.AccessTTL,
		issuer: cfg.Issuer,
	}
}

func (m *Manager) Issue(subject uuid.UUID, role auth.Role) (string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(m.ttl)

	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.issuer,
			Subject:   subject.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
		Role: role,
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(m.secret)
	if err != nil {
		return "", time.Time{}, err
	}
	return token, expiresAt, nil
}

func (m *Manager) Parse(token string) (uuid.UUID, auth.Role, error) {
	var claims Claims

	parsed, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return m.secret, nil
	}, jwt.WithIssuer(m.issuer))
	if err != nil {
		return uuid.Nil, "", err
	}
	if !parsed.Valid {
		return uuid.Nil, "", errors.New("invalid token")
	}

	id, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, "", err
	}

	switch claims.Role {
	case auth.RoleAdmin, auth.RoleCourier:
	default:
		return uuid.Nil, "", errors.New("unknown role")
	}

	return id, claims.Role, nil
}
