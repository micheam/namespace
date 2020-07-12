package ns

import (
	"github.com/google/uuid"
)

type NodeID string

func NewNodeID() *NodeID {
	id := uuid.New().String()
	nodeID := NodeID(id)
	return &nodeID
}

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
}

func NewNode(name string) *Node {
	id := NewNodeID()
	return &Node{
		ID:   *id,
		Name: NodeName(name),
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
