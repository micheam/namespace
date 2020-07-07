package postgres

import (
	"database/sql"
	"time"

	"micheam.com/ns"
)

type RowNode struct {
	Id          string         `db:"id"`
	Name        string         `db:"name"`
	Description sql.NullString `db:"description"`
	CreatedAt   time.Time      `db:"created_at"`
	UpdatedAt   time.Time      `db:"updated_at"`
}

func (r *RowNode) AsNode() (*ns.Node, error) {
	node := &ns.Node{
		ID:   ns.NodeID(r.Id),
		Name: ns.NodeName(r.Name),
	}
	if r.Description.Valid {
		desc := ns.NodeDescription(r.Description.String)
		node.Description = &desc
	}
	return node, nil
}
