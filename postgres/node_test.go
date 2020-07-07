package postgres

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
