package database

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/simapp/params"

	"github.com/desmos-labs/juno/config"
	"github.com/desmos-labs/juno/db"
	"github.com/desmos-labs/juno/db/postgresql"
	"github.com/jmoiron/sqlx"
)

var _ db.Database = &BigDipperDb{}

// Cast allows to cast the given db to a BigDipperDb instance
func Cast(db db.Database) *BigDipperDb {
	bdDatabase, ok := db.(*BigDipperDb)
	if !ok {
		panic(fmt.Errorf("given database instance is not a BigDipperDb"))
	}
	return bdDatabase
}

// Builder allows to create a new BigDipperDb instance implementing the database.Builder type
func Builder(cfg *config.Config, codec *params.EncodingConfig) (db.Database, error) {
	database, err := postgresql.Builder(cfg.Database, codec)
	if err != nil {
		return nil, err
	}

	psqlDb, ok := (database).(*postgresql.Database)
	if !ok {
		return nil, fmt.Errorf("invalid configuration database, must be PostgreSQL")
	}

	return &BigDipperDb{
		Database: psqlDb,
		Sqlx:     sqlx.NewDb(psqlDb.Sql, "postgresql"),
	}, nil
}

// BigDipperDb represents a PostgreSQL database with expanded features.
// so that it can properly store custom BigDipper-related data.
type BigDipperDb struct {
	*postgresql.Database
	Sqlx *sqlx.DB
}
