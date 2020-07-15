package ns

import (
	"regexp"
	"time"

	"github.com/google/uuid"
)

// =============================================
// NodeID {{{1

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

// =============================================
// NodeName {{{1

// NodeName is a Name of Node
type NodeName struct {
	str   string
	valid bool
}

// NewNodeName returns a NodeName generated from s
func NewNodeName(s string) *NodeName {

	valid := func(s string) bool {
		// とりあえず、スラッシュだけ検査
		hasSlash, _ := regexp.Match(`/`, []byte(s))
		return !hasSlash
	}(s)

	return &NodeName{
		str:   s,
		valid: valid,
	}
}

// Valid returns Validation result
func (n *NodeName) Valid() bool {
	return n.valid
}

// String returns a string representation of Node.
func (n *NodeName) String() string {
	return n.str
}

// =============================================
// NodeDescription {{{1

type NodeDescription string

func (n *NodeDescription) String() string {
	return string(*n)
}

// =============================================
// Node {{{1

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
		Name:      *NewNodeName(name),
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
