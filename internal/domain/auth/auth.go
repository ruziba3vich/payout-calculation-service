package auth

import "github.com/ruziba3vich/payout-calculation-service/internal/domain/errs"

type Role string

const (
	RoleAdmin   Role = "admin"
	RoleCourier Role = "courier"
)

var ErrInvalidCredentials = errs.Unauthorizedf("invalid credentials")
