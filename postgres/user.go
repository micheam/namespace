package postgres

import (
	"database/sql"

	"micheam.com/ns"
)

// RowUser is a row of user table
type RowUser struct {
	ID        string       `db:"id"`
	Name      string       `db:"name"`
	CreatedAt sql.NullTime `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
}

// AsEntity convert table model into entity
func (r *RowUser) AsEntity() (*ns.User, error) {
	entity := &ns.User{
		ID:   ns.UserID(r.ID),
		Name: ns.UserName(r.Name),
	}
	return entity, nil
}

// UserRepository ...
type UserRepository struct{}

// GetByID ...
func (*UserRepository) GetByID(id ns.UserID) (*ns.User, error) {
	panic("Not Implemented Yet")
}
