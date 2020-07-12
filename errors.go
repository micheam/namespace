package ns

import "errors"

var (
	ErrNotFound         = errors.New("not found")
	ErrDuplicatedEntity = errors.New("duplicated")
)
