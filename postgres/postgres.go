package postgres

import (
	"github.com/jmoiron/sqlx"
)

// GetConn は、postgres に接続する
func GetConn() (*sqlx.DB, error) {
	return sqlx.Connect(
		"postgres",
		"user=postgres dbname=ns password=passwd sslmode=disable",
	)
}
