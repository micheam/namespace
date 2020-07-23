package postgres

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"micheam.com/ns"
)

func TestRowNode_AsNode_returns_a_node(t *testing.T) {
	assert := assert.New(t)
	// setup
	var (
		id   = uuid.New()
		name = "this is a name"
		desc = "this is a desc"
	)
	sut := &RowNode{
		ID:          id.String(),
		Name:        name,
		Description: sql.NullString{Valid: true, String: desc},
	}
	// exercise
	got, gotErr := sut.AsEntity()
	// verification
	assert.NoError(gotErr)
	assert.NotNil(got)
	assert.EqualValues(id.String(), got.ID)
	if assert.NotEmpty(got.Name) {
		assert.EqualValues(name, got.Name.String())
	}
	assert.NotNil(got.Description)
	assert.EqualValues(desc, *got.Description)
}

func TestPostgresNodeReader_GetByID(t *testing.T) {
	var (
		assert = assert.New(t)
		db     = MustGetConn()
		owner  = new(ns.User)
		id     = uuid.New().String()
	)
	user := &RowUser{ID: uuid.New().String(), Name: "test user"}
	MustInsertUser(db, user)
	MustInsertNode(db, &RowNode{ID: id, Name: "test", UserID: user.ID})
	t.Cleanup(func() { CleanupAll(db) })
	sut := NewNodeRepository(db)

	got, gotErr := sut.GetByID(owner, ns.NodeID(id))
	assert.NoError(gotErr)
	assert.NotNil(got)
	assert.EqualValues(got.ID, id)
}

func TestPostgresNodeReader_GetByID_NoExist(t *testing.T) {
	// Setup
	assert := assert.New(t)
	db := MustGetConn()
	owner := new(ns.User)
	id := uuid.New().String()
	t.Cleanup(func() { CleanupAll(db) })
	sut := NewNodeRepository(db)
	// Exercise
	_, err := sut.GetByID(owner, ns.NodeID(id))
	// Verification
	if assert.Error(err) {
		assert.True(errors.Is(err, ns.ErrNotFound))
	}
}

func TestNodeRepository_Save(t *testing.T) {
	// Setup
	assert := assert.New(t)
	db := MustGetConn()
	owner := PrepareTestUser(db)

	sut := &NodeRepository{db: db}
	node := &ns.Node{
		ID:   *ns.NewNodeID(),
		Name: *ns.NewNodeName(uuid.New().String())}
	// Exercise
	gotErr := sut.Save(owner, node)
	// Verify
	if assert.NoError(gotErr) {
		assert.NotEmpty(node.CreatedAt)
		assert.NotEmpty(node.UpdatedAt)
		assert.EqualValues(node.CreatedAt, node.UpdatedAt)

		foundNode, err := sut.GetByID(owner, node.ID)
		assert.NoError(err)
		assert.EqualValues(node.ID, foundNode.ID)
	}
}

func TestNodeRepository_Save_Duplicated(t *testing.T) {
	// Setup
	assert := assert.New(t)
	db := MustGetConn()
	owner := PrepareTestUser(db)
	sut := &NodeRepository{db: db}
	node := &ns.Node{
		ID:   *ns.NewNodeID(),
		Name: *ns.NewNodeName(uuid.New().String())}
	err := sut.Save(owner, node)
	assert.NoError(err)
	// Exercise
	gotErr := sut.Save(owner, node)
	// Verify
	if assert.Error(gotErr) {
		assert.True(errors.Is(gotErr, ns.ErrDuplicatedEntity))
	}
}

func TestNodeRepository_Save_IllegalNode(t *testing.T) {
	assert := assert.New(t)
	db := MustGetConn()
	sut := &NodeRepository{db: db}
	owner := &ns.User{}
	node := &ns.Node{ /* ID の指定なし */ Name: *ns.NewNodeName("aaa")}
	gotErr := sut.Save(owner, node)
	assert.Error(gotErr)
}
