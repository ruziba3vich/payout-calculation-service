// Package errs holds the error type every layer returns.
// Repos wrap driver errors, services wrap repo errors, handlers map Kind to a status.
package errs

import "errors"

type Kind int

const (
	Internal Kind = iota
	NotFound
	Conflict
	Invalid
	Unauthorized
	Forbidden
)

type Error struct {
	Kind    Kind
	Message string
	Cause   error
}

func (e *Error) Error() string {
	if e.Cause != nil {
		return e.Message + ": " + e.Cause.Error()
	}
	return e.Message
}

func (e *Error) Unwrap() error {
	return e.Cause
}

func New(kind Kind, msg string) *Error {
	return &Error{Kind: kind, Message: msg}
}

func NotFoundf(msg string) *Error     { return New(NotFound, msg) }
func Conflictf(msg string) *Error     { return New(Conflict, msg) }
func Invalidf(msg string) *Error      { return New(Invalid, msg) }
func Unauthorizedf(msg string) *Error { return New(Unauthorized, msg) }
func Forbiddenf(msg string) *Error    { return New(Forbidden, msg) }

// Wrap adds context and keeps the kind of the wrapped error.
// Errors without a kind (driver, bcrypt, jwt...) become Internal.
func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}
	return &Error{Kind: KindOf(err), Message: msg, Cause: err}
}

func KindOf(err error) Kind {
	var e *Error
	if errors.As(err, &e) {
		return e.Kind
	}
	return Internal
}

// Public returns the message of the innermost typed error, the one meant
// for the client. Wrapping context like "courier repo: create" is not included.
func Public(err error) string {
	var last *Error
	for err != nil {
		var e *Error
		if !errors.As(err, &e) {
			break
		}
		last = e
		err = e.Cause
	}
	if last == nil {
		return "internal error"
	}
	return last.Message
}
