package payout

import "errors"

var (
	ErrNotFound      = errors.New("payout not found")
	ErrAlreadyExists = errors.New("payout for this courier and period already exists")
	ErrInvalidPeriod = errors.New("period must be in YYYY-MM format")
	ErrFuturePeriod  = errors.New("period must not be in the future")
)
