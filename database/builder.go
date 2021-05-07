package database

import (
	"fmt"

	"github.com/forbole/bdjuno/types/config"

	"github.com/cosmos/cosmos-sdk/simapp/params"
	"github.com/desmos-labs/juno/db"
	juno "github.com/desmos-labs/juno/types"

	bigdipperdb "github.com/forbole/bdjuno/database/bigdipper"
	forbolexdb "github.com/forbole/bdjuno/database/forbolex"
)

// Builder represents a db.Builder implementation that returns the proper database instance based on the type
// of application that this parser should be run for.
func Builder(junoCfg juno.Config, codec *params.EncodingConfig) (db.Database, error) {
	cfg, ok := junoCfg.(*config.Config)
	if !ok {
		panic(fmt.Errorf("invalid configuration type: %T", junoCfg))
	}

	switch cfg.GetDataType() {
	case config.DataTypeUpdated:
		return bigdipperdb.Builder(cfg, codec)

	case config.DataTypeHistoric:
		return forbolexdb.Builder(cfg, codec)

	default:
		panic(fmt.Errorf("invalid application type: %s", cfg.GetDataType()))
	}
}
