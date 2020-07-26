package ns

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// GetPSQLConn は、postgres に接続する
func GetPSQLConn() (*sqlx.DB, error) {
	return sqlx.Connect(
		"postgres",
		"user=postgres dbname=ns password=passwd sslmode=disable",
	)
}

const (
	uniquenessViolation = pq.ErrorCode("23505")
)

// IsUniqueViolation は "unique_violation" エラー かどうかを判定する
func IsUniqueViolation(pqerr *pq.Error) bool {
	return pqerr.Code == uniquenessViolation
}

// psqlNodeRow is a row of node table
type psqlNodeRow struct {
	UserID      string         `db:"user_id"`
	ID          string         `db:"id"`
	Name        string         `db:"name"`
	Description sql.NullString `db:"description"`
	CreatedAt   sql.NullTime   `db:"created_at"`
	UpdatedAt   sql.NullTime   `db:"updated_at"`
}

// AsEntity convert table model into entity
func (r *psqlNodeRow) AsEntity() (*Node, error) {
	node := &Node{
		ID:   NodeID(r.ID),
		Name: *NewNodeName(r.Name),
	}
	if r.Description.Valid {
		desc := NodeDescription(r.Description.String)
		node.Description = &desc
	}
	return node, nil
}

// The PSQLNodes provides access to node information.
type PSQLNodes struct {
	db *sqlx.DB
}

// NewPSQLNodes ...
func NewPSQLNodes(db *sqlx.DB) NodeReadWriter {
	return &PSQLNodes{db: db}
}

// GetByID は、Postgresql から node を抽出して返却する
func (n *PSQLNodes) GetByID(owner *User, id NodeID) (*Node, error) {
	var (
		row = new(psqlNodeRow)
		err = n.db.Get(row, "SELECT * FROM nodes WHERE id = $1", id)
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get node(id %s): %w", id, err)
	}
	return row.AsEntity()
}

// Save は、指定されたノードを Postgresql に登録する
//
// すでに登録されている（一制約違反する）場合は、 ErrDuplicatedEntity を返却する。
func (n *PSQLNodes) Save(owner *User, node *Node) error {

	var desc sql.NullString
	if node.Description != nil {
		desc = sql.NullString{
			String: node.Description.String(),
			Valid:  true,
		}
	}
	now := time.Now()
	row := psqlNodeRow{
		UserID:      owner.ID.String(),
		ID:          node.ID.String(),
		Name:        node.Name.String(),
		Description: desc,
		CreatedAt:   sql.NullTime{Valid: true, Time: now},
		UpdatedAt:   sql.NullTime{Valid: true, Time: now},
	}
	if _, err := n.db.NamedExec(
		`INSERT INTO nodes
         (id, name, description, user_id) 
         VALUES (:id, :name, :description, :user_id)`, row); err != nil {

		if pqerr, ok := err.(*pq.Error); ok {
			if IsUniqueViolation(pqerr) {
				return ErrDuplicatedEntity
			}
		}

		return err
	}

	// 結果をエンティティに反映
	node.CreatedAt = row.CreatedAt.Time
	node.UpdatedAt = row.UpdatedAt.Time

	return nil
}

// psqlUserRow is a row of user table
type psqlUserRow struct {
	ID        string       `db:"id"`
	Name      string       `db:"name"`
	CreatedAt sql.NullTime `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
}

// toEntity convert table model into entity
func (r *psqlUserRow) toEntity() (*User, error) {
	entity := &User{
		ID:   UserID(r.ID),
		Name: UserName(r.Name),
	}
	return entity, nil
}

// PSQLUsers ...
type PSQLUsers struct {
	db *sqlx.DB
}

// NewPSQLUsers ...
func NewPSQLUsers(db *sqlx.DB) *PSQLUsers {
	return &PSQLUsers{db: db}
}

// GetByID ...
func (*PSQLUsers) GetByID(id UserID) (*User, error) {
	panic("Not Implemented Yet")
}
