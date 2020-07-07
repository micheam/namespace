package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"micheam.com/ns"
)

type nodeReader struct {
	db *sqlx.DB
}

func NewPostgresNodeReader() (ns.NodeReader, error) {
	db, err := GetConn()
	if err != nil {
		return nil, fmt.Errorf("failed to init PostgresNodeReader: %w", err)
	}
	return &nodeReader{db: db}, nil
}

// GetByID は、Postgresql から node を抽出して返却する。
//
// TODO(micheam): interface sqlx.Queryer を使って sqlx.DB と sqlx.Tx を透過的に扱う
func (p *nodeReader) GetByID(owner ns.User, id ns.NodeID) (*ns.Node, error) {
	row := &RowNode{}
	if err := p.db.Get(row, "SELECT * FROM node WHERE id = $1", id); err != nil {
		return nil, fmt.Errorf("failed to get node(id %s): %w", id, err)
	}
	return row.AsNode()
}
