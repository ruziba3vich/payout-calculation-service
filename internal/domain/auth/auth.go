package auth

import "errors"

type Role string

const (
	RoleAdmin   Role = "admin"
	RoleCourier Role = "courier"
)

var ErrInvalidCredentials = errors.New("invalid credentials")
