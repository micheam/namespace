package ns

import (
	"github.com/google/uuid"
)

type UID uuid.UUID
type UserName string

type User struct {
	id   UID
	name UserName
}
