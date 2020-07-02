package ns

import (
	"github.com/google/uuid"
)

type ID uuid.UUID
type Name string
type Description string

type Node struct {
	id          ID
	name        Name
	description Description
	children    []*Node
}

type NodeReader interface {
	GetByID(root *Node, id ID) *Node
	FindByName(root *Node, name Name) []*Node
}
