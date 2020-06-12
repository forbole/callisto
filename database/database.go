package database

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/desmos-labs/juno/config"
	"github.com/desmos-labs/juno/db"
	"github.com/desmos-labs/juno/db/postgresql"
	"github.com/jmoiron/sqlx"
)

// BigDipperDb represents a PostgreSQL database with expanded features.
// so that it can properly store custom BigDipper-related data.
type BigDipperDb struct {
	postgresql.Database
	sqlx *sqlx.DB
}

// Builder allows to create a new BigDipperDb instance implementing the database.Builder type
func Builder(cfg config.Config, codec *codec.Codec) (*db.Database, error) {
	psqlConfig, ok := cfg.DatabaseConfig.Config.(*config.PostgreSQLConfig)
	if !ok {
		// TODO: Support MongoDB
		return nil, fmt.Errorf("MongoDB configuration is not supported on BigDipper")
	}

	database, err := postgresql.Builder(*psqlConfig, codec)
	if err != nil {
		return nil, err
	}

	psqlDb, _ := (*database).(postgresql.Database)
	var bigDipperDb db.Database = BigDipperDb{
		Database: psqlDb,
		sqlx:     sqlx.NewDb(psqlDb.Sql, "postgresql"),
	}

	return &bigDipperDb, nil
}
