package apperror

import "errors"

var (
	ErrSessionExpired = errors.New("Session expired, please login again")
	ErrInvalidSession = errors.New("Invalid session, please login again")
	ErrInternal       = errors.New("Internal server error")
	ErrLogoutFailed   = errors.New("Failed to logout")
)
