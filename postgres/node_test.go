package postgres

import (
	"database/sql"
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
	got, gotErr := sut.AsNode()
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
		assert  = assert.New(t)
		db, err = GetConn()
		owner   = ns.User{}
		id      = uuid.New().String()
	)
	assert.NoError(err)
	MustInsertNode(db, &RowNode{Id: id, Name: "test"})
	t.Cleanup(func() { CleanupAll(db) })

	sut, err := NewPostgresNodeReader()
	assert.NoError(err)

	got, gotErr := sut.GetByID(owner, ns.NodeID(id))
	assert.NoError(gotErr)
	assert.NotNil(got)
	assert.EqualValues(got.ID, id, "指定されたIDのNODEが抽出される")
}

func TestPostgresNodeReader_GetByID_NoExist(t *testing.T) {
	var (
		assert  = assert.New(t)
		db, err = GetConn()
		owner   = ns.User{}
		id      = uuid.New().String()
	)
	assert.NoError(err)
	t.Cleanup(func() { CleanupAll(db) })

	sut, err := NewPostgresNodeReader()
	assert.NoError(err)

	_, gotErr := sut.GetByID(owner, ns.NodeID(id))
	t.Logf("%T", gotErr)
	assert.Error(gotErr)
}
