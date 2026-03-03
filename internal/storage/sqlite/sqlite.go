package sqlite

import "database/sql"

type Storage struct {
	db *sql.DB
}

func New() {
	const op = "storage.sqlite.New"

}
