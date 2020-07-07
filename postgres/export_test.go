package postgres

import "github.com/jmoiron/sqlx"

func MustInsertNode(db *sqlx.DB, n *RowNode) {
	_, err := db.NamedExec("INSERT INTO node (id, name) VALUES (:id, :name)", n)
	if err != nil {
		panic(err)
	}
}
func CleanupAll(db *sqlx.DB) {
	db.MustExec("TRUNCATE TABLE node;")
}
