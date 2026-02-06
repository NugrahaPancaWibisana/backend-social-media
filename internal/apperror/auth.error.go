package apperror

import "errors"

var (
	ErrEmailAlreadyExists = errors.New("Email already exists")
	ErrInvalidEmailFormat = errors.New("Email must be a valid email address")
	ErrUserNotFound       = errors.New("User not found")
	ErrInvalidCredential  = errors.New("Invalid email or password")
)
