package ns

import (
	"database/sql"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRowNode_AsNode_returns_a_node(t *testing.T) {
	assert := assert.New(t)
	// setup
	var (
		id   = uuid.New()
		name = "this is a name"
		desc = "this is a desc"
	)
	sut := &psqlNodeRow{
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

func TestPSQLNodes_GetByID(t *testing.T) {
	var (
		assert = assert.New(t)
		db     = MustGetConn()
		owner  = new(User)
		id     = uuid.New().String()
	)
	user := &psqlUserRow{ID: uuid.New().String(), Name: "test user"}
	MustInsertUser(db, user)
	MustInsertNode(db, &psqlNodeRow{ID: id, Name: "test", UserID: user.ID})
	t.Cleanup(func() { CleanupAll(db) })
	sut := NewPSQLNodes(db)

	got, gotErr := sut.GetByID(owner, NodeID(id))
	assert.NoError(gotErr)
	assert.NotNil(got)
	assert.EqualValues(got.ID, id)
}

func TestPSQLNodes_GetByID_NoExist(t *testing.T) {
	// Setup
	assert := assert.New(t)
	db := MustGetConn()
	owner := new(User)
	id := uuid.New().String()
	t.Cleanup(func() { CleanupAll(db) })
	sut := NewPSQLNodes(db)
	// Exercise
	_, err := sut.GetByID(owner, NodeID(id))
	// Verification
	if assert.Error(err) {
		assert.ErrorIs(err, ErrNotFound)
	}
}

func TestPSQLNodes_Save(t *testing.T) {
	// Setup
	assert := assert.New(t)
	db := MustGetConn()
	owner := PrepareTestUser(db)

	sut := &PSQLNodes{db: db}
	node := &Node{
		ID:   *NewNodeID(),
		Name: *NewNodeName(uuid.New().String())}
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

func TestPSQLNodes_Save_Duplicated(t *testing.T) {
	// Setup
	assert := assert.New(t)
	db := MustGetConn()
	owner := PrepareTestUser(db)
	sut := &PSQLNodes{db: db}
	node := &Node{
		ID:   *NewNodeID(),
		Name: *NewNodeName(uuid.New().String())}
	err := sut.Save(owner, node)
	assert.NoError(err)
	// Exercise
	gotErr := sut.Save(owner, node)
	// Verify
	if assert.Error(gotErr) {
		assert.ErrorIs(gotErr, ErrDuplicatedEntity)
	}
}

func TestPSQLNodes_Save_IllegalNode(t *testing.T) {
	assert := assert.New(t)
	db := MustGetConn()
	sut := &PSQLNodes{db: db}
	owner := &User{}
	node := &Node{ /* ID の指定なし */ Name: *NewNodeName("aaa")}
	gotErr := sut.Save(owner, node)
	assert.Error(gotErr)
}

func TestPSQLNodes_Save_NilOwner(t *testing.T) {
	assert := assert.New(t)
	db := MustGetConn()
	sut := &PSQLNodes{db: db}
	owner := (*User)(nil)
	node := &Node{ /* ID の指定なし */ Name: *NewNodeName("aaa")}
	gotErr := sut.Save(owner, node)
	if assert.Error(gotErr) {
		assert.ErrorIs(gotErr, ErrIllegalArgument)
	}
}

func TestRowUser_AsNode(t *testing.T) {
	assert := assert.New(t)
	// setup
	var (
		id   = uuid.New().String()
		name = "taro"
	)
	sut := &psqlUserRow{
		ID:   id,
		Name: name,
	}
	// exercise
	got, gotErr := sut.toEntity()
	// verification
	assert.NoError(gotErr)
	assert.NotNil(got)
	assert.EqualValues(id, got.ID)
	assert.EqualValues(name, got.Name)
}

func TestPSQLUsers_GetByID(t *testing.T) {
	// FIXME
}
