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

	valid := true
	if matched, _ := regexp.Match(`/`, []byte(s)); matched {
		valid = false
	}
	if s == "" {
		valid = false
	}

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

// NodeReader ...
type NodeReader interface {
	GetByID(owner *User, id NodeID) (*Node, error)
}

// NodeWriter ...
type NodeWriter interface {
	Save(owner *User, node *Node) error
}

// NodeReadWriter ...
type NodeReadWriter interface {
	NodeReader
	NodeWriter
}

// NodeCreation is a UseCase.
type NodeCreation struct {
	nodeWriter NodeWriter
	userReader UserReader
	handleRes  NodeCreationResponseHandler
}

// NewNodeCreation return NodeCreation interactor.
func NewNodeCreation(
	w NodeWriter, u UserReader, p NodeCreationResponseHandler) *NodeCreation {
	return &NodeCreation{
		nodeWriter: w,
		userReader: u,
		handleRes:  p,
	}
}

type (
	// NodeCreationRequest is a request data of new node creation.
	NodeCreationRequest struct {
		Name string
		UID  string
	}
	// NodeCreationResponse is a response data of new node creation.
	NodeCreationResponse struct {
		Created *Node
	}
)

// NodeCreationResponseHandler defines how to output the result on new node creation.
type NodeCreationResponseHandler func(ctx context.Context, resp *NodeCreationResponse) error

// Exec executes the process of creating a new node.
func (c *NodeCreation) Exec(ctx context.Context, request NodeCreationRequest) error {

	var err error
	var user *User
	if user, err = c.userReader.GetByID(UserID(request.UID)); err != nil {
		return fmt.Errorf("user %s: %w", request.UID, err)
	}

	nodeName := NewNodeName(request.Name)
	if !nodeName.Valid() {
		return fmt.Errorf("node name: %w", ErrIllegalArgument)
	}

	node := NewNode(*nodeName)
	if err := c.nodeWriter.Save(user, node); err != nil {
		return fmt.Errorf("failed to save new node: %w", err)
	}

	response := new(NodeCreationResponse)
	response.Created = node

	return c.handleRes(ctx, response)
}
