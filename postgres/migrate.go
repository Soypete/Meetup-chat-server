package postgres

import (
	"github.com/pressly/goose"
)

const (
	migrationTable = "goose_db_migrations"
	migrationPath  = "postgres/sql"
)

func (pg *PG) Migrate() error {
	goose.SetTableName(migrationTable)
	goose.Up(pg.Client.DB, migrationPath)
	return nil
}
