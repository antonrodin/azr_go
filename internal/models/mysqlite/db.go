package mysqlite

import "database/sql"

type DbSqlite struct {
	DB *sql.DB
}
