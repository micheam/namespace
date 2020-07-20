package ns

import "errors"

var (
	// ErrNotFound occurs when the specified entity does not exist.
	ErrNotFound = errors.New("not found")

	// ErrDuplicatedEntity occurs when there is a duplicate entity that is attempted to be created.
	ErrDuplicatedEntity = errors.New("duplicated")
)
