package v5

import (
	"github.com/jmoiron/sqlx"

	"github.com/forbole/juno/v5/database"
	"github.com/forbole/juno/v5/database/postgresql"
)

var _ database.Migrator = &Migrator{}

// Migrator represents the database migrator that should be used to migrate from v4 of the database to v5
type Migrator struct {
	SQL *sqlx.DB
}

func NewMigrator(db *postgresql.Database) *Migrator {
	return &Migrator{
		SQL: db.SQL,
	}
}
