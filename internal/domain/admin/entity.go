package admin

import (
	"time"

	"github.com/google/uuid"
)

type Admin struct {
	ID           uuid.UUID
	FullName     string
	Username     string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type UpdateParams struct {
	ID           uuid.UUID
	FullName     string
	Username     string
	PasswordHash string
}

type ListParams struct {
	Limit  int32
	Offset int32
}
