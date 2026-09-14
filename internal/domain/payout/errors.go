package payout

import "github.com/ruziba3vich/payout-calculation-service/internal/domain/errs"

var (
	ErrNotFound      = errs.NotFoundf("payout not found")
	ErrAlreadyExists = errs.Conflictf("payout for this courier and period already exists")
	ErrInvalidPeriod = errs.Invalidf("period must be in YYYY-MM format")
	ErrFuturePeriod  = errs.Invalidf("period must not be in the future")
)
