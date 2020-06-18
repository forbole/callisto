package staking

import (
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/desmos-labs/juno/parse/client"
	"github.com/forbole/bdjuno/database"
	"github.com/rs/zerolog/log"
)

// updateValidatorsUptime reads from the LCD the current staking pool and stores its value inside the database
func updateStakingPool(cp client.ClientProxy, db database.BigDipperDb) error {
	log.Debug().
		Str("module", "staking").
		Str("operation", "staking_pool").
		Msg("getting staking pool")

	var pool staking.Pool
	height, err := cp.QueryLCDWithHeight("/staking/pool", &pool)
	if err != nil {
		return err
	}

	log.Debug().
		Str("module", "staking").
		Str("operation", "staking_pool").
		Msg("saving staking pool")
	if err := db.SaveStakingPool(height, pool); err != nil {
		return err
	}

	return nil
}
