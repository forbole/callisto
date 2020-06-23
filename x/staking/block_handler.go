package staking

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/desmos-labs/juno/parse/client"
	"github.com/desmos-labs/juno/parse/worker"
	juno "github.com/desmos-labs/juno/types"
	"github.com/forbole/bdjuno/database"
	"github.com/rs/zerolog/log"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// BlockHandler represents a method that is called each time a new block is created
func BlockHandler(block *tmctypes.ResultBlock, _ []juno.Tx, _ *tmctypes.ResultValidators, w worker.Worker) error {
	bigDipperDb, ok := w.Db.(database.BigDipperDb)
	if !ok {
		return fmt.Errorf("provided database is not a BigDipper database")
	}

	// Update the staking pool
	if err := updateStakingPool(block.Block.Height, w.ClientProxy, bigDipperDb); err != nil {
		return err
	}

	return nil
}

// updateStakingPool reads from the LCD the current staking pool and stores its value inside the database
func updateStakingPool(height int64, cp client.ClientProxy, db database.BigDipperDb) error {
	log.Debug().
		Str("module", "staking").
		Str("operation", "staking_pool").
		Msg("getting staking pool")

	var pool staking.Pool
	endpoint := fmt.Sprintf("/staking/pool?height=%d", height)
	height, err := cp.QueryLCDWithHeight(endpoint, &pool)
	if err != nil {
		return err
	}

	log.Debug().
		Str("module", "staking").
		Str("operation", "staking_pool").
		Msg("saving staking pool")
	if err := db.SaveStakingPool(pool, height, time.Now()); err != nil {
		return err
	}

	return nil
}
