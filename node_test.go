package ns

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNode(t *testing.T) {
	id := NewNodeID()
	node := Node{
		ID:          *id,
		Name:        "foo",
		Description: nil,
	}
	assert := assert.New(t)
	assert.Equal(*id, node.ID)
}

func TestNodeID_New(t *testing.T) {
	got := NewNodeID()
	assert.NotNil(t, got)
}

func TestNodeID_String(t *testing.T) {
	// Nill
	nilID := (*NodeID)(nil)
	assert.Equal(t, "", nilID.String())

	// Not Nil
	id := NodeID("foo")
	assert.Equal(t, "foo", id.String())
}

func TestNode_WithDesc(t *testing.T) {
	sut := NewNode("foo")
	sut.WithDesc("bar")
	assert.EqualValues(t, "foo", sut.Name)
	assert.NotNil(t, sut.Description)
	assert.EqualValues(t, "bar", *sut.Description)
	assert.NotEqualValues(t, new(time.Time), sut.CreatedAt)
	assert.NotEqualValues(t, new(time.Time), sut.UpdatedAt)
}
