package database

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/simapp/params"
	"github.com/desmos-labs/juno/db"
	juno "github.com/desmos-labs/juno/types"

	bigdipperdb "github.com/forbole/bdjuno/database/bigdipper"
	forbolexdb "github.com/forbole/bdjuno/database/forbolex"
	"github.com/forbole/bdjuno/types"
)

// Builder represents a db.Builder implementation that returns the proper database instance based on the type
// of application that this parser should be run for.
func Builder(cfg juno.Config, codec *params.EncodingConfig) (db.Database, error) {
	config, ok := cfg.(types.Config)
	if !ok {
		panic(fmt.Errorf("invalid configuration type: %T", cfg))
	}

	switch config.GetApplicationType() {
	case types.ApplicationTypeExplorer:
		return bigdipperdb.Builder(cfg, codec)

	case types.ApplicationTypeUtility:
		return forbolexdb.Builder(cfg, codec)

	default:
		panic(fmt.Errorf("invalid application type: %s", config.GetApplicationType()))
	}
}
