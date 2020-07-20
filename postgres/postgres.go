package postgres

import (
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// GetConn は、postgres に接続する
func GetConn() (*sqlx.DB, error) {
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
