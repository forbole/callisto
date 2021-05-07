package forbolex

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/simapp/params"
	"github.com/desmos-labs/juno/db"
	"github.com/desmos-labs/juno/db/postgresql"
	juno "github.com/desmos-labs/juno/types"
	"github.com/jmoiron/sqlx"
)

var _ db.Database = &Db{}

// Cast allows to cast the given db to a Db instance
func Cast(db db.Database) *Db {
	bdDatabase, ok := db.(*Db)
	if !ok {
		panic(fmt.Errorf("given bigdipper instance is not a Db"))
	}
	return bdDatabase
}

// Builder allows to create a new Db instance implementing the db.Builder type
func Builder(cfg juno.Config, codec *params.EncodingConfig) (db.Database, error) {
	database, err := postgresql.Builder(cfg.GetDatabaseConfig(), codec)
	if err != nil {
		return nil, err
	}

	psqlDb, ok := (database).(*postgresql.Database)
	if !ok {
		return nil, fmt.Errorf("invalid configuration for ForboleX database, must be PostgreSQL")
	}

	return &Db{
		Database: psqlDb,
		Sqlx:     sqlx.NewDb(psqlDb.Sql, "postgresql"),
	}, nil
}

// Db represents a PostgreSQL database with expanded features.
// so that it can properly store custom ForboleX-related data.
type Db struct {
	*postgresql.Database
	Sqlx *sqlx.DB
}
