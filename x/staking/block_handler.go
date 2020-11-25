package staking

import (
	"fmt"
	"github.com/forbole/bdjuno/x/staking/utils"
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
	err := updateStakingPool(block.Block.Height, w.ClientProxy, bigDipperDb)
	if err != nil {
		return err
	}

	// Update the delegations
	err = updateValidatorsDelegations(block.Block.Height, block.Block.Time, w.ClientProxy, bigDipperDb)
	if err != nil {
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

// updateDelegations reads from the LCD the current delegations and stores them inside the database
func updateValidatorsDelegations(height int64, timestamp time.Time, cp client.ClientProxy, db database.BigDipperDb) error {
	log.Debug().
		Str("module", "staking").
		Str("operation", "delegations").
		Msg("getting delegations")

	// Get the params
	params, err := db.GetStakingParams()
	if err != nil {
		return err
	}

	// Get the validators
	validators, err := db.GetValidators()
	if err != nil {
		return err
	}

	for _, validator := range validators {
		// Update the delegations
		delegations, err := utils.GetDelegations(validator.GetOperator(), height, timestamp, cp)
		if err != nil {
			return err
		}

		err = db.SaveCurrentDelegations(delegations)
		if err != nil {
			return err
		}

		// Update the unbonding delegations
		unDels, err := utils.GetUnbondingDelegations(validator.GetOperator(), params.BondName, height, timestamp, cp)
		if err != nil {
			return err
		}

		err = db.SaveCurrentUnbondingDelegations(unDels)
		if err != nil {
			return err
		}
	}

	return nil
}
