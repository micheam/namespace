package ns

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// ============================================
// Node
func TestNode_WithDesc(t *testing.T) {
	sut := NewNode("foo")
	sut.WithDesc("bar")
	assert.EqualValues(t, "foo", sut.Name.String())
	assert.NotNil(t, sut.Description)
	assert.EqualValues(t, "bar", *sut.Description)
	assert.NotEqualValues(t, new(time.Time), sut.CreatedAt)
	assert.NotEqualValues(t, new(time.Time), sut.UpdatedAt)
}

func TestNode(t *testing.T) {
	id := NewNodeID()
	node := Node{
		ID:          *id,
		Name:        *NewNodeName("foo"),
		Description: nil,
	}
	assert := assert.New(t)
	assert.Equal(*id, node.ID)
}

// ============================================
// NodeID
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

// ============================================
// NodeName
func TestNodeName(t *testing.T) {
	t.Run("generate from valid name", func(t *testing.T) {
		assert := assert.New(t)
		n := NewNodeName("valid name")
		assert.NotEmpty(n.String())
		assert.True(n.Valid())
	})
	t.Run("name with slash", func(t *testing.T) {
		assert := assert.New(t)
		got := NewNodeName("this is / invalid name")
		assert.NotEmpty(got)
		assert.False(got.Valid())
	})
}
