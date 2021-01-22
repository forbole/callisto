package database

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/simapp/params"

	"github.com/rs/zerolog/log"

	"github.com/desmos-labs/juno/config"
	"github.com/desmos-labs/juno/db"
	"github.com/desmos-labs/juno/db/postgresql"
	"github.com/jmoiron/sqlx"
)

// BigDipperDb represents a PostgreSQL database with expanded features.
// so that it can properly store custom BigDipper-related data.
type BigDipperDb struct {
	*postgresql.Database
	Sqlx *sqlx.DB
}

// Builder allows to create a new BigDipperDb instance implementing the database.Builder type
func Builder(cfg *config.Config, codec *params.EncodingConfig) (db.Database, error) {
	psqlConfig, ok := cfg.DatabaseConfig.Config.(*config.PostgreSQLConfig)
	if !ok {
		return nil, fmt.Errorf("MongoDB configuration is not supported on BigDipper")
	}

	database, err := postgresql.Builder(psqlConfig, codec)
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

// Cast allows to cast the given db to a BigDipperDb instance
func Cast(db db.Database) *BigDipperDb {
	bdDatabase, ok := db.(*BigDipperDb)
	if !ok {
		log.Fatal().Str("module", "supply").Msg("given database instance is not a BigDipperDb")
	}
	return bdDatabase
}
