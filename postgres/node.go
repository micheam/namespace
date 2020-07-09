package postgres

import (
	"fmt"

	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"micheam.com/ns"
)

// RowNode は、TODO
type RowNode struct {
	Id          string         `db:"id"`
	Name        string         `db:"name"`
	Description sql.NullString `db:"description"`
	CreatedAt   time.Time      `db:"created_at"`
	UpdatedAt   time.Time      `db:"updated_at"`
}

// AsNode は、 TODO
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

// nodeRepository は、 TODO
type nodeRepository struct {
	db *sqlx.DB
}

func NewNodeRepository() (ns.NodeReader, error) {
	db, err := GetConn()
	if err != nil {
		return nil, fmt.Errorf("failed to init PostgresNodeReader: %w", err)
	}
	return &nodeRepository{db: db}, nil
}

// GetByID は、Postgresql から node を抽出して返却する。
//
// TODO(micheam): interface sqlx.Queryer を使って sqlx.DB と sqlx.Tx を透過的に扱う
func (p *nodeRepository) GetByID(owner ns.User, id ns.NodeID) (*ns.Node, error) {

	row := new(RowNode)
	err := p.db.Get(row, "SELECT * FROM node WHERE id = $1", id)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ns.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get node(id %s): %w", id, err)
	}
	return row.AsNode()
}

func (n *nodeRepository) Save(owner ns.User, node ns.Node) error {
	panic("NOT IMPLEMENTED YET")
}
