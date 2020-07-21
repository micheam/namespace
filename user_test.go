package ns

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUID_String(t *testing.T) {
	uuid := uuid.New()
	sut := UserID(uuid.String())
	assert.EqualValues(t, uuid.String(), sut.String())
}
