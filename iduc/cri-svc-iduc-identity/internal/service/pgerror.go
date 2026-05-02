package service

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

// isUniqueViolation detecta el error 23505 de PostgreSQL.
func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}
