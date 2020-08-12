package ns

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

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

type MockNodeWriter struct {
	mock.Mock
}

func (m *MockNodeWriter) Save(owner *User, node *Node) error {
	args := m.Called(owner, node)
	return args.Error(0)
}

type MockUserReader struct {
	mock.Mock
}

func (m *MockUserReader) GetByID(id UserID) (*User, error) {
	args := m.Called(id)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), nil
}

type TestNodeCreationContext struct {
	Request     NodeCreationRequest
	RespHandler NodeCreationResponseHandler
}

func SetupTestNodeCreationContext(t *testing.T) *TestNodeCreationContext {
	c := new(TestNodeCreationContext)
	c.Request = NodeCreationRequest{
		Name: "my namespace",
		UID:  uuid.New().String(),
	}
	c.RespHandler = func(ctx context.Context, resp *NodeCreationResponse) error {
		assert.NotNil(t, resp)
		assert.NotEmpty(t, resp, "resp must not be empty")
		return nil
	}
	return c
}

func TestNodeCreation_Exec(t *testing.T) {

	t.Run("Should launch a callback", func(t *testing.T) {
		// Setup
		assert := assert.New(t)
		c := SetupTestNodeCreationContext(t)

		// Setup Mocks
		nodeWriter := new(MockNodeWriter)
		nodeWriter.On("Save", mock.Anything, mock.Anything).Return(nil)

		userReader := new(MockUserReader)
		aUser := &User{ID: UserID(c.Request.UID)}
		userReader.On("GetByID", mock.Anything).Return(aUser, nil)

		// Exercise & Verify
		got := NewNodeCreation(nodeWriter, userReader, c.RespHandler).Exec(context.TODO(), c.Request)
		if assert.NoError(got) {
			nodeWriter.AssertNumberOfCalls(t, "Save", 1)
		}
	})
	t.Run("Must return an error on No-User-Found", func(t *testing.T) {
		// Setup
		assert := assert.New(t)
		c := SetupTestNodeCreationContext(t)
		// Setup Mocks
		var nodeWriter, userReader = new(MockNodeWriter), new(MockUserReader)
		{
			notFoundErr := fmt.Errorf("foo: %w", ErrNotFound)
			userReader.On("GetByID", mock.Anything).Return(nil, notFoundErr)
		}
		// Exercise & Verify
		got := NewNodeCreation(nodeWriter, userReader, c.RespHandler).Exec(context.TODO(), c.Request)
		if assert.Error(got) {
			assert.ErrorIs(got, ErrNotFound)
		}
	})
	t.Run("Must return an error on Failed-to-save", func(t *testing.T) {
		// Setup
		assert := assert.New(t)
		c := SetupTestNodeCreationContext(t)
		// Setup Mocks
		var nodeWriter, userReader = new(MockNodeWriter), new(MockUserReader)
		orgError := errors.New("this is a original error")
		{
			aUser := &User{ID: UserID(c.Request.UID)}
			userReader.On("GetByID", mock.Anything).Return(aUser, nil)
			nodeWriter.On("Save", aUser, mock.Anything).Return(orgError)
		}
		// Exercise & Verify
		got := NewNodeCreation(nodeWriter, userReader, c.RespHandler).Exec(context.TODO(), c.Request)
		if assert.Error(got) {
			assert.ErrorIs(got, orgError)
		}
	})
	t.Run("Must validate node-name", func(t *testing.T) {
		// Setup
		c := SetupTestNodeCreationContext(t)
		c.RespHandler = func(ctx context.Context, resp *NodeCreationResponse) error {
			t.Error("Must not be called XD")
			return nil
		}
		//  Setup Mocks
		var nodeWriter, userReader = new(MockNodeWriter), new(MockUserReader)

		aUser := &User{ID: UserID(c.Request.UID)}
		userReader.On("GetByID", mock.Anything).Return(aUser, nil)
		nodeWriter.On("Save", aUser, mock.Anything).Return(nil)

		sut := NewNodeCreation(nodeWriter, userReader, c.RespHandler)
		ctx := context.TODO()
		t.Run("node name with slash", func(t *testing.T) {
			got := sut.Exec(ctx, NodeCreationRequest{Name: "name/with/slash"})
			assert.Error(t, got)
			assert.ErrorIs(t, got, ErrIllegalArgument)
		})
		t.Run("empty node name", func(t *testing.T) {
			got := sut.Exec(ctx, NodeCreationRequest{Name: ""})
			assert.Error(t, got)
			assert.ErrorIs(t, got, ErrIllegalArgument)
		})
	})
}
