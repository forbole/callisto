package database

import (
	"fmt"

	codec "github.com/cosmos/cosmos-sdk/codec"
	db "github.com/forbole/juno/v6/database"
	"github.com/forbole/juno/v6/database/postgresql"
	"github.com/jmoiron/sqlx"
)

var _ db.Database = &Db{}

// Db represents a PostgreSQL database with expanded features.
// so that it can properly store custom BigDipper-related data.
type Db struct {
	*postgresql.Database
	Sqlx *sqlx.DB
	cdc  codec.Codec
}

// Builder allows to create a new Db instance implementing the db.Builder type
func Builder(cdc codec.Codec) db.Builder {
	return func(ctx *db.Context) (db.Database, error) {
		database, err := postgresql.Builder(ctx)
		if err != nil {
			return nil, err
		}

		psqlDb, ok := (database).(*postgresql.Database)
		if !ok {
			return nil, fmt.Errorf("invalid configuration database, must be PostgreSQL")
		}

		return &Db{
			Database: psqlDb,
			Sqlx:     sqlx.NewDb(psqlDb.SQL.DB, "postgresql"),
			cdc:      cdc,
		}, nil
	}
}

// Cast allows to cast the given db to a Db instance
func Cast(db db.Database) *Db {
	bdDatabase, ok := db.(*Db)
	if !ok {
		panic(fmt.Errorf("given database instance is not a Db"))
	}
	return bdDatabase
}
