package postgres

import (
	"testing"

	"github.com/google/uuid"
	"micheam.com/ns"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

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
