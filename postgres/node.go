package postgres

import (
	"fmt"
	"time"

	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"micheam.com/ns"
)

// RowNode is a row of node table
type RowNode struct {
	ID          string         `db:"id"`
	Name        string         `db:"name"`
	Description sql.NullString `db:"description"`
	CreatedAt   sql.NullTime   `db:"created_at"`
	UpdatedAt   sql.NullTime   `db:"updated_at"`
}

// AsEntity convert table model into entity
func (r *RowNode) AsEntity() (*ns.Node, error) {
	node := &ns.Node{
		ID:   ns.NodeID(r.ID),
		Name: *ns.NewNodeName(r.Name),
	}
	if r.Description.Valid {
		desc := ns.NodeDescription(r.Description.String)
		node.Description = &desc
	}
	return node, nil
}

// The NodeRepository provides access to node information.
type NodeRepository struct {
	db *sqlx.DB
}

// NewNodeRepository は、nodeReader を初期化して返却する
func NewNodeRepository(db *sqlx.DB) *NodeRepository {
	return &NodeRepository{db: db}
}

// GetByID は、Postgresql から node を抽出して返却する
func (n *NodeRepository) GetByID(owner *ns.User, id ns.NodeID) (*ns.Node, error) {
	var (
		row = new(RowNode)
		err = n.db.Get(row, "SELECT * FROM node WHERE id = $1", id)
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ns.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get node(id %s): %w", id, err)
	}
	return row.AsEntity()
}

// Save は、指定されたノードを Postgresql に登録する
//
// すでに登録されている（一制約違反する）場合は、 ns.ErrDuplicatedEntity を返却する。
func (n *NodeRepository) Save(owner *ns.User, node *ns.Node) error {

	var desc sql.NullString
	if node.Description != nil {
		desc = sql.NullString{
			String: node.Description.String(),
			Valid:  true,
		}
	}
	now := time.Now()
	row := RowNode{
		ID:          node.ID.String(),
		Name:        node.Name.String(),
		Description: desc,
		CreatedAt:   sql.NullTime{Valid: true, Time: now},
		UpdatedAt:   sql.NullTime{Valid: true, Time: now},
	}
	if _, err := n.db.NamedExec(
		"INSERT INTO node (id, name, description) VALUES (:id, :name, :description)", row); err != nil {

		if pqerr, ok := err.(*pq.Error); ok {
			if IsUniqueViolation(pqerr) {
				return ns.ErrDuplicatedEntity
			}
		}

		return err
	}

	// 結果をエンティティに反映
	node.CreatedAt = row.CreatedAt.Time
	node.UpdatedAt = row.UpdatedAt.Time

	return nil
}
