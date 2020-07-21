package ns

import (
	"context"
	"fmt"
	"regexp"
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

// NodeDescription ...
type NodeDescription string

// String ...
func (n *NodeDescription) String() string {
	return string(*n)
}

// Node ...
type Node struct {
	ID          NodeID
	Name        NodeName
	Description *NodeDescription
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewNode ...
func NewNode(name NodeName) *Node {
	id := NewNodeID()
	now := time.Now()
	return &Node{
		ID:        *id,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// WithDesc ...
func (n *Node) WithDesc(d string) *Node {
	desc := NodeDescription(d)
	n.Description = &desc
	return n
}

// UseCases ------------------------------------------------

// NodeReader ...
type NodeReader interface {
	GetByID(owner *User, id NodeID) (*Node, error)
}

// NodeWriter ...
type NodeWriter interface {
	Save(owner *User, node *Node) error
}

// UseCases ------------------------------------------------

// NodeCreation is a UseCase.
type NodeCreation struct {
	nodeWriter NodeWriter
	presenter  NodeCreationResponseOutput
}

// NewNodeCreation return NodeCreation interactor.
func NewNodeCreation(w NodeWriter, p NodeCreationResponseOutput) *NodeCreation {
	return &NodeCreation{
		nodeWriter: w,
		presenter:  p,
	}
}

// NodeCreationRequest is a request data of new node creation.
type NodeCreationRequest struct {
	Name string
}

// NodeCreationResponse is a response data of new node creation.
type NodeCreationResponse struct {
	Created *Node
}

// NodeCreationResponseOutput defines how to output the result on new node creation.
type NodeCreationResponseOutput func(ctx context.Context, resp *NodeCreationResponse) error

// Exec executes the process of creating a new node.
func (c *NodeCreation) Exec(ctx context.Context, request NodeCreationRequest) error {

	name := NewNodeName(request.Name)
	node := NewNode(*name)
	if err := c.nodeWriter.Save(nil, node); err != nil {
		return fmt.Errorf("failed to save new node: %w", err)
	}

	response := new(NodeCreationResponse)
	response.Created = node

	return c.presenter(ctx, response)
}
