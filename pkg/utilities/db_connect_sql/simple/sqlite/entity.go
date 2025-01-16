package sqlite

import "database/sql"

type Service interface {
	Init() *sql.DB
}
