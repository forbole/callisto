package operations

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/desmos-labs/juno/parse/client"
	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/x/staking/types"
	"github.com/rs/zerolog/log"
)

// UpdateValidatorsUptime fetches from the REST APIs the uptimes for all the currently
// stored validators, later saving them into the database.
func UpdateValidatorsUptime(cp client.ClientProxy, db database.BigDipperDb) error {
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

func updateValidators(height int64, cp client.ClientProxy, db database.BigDipperDb) error {
	log.Debug().
		Str("module", "staking").
		Str("operation", "validators").
		Msg("getting validators data")

	statuses := []string{"bonded", "unbonded", "unbonding"}

	// Get all the validators in any state
	var validators []types.Validator
	for _, status := range statuses {
		var validatorSet []staking.Validator
		endpoint := fmt.Sprintf("/staking/validators?status=%s&height=%d", status, height)
		if _, err := cp.QueryLCDWithHeight(endpoint, &validatorSet); err != nil {
			return err
		}

		for _, validator := range validatorSet {
			validators = append(validators, validator)
		}
	}

	log.Debug().
		Str("module", "staking").
		Str("operation", "validators").
		Msg("saving validators data")
	return db.SaveValidatorsData(validators)
}
