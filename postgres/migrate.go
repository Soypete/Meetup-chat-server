package postgres

import (
	"github.com/pkg/errors"
	"github.com/pressly/goose"
)

/* TODO: This index is broken. need something to force uniqueness
CREATE INDEX IF NOT EXISTS index_users ON users (username, source); */

const (
	migrationTable = "goose_db_migrations"
	migrationPath  = "postgres/sql"
)

func (pg *PG) Migrate() error {
	goose.SetTableName(migrationTable)
	err := goose.Up(pg.Client.DB, migrationPath)
	if err != nil {
		return errors.Wrap(err, "failed migration")
	}

	return nil
}
