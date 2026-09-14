package courier

import "errors"

var (
	ErrNotFound   = errors.New("courier not found")
	ErrPhoneTaken = errors.New("courier phone already taken")
)
