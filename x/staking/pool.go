package staking

import (
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/desmos-labs/juno/parse/client"
	"github.com/forbole/bdjuno/database"
	"github.com/rs/zerolog/log"
)

// updateStakingPool reads from the LCD the current staking pool and stores its value inside the database
func updateStakingPool(cp client.ClientProxy, db database.BigDipperDb) {
	var pool staking.Pool
	height, err := cp.QueryLCDWithHeight("/staking/pool", &pool)
	if err != nil {
		log.Error().Err(err)
	}

	if err := db.SaveStakingPool(height, pool); err != nil {
		log.Error().Err(err)
	}
}
