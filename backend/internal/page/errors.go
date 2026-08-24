package page

import "errors"

var (
	ErrNotFound     = errors.New("not found")
	ErrConflict     = errors.New("conflict")
	ErrValidation   = errors.New("validation")
	ErrUnauthorized = errors.New("unauthorized")
	ErrLocked       = errors.New("locked")
)

type FieldError struct {
	Fields map[string]string
}

func (e *FieldError) Error() string { return ErrValidation.Error() }
func (e *FieldError) Unwrap() error { return ErrValidation }

func Invalid(field, msg string) error {
	if msg == "" {
		msg = "invalid"
	}
	return &FieldError{Fields: map[string]string{field: msg}}
}
