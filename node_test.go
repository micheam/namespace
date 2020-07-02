package ns

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNode(t *testing.T) {
	id := ID(uuid.New())
	node := Node{
		id:          id,
		name:        "foo",
		description: "my first node",
	}
	assert := assert.New(t)
	assert.Equal(id, node.id)
}
