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

func TestNodeCreation_Exec_Ordinal(t *testing.T) {
	assert := assert.New(t)
	request := NodeCreationRequest{
		Name: "my namespace",
		UID:  uuid.New().String(),
	}

	presenter := func(ctx context.Context, resp *NodeCreationResponse) error {
		assert.NotNil(resp)
		assert.NotEmpty(resp, "resp must not be empty")
		return nil
	}

	nodeWriter := new(MockNodeWriter)
	nodeWriter.On("Save", mock.Anything, mock.Anything).Return(nil)

	userReader := new(MockUserReader)
	aUser := &User{ID: UserID(request.UID)}
	userReader.On("GetByID", mock.Anything).Return(aUser, nil)

	sut := NewNodeCreation(nodeWriter, userReader, presenter)
	got := sut.Exec(context.TODO(), request)
	if assert.NoError(got) {
		nodeWriter.AssertNumberOfCalls(t, "Save", 1)
	}
}

func TestNodeCreation_Exec_NoUserFound(t *testing.T) {
	assert := assert.New(t)
	request := NodeCreationRequest{
		Name: "my namespace",
		UID:  uuid.New().String(),
	}

	presenter := func(ctx context.Context, resp *NodeCreationResponse) error {
		assert.NotNil(resp)
		assert.NotEmpty(resp, "resp must not be empty")
		return nil
	}

	nodeWriter := new(MockNodeWriter)
	userReader := new(MockUserReader)

	notFoundErr := fmt.Errorf("foo: %w", ErrNotFound)
	userReader.On("GetByID", mock.Anything).Return(nil, notFoundErr)

	sut := NewNodeCreation(nodeWriter, userReader, presenter)
	got := sut.Exec(context.TODO(), request)
	if assert.Error(got) {
		assert.ErrorIs(got, ErrNotFound)
	}
}

func TestNodeCreation_Exec_ValidationArg(t *testing.T) {
	request := NodeCreationRequest{
		Name: "my namespace",
		UID:  uuid.New().String(),
	}
	presenter := func(ctx context.Context, resp *NodeCreationResponse) error {
		t.Error("expected to not be called, but was")
		return nil
	}

	aUser := &User{ID: UserID(request.UID)}
	nodeWriter := new(MockNodeWriter)
	userReader := new(MockUserReader)

	nodeWriter.On("Save", aUser, mock.Anything).Return(nil)
	userReader.On("GetByID", mock.Anything).Return(aUser, nil)
	sut := NewNodeCreation(nodeWriter, userReader, presenter)

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
}

func TestNodeCreation_Exec_FaileToSave(t *testing.T) {
	assert := assert.New(t)
	request := NodeCreationRequest{
		Name: "my namespace",
		UID:  uuid.New().String(),
	}
	presenter := func(ctx context.Context, resp *NodeCreationResponse) error {
		t.Error("expected to not be called, but was")
		return nil
	}

	aUser := &User{ID: UserID(request.UID)}

	orgError := errors.New("this is a original error")
	nodeWriter := new(MockNodeWriter)
	nodeWriter.On("Save", aUser, mock.Anything).Return(orgError)

	userReader := new(MockUserReader)
	userReader.On("GetByID", mock.Anything).Return(aUser, nil)
	sut := NewNodeCreation(nodeWriter, userReader, presenter)
	got := sut.Exec(context.TODO(), request)

	if assert.Error(got) {
		assert.True(nodeWriter.AssertNumberOfCalls(t, "Save", 1))
		assert.ErrorIs(got, orgError)
	}
}
