package apperror

import "errors"

var (
	ErrNoFieldsToUpdate     = errors.New("No fields to update")
	ErrCannotFollowYourself = errors.New("Cannot follow yourself")
)
