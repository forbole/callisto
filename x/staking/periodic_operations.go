package staking

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/rpc"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/desmos-labs/juno/client"
	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/x/staking/types"
	"github.com/forbole/bdjuno/x/utils"
)

// RegisterPeriodicOps returns the AdditionalOperation that periodically runs fetches from
// the LCD to make sure that constantly changing data are synced properly.
func RegisterPeriodicOps(scheduler *gocron.Scheduler, cp *client.Proxy, db *database.BigDipperDb) error {
	log.Debug().Str("module", "staking").Msg("setting up periodic tasks")

	// Setup a cron job to run every 15 seconds
	if _, err := scheduler.Every(15).Second().StartImmediately().Do(func() {
		utils.WatchMethod(func() error { return updateValidatorsUptime(cp, db) })
	}); err != nil {
		return err
	}

	if _, err := scheduler.Every(1).Minute().StartImmediately().Do(func() {
		utils.WatchMethod(func() error { return UpdateValidatorVotingPower(cp, db) })
	}); err != nil {
		return err
	}

	return nil
}

// updateValidatorsUptime fetches from the REST APIs the uptimes for all the currently
// stored validators, later saving them into the database.
func updateValidatorsUptime(cp *client.Proxy, db *database.BigDipperDb) error {
	log.Debug().
		Str("module", "staking").
		Str("operation", "uptime").
		Msg("getting validators uptime")

	// Get the staking parameters
	var params slashing.Params
	height, err := cp.QueryLCDWithHeight("/slashing/parameters", &params)
	if err != nil {
		return err
	}

	// Update the validators
	if err := updateValidators(height, cp, db); err != nil {
		return err
	}

	// Get the validator signing info
	var signingInfo []slashing.ValidatorSigningInfo
	endpoint := fmt.Sprintf("/slashing/signing_infos?height=%d", height)
	if _, err := cp.QueryLCDWithHeight(endpoint, &signingInfo); err != nil {
		return err
	}

	// Store the signing infos into the database
	log.Debug().
		Str("module", "staking").
		Str("operation", "uptime").
		Msg("saving validators uptime")

	for _, info := range signingInfo {
		validatorUptime := types.ValidatorUptime{
			Height:              height,
			ValidatorAddress:    info.Address,
			SignedBlocksWindow:  params.SignedBlocksWindow,
			MissedBlocksCounter: info.MissedBlocksCounter,
		}

		// Skip non existing validators
		if found, _ := db.HasValidator(info.Address.String()); !found {
			continue
		}

		// Save the validator uptime information
		if err := db.SaveValidatorUptime(validatorUptime); err != nil {
			return err
		}
	}

	return nil
}

// Update all the validators
func updateValidators(height int64, cp *client.Proxy, db *database.BigDipperDb) error {
	log.Debug().
		Str("module", "staking").
		Str("operation", "validators").
		Msg("getting validators data")

	// Get all the validators in any state
	var validators []types.Validator
	statuses := []string{"bonded", "unbonded", "unbonding"}
	for _, status := range statuses {
		var validatorSet []stakingtypes.Validator
		endpoint := fmt.Sprintf("/staking/validators?status=%s&height=%d", status, height)
		if _, err := cp.QueryLCDWithHeight(endpoint, &validatorSet); err != nil {
			return err
		}

		for _, validator := range validatorSet {
			validators = append(validators, types.NewValidator(
				validator.GetConsAddr(),
				validator.GetOperator(),
				validator.GetConsPubKey(),
				sdk.AccAddress(validator.GetOperator()),
				&validator.Commission.MaxChangeRate,
				&validator.Commission.MaxRate,
			))
		}
	}

	return db.SaveValidators(validators)
}

// UpdateValidatorVotingPower fetches and stores into the database all the current validators' voting powers
func UpdateValidatorVotingPower(cp *client.Proxy, db *database.BigDipperDb) error {
	log.Debug().
		Str("module", "staking").
		Str("operation", " voting percentage").
		Msg("getting validators voting percentage")

	// First, get the latest block height
	var block tmctypes.ResultBlock
	if err := cp.QueryLCD("/blocks/latest", &block); err != nil {
		return err
	}

	// Second, get the validators
	var validators rpc.ResultValidatorsOutput
	endpoint := fmt.Sprintf("/validatorsets/%d", block.Block.Height)
	_, err := cp.QueryLCDWithHeight(endpoint, &validators)
	if err != nil {
		return err
	}

	// Store the signing infos into the database
	log.Debug().
		Str("module", "staking").
		Str("operation", "uptime").
		Msg("saving voting powers")

	for _, validator := range validators.Validators {
		if found, _ := db.HasValidator(validator.Address.String()); !found {
			continue
		}
		consAddress, err := sdk.ConsAddressFromBech32(validator.Address.String())
		if err != nil {
			return err
		}

		err = db.SaveValidatorVotingPower(types.NewValidatorVotingPower(
			consAddress,
			validator.VotingPower,
			block.Block.Height,
			block.Block.Time,
		))
		if err != nil {
			return err
		}
	}

	return nil
}
