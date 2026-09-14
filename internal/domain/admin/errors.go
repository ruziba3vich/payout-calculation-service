package admin

import "github.com/ruziba3vich/payout-calculation-service/internal/domain/errs"

var (
	ErrNotFound      = errs.NotFoundf("admin not found")
	ErrUsernameTaken = errs.Conflictf("admin username already taken")
)
