package courier

import "github.com/ruziba3vich/payout-calculation-service/internal/domain/errs"

var (
	ErrNotFound   = errs.NotFoundf("courier not found")
	ErrPhoneTaken = errs.Conflictf("courier phone already taken")
)
