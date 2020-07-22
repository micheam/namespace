package postgres

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRowUser_AsNode(t *testing.T) {
	assert := assert.New(t)
	// setup
	var (
		id   = uuid.New().String()
		name = "taro"
	)
	sut := &RowUser{
		ID:   id,
		Name: name,
	}
	// exercise
	got, gotErr := sut.AsEntity()
	// verification
	assert.NoError(gotErr, "エラーとならないこと")
	assert.NotNil(got, "nilでないNodeをかえすこと")
	assert.EqualValues(id, got.ID, "Idが付与されていること")
	assert.EqualValues(name, got.Name, "Nameが付与されていること")
}

func TestUserRepository_GetByID(t *testing.T) {
}
