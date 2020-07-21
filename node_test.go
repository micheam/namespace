package ns

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Node {{{1

func TestNode_WithDesc(t *testing.T) {
	name := NewNodeName("foo")
	sut := NewNode(*name)
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

// NodeID {{{1

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

// NodeName {{{1

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

// UseCases {{{1

type MockNodeWriter struct {
	mock.Mock
}

func (m *MockNodeWriter) Save(owner *User, node *Node) error {
	args := m.Called(owner, node)
	return args.Error(0)
}

func TestCreateNewNode_Exec(t *testing.T) {

	t.Run("ordinal", func(t *testing.T) {
		assert := assert.New(t)

		presenter := func(ctx context.Context, resp *NodeCreationResponse) error {
			if assert.NotNil(resp,
				"presenter must be executed with non nil data") {
				assert.NotEmpty(resp, "resp must not be empty")
			}
			return nil
		}

		nodeWriter := new(MockNodeWriter)
		nodeWriter.On("Save", mock.Anything, mock.Anything).Return(nil)

		sut := NewNodeCreation(nodeWriter, presenter)
		got := sut.Exec(context.TODO(), NodeCreationRequest{})
		if assert.NoError(got) {
			nodeWriter.AssertNumberOfCalls(t, "Save", 1)
		}
	})

	t.Run("faile to save", func(t *testing.T) {
		assert := assert.New(t)
		presenter := func(ctx context.Context, resp *NodeCreationResponse) error {
			t.Error("expected to not be called, but was")
			return nil
		}

		orgError := errors.New("this is a original error")
		nodeWriter := new(MockNodeWriter)
		nodeWriter.On("Save", mock.Anything, mock.Anything).Return(orgError)

		sut := NewNodeCreation(nodeWriter, presenter)
		got := sut.Exec(context.TODO(), NodeCreationRequest{})
		if assert.Error(got) {
			assert.True(nodeWriter.AssertNumberOfCalls(t, "Save", 1))
			assert.True(errors.Is(got, orgError))
		}
	})
}
