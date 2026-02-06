package apperror

import "errors"

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidEmailFormat = errors.New("email must be a valid email address")
)
