package apperror

import "errors"

var (
	ErrEmailAlreadyExists = errors.New("Email already exists")
	ErrInvalidEmailFormat = errors.New("Email must be a valid email address")
)
