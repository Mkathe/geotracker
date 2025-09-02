package postgres

import (
	"database/sql"
	_ "github.com/lib/pq"
)

func LoadDatabase(connString string) (*sql.DB, error) {
	return sql.Open("postgres", connString)
}
