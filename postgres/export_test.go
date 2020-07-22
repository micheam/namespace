package postgres

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"micheam.com/ns"
)

func MustGetConn() *sqlx.DB {
	db, err := sqlx.Connect(
		"postgres", "user=postgres dbname=ns password=passwd sslmode=disable",
	)
	if err != nil {
		panic(err)
	}
	return db
}

func PrepareTestUser(db *sqlx.DB) *ns.User {
	// TODO: Add ns.NewUser then use it
	// TODO: Add NewUserRow then use it
	user := new(RowUser)
	user.ID = uuid.New().String()
	user.Name = "test user"
	MustInsertUser(db, user)
	e, _ := user.AsEntity()
	return e
}

func MustInsertUser(db *sqlx.DB, n *RowUser) {
	_, err := db.NamedExec("INSERT INTO users (id, name) VALUES (:id, :name)", n)
	if err != nil {
		panic(err)
	}
}

func MustInsertNode(db *sqlx.DB, n *RowNode) {
	_, err := db.NamedExec("INSERT INTO node (id, name, user_id) VALUES (:id, :name, :user_id)", n)
	if err != nil {
		panic(err)
	}
}

func CleanupAll(db *sqlx.DB) {
	db.MustExec("TRUNCATE TABLE users CASCADE;")
	db.MustExec("TRUNCATE TABLE node;")
}
