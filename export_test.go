package ns

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
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

func PrepareTestUser(db *sqlx.DB) *User {
	// TODO: Add ns.NewUser then use it
	// TODO: Add NewUserRow then use it
	user := new(psqlUserRow)
	user.ID = uuid.New().String()
	user.Name = "test user"
	MustInsertUser(db, user)
	e, _ := user.toEntity()
	return e
}

func MustInsertUser(db *sqlx.DB, n *psqlUserRow) {
	_, err := db.NamedExec("INSERT INTO users (id, name) VALUES (:id, :name)", n)
	if err != nil {
		panic(err)
	}
}

func MustInsertNode(db *sqlx.DB, n *psqlNodeRow) {
	_, err := db.NamedExec("INSERT INTO nodes (id, name, user_id) VALUES (:id, :name, :user_id)", n)
	if err != nil {
		panic(err)
	}
}

func CleanupAll(db *sqlx.DB) {
	db.MustExec("TRUNCATE TABLE users CASCADE;")
	db.MustExec("TRUNCATE TABLE nodes;")
}
