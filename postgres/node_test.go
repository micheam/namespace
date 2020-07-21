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
	assert.NoError(gotErr, "エラーとならないこと")
	assert.NotNil(got, "nilでないNodeをかえすこと")
	assert.EqualValues(id.String(), got.ID, "Idが付与されていること")
	if assert.True(got.Name.Valid()) {
		assert.EqualValues(name, got.Name.String(), "Nameが付与されていること")
	}
	assert.NotNil(got.Description, "DescriptionがNilでないこと")
	assert.EqualValues(desc, *got.Description, "Descriptionがふよされていること")
}

func TestPostgresNodeReader_GetByID(t *testing.T) {
	var (
		assert = assert.New(t)
		db     = MustGetConn()
		owner  = new(ns.User)
		id     = uuid.New().String()
	)
	MustInsertNode(db, &RowNode{ID: id, Name: "test"})
	t.Cleanup(func() { CleanupAll(db) })
	sut := NewNodeRepository(db)

	got, gotErr := sut.GetByID(owner, ns.NodeID(id))
	assert.NoError(gotErr)
	assert.NotNil(got)
	assert.EqualValues(got.ID, id, "指定されたIDのNODEが抽出される")
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
	_, gotErr := sut.GetByID(owner, ns.NodeID(id))
	// Verification
	t.Logf("%T", gotErr)
	if assert.Error(gotErr) {
		assert.True(errors.Is(gotErr, ns.ErrNotFound),
			"got error must be ns.ErrNotFound")
	}
}

func TestNodeRepository_Save(t *testing.T) {
	// Setup
	assert := assert.New(t)
	db := MustGetConn()
	sut := &NodeRepository{db: db}
	owner := &ns.User{}
	node := &ns.Node{
		ID:   *ns.NewNodeID(),
		Name: *ns.NewNodeName(uuid.New().String())}
	// Exercise
	gotErr := sut.Save(owner, node)
	// Verify
	if assert.NoError(gotErr) {
		assert.False(node.CreatedAt.IsZero(), "CreatedAt が登録時にセットされるべき")
		assert.False(node.UpdatedAt.IsZero(), "UpdatedAt が登録時にセットされるべき")
		assert.EqualValues(node.CreatedAt, node.UpdatedAt, "初期登録時は、CreatedAt と UpdatedAt が一致する")

		got, gotErr := sut.GetByID(owner, node.ID)
		assert.NoError(gotErr, "登録に")
		assert.EqualValues(node.ID, got.ID, "エンティティに指定された ID で登録されるべき")
	}
}

func TestNodeRepository_Save_Duplicated(t *testing.T) {
	// Setup
	assert := assert.New(t)
	db := MustGetConn()
	sut := &NodeRepository{db: db}
	owner := &ns.User{}
	node := &ns.Node{
		ID:   *ns.NewNodeID(),
		Name: *ns.NewNodeName(uuid.New().String())}
	err := sut.Save(owner, node)
	assert.NoError(err)
	// Exercise
	gotErr := sut.Save(owner, node)
	// Verify
	if assert.Error(gotErr) {
		assert.True(errors.Is(gotErr, ns.ErrDuplicatedEntity),
			"ns.ErrDuplicatedEntity が返却される")
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
