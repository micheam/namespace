package postgres

import (
	"fmt"

	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"micheam.com/ns"
)

// RowNode is a row of node table
type RowNode struct {
	Id          string         `db:"id"`
	Name        string         `db:"name"`
	Description sql.NullString `db:"description"`
	CreatedAt   time.Time      `db:"created_at"`
	UpdatedAt   time.Time      `db:"updated_at"`
}

// AsEntity convert table model into entity
func (r *RowNode) AsEntity() (*ns.Node, error) {
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

type nodeRepository struct {
	db *sqlx.DB
}

// NewNodeRepository は、ns.NodeReader を初期化して返却する
func NewNodeRepository() (ns.NodeReader, error) {
	db, err := GetConn()
	if err != nil {
		return nil, fmt.Errorf("failed to init PostgresNodeReader: %w", err)
	}
	return &nodeRepository{db: db}, nil
}

// GetByID は、Postgresql から node を抽出して返却する
func (p *nodeRepository) GetByID(owner ns.User, id ns.NodeID) (*ns.Node, error) {
	var (
		row = new(RowNode)
		err = p.db.Get(row, "SELECT * FROM node WHERE id = $1", id)
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ns.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get node(id %s): %w", id, err)
	}
	return row.AsEntity()
}

// Save は、指定されたノードを保存する
//
// 指定された node が重複している場合は、 ns.ErrDuplicatedEntity を返却する。
func (n *nodeRepository) Save(owner *ns.User, node *ns.Node) error {
	var desc sql.NullString
	if node.Description != nil {
		desc = sql.NullString{
			String: node.Description.String(),
			Valid:  true,
		}
	}
	row := RowNode{
		Id:          node.ID.String(),
		Name:        node.Name.String(),
		Description: desc,
	}
	if _, err := n.db.NamedExec(
		"INSERT INTO node (id, name, description) VALUES (:id, :name, :description)", row); err != nil {
		return err
	}
	return nil
}
