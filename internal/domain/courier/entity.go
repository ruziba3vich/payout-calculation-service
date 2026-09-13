package courier

import (
	"time"

	"github.com/google/uuid"
)

type Courier struct {
	ID           uuid.UUID
	FullName     string
	Phone        string
	PasswordHash string
	HiredAt      time.Time
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type UpdateParams struct {
	ID           uuid.UUID
	FullName     *string
	Phone        *string
	PasswordHash *string
	HiredAt      *time.Time
	IsActive     *bool
}

type ListParams struct {
	Limit  int32
	Offset int32
}
