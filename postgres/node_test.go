package postgres

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
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
		Id:          id.String(),
		Name:        name,
		Description: sql.NullString{Valid: true, String: desc},
	}
	// exercise
	got, gotErr := sut.AsEntity()
	// verification
	assert.NoError(gotErr, "エラーとならないこと")
	assert.NotNil(got, "nilでないNodeをかえすこと")
	assert.EqualValues(id.String(), got.ID, "Idが付与されていること")
	assert.EqualValues(name, got.Name, "Nameが付与されていること")
	assert.NotNil(got.Description, "DescriptionがNilでないこと")
	assert.EqualValues(desc, *got.Description, "Descriptionがふよされていること")
}

func TestPostgresNodeReader_GetByID(t *testing.T) {
	var (
		assert = assert.New(t)
		db     = MustGetConn()
		owner  = ns.User{}
		id     = uuid.New().String()
	)
	MustInsertNode(db, &RowNode{Id: id, Name: "test"})
	t.Cleanup(func() { CleanupAll(db) })

	sut, err := NewNodeRepository()
	assert.NoError(err)

	got, gotErr := sut.GetByID(owner, ns.NodeID(id))
	assert.NoError(gotErr)
	assert.NotNil(got)
	assert.EqualValues(got.ID, id, "指定されたIDのNODEが抽出される")
}

func TestPostgresNodeReader_GetByID_NoExist(t *testing.T) {
	// Setup
	assert := assert.New(t)
	db := MustGetConn()
	owner := ns.User{}
	id := uuid.New().String()
	t.Cleanup(func() { CleanupAll(db) })
	sut, err := NewNodeRepository()
	assert.NoError(err)
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
	assert := assert.New(t)
	db := MustGetConn()
	sut := &nodeRepository{db: db}
	owner := &ns.User{}
	node := ns.NewNode("foo")
	gotErr := sut.Save(owner, node)
	if assert.NoError(gotErr) {
		got, gotErr := sut.GetByID(*owner, node.ID)
		assert.NoError(gotErr)
		assert.EqualValues(node.ID, got.ID, "Nodeが登録されていること")
	}
}

func TestNodeRepository_Save_IllegalNode(t *testing.T) {
	assert := assert.New(t)
	db := MustGetConn()
	sut := &nodeRepository{db: db}
	owner := &ns.User{}
	node := &ns.Node{ /* ID の指定なし */ Name: ns.NodeName("aaa")}
	gotErr := sut.Save(owner, node)
	assert.Error(gotErr)
}
