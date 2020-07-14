package ns

import (
	"time"

	"github.com/google/uuid"
)

// NodeID is an identity of node
type NodeID string

// NewNodeID return generated NodeID
func NewNodeID() *NodeID {
	id := uuid.New().String()
	nodeID := NodeID(id)
	return &nodeID
}

// String returns string
func (n *NodeID) String() string {
	if n == nil {
		return ""
	}
	return string(*n)
}

type NodeName string

func (n *NodeName) String() string {
	return string(*n)
}

type NodeDescription string

func (n *NodeDescription) String() string {
	return string(*n)
}

type Node struct {
	ID          NodeID
	Name        NodeName
	Description *NodeDescription
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewNode(name string) *Node {
	id := NewNodeID()
	now := time.Now()
	return &Node{
		ID:        *id,
		Name:      NodeName(name),
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (n *Node) WithDesc(d string) *Node {
	desc := NodeDescription(d)
	n.Description = &desc
	return n
}

type NodeReader interface {
	GetByID(owner User, id NodeID) (*Node, error)
}
