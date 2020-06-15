package staking

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/desmos-labs/juno/parse/client"
	"github.com/forbole/bdjuno/database"
	"github.com/rs/zerolog/log"
)

func updateValidatorsUptime(cp client.ClientProxy, db database.BigDipperDb) error {
	log.Debug().Msg("updating validators uptime")

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
	for _, info := range signingInfo {
		validatorUptime := database.ValidatorUptime{
			Height:              height,
			ValidatorAddress:    info.Address,
			SignedBlocksWindow:  params.SignedBlocksWindow,
			MissedBlocksCounter: info.MissedBlocksCounter,
		}

		if err := db.SaveValidatorUptime(validatorUptime); err != nil {
			log.Debug().Str("validator_address", info.Address.String()).Msg("saving validator update")
			return err
		}
	}

	return nil
}

func updateValidators(height int64, cp client.ClientProxy, db database.BigDipperDb) error {
	statuses := []string{"bonded", "unbonded", "unbonding"}

	// Get all the validators in any state
	var validators staking.Validators
	for _, status := range statuses {
		var validatorSet staking.Validators
		endpoint := fmt.Sprintf("/staking/validators?status=%s&height=%d", status, height)
		if _, err := cp.QueryLCDWithHeight(endpoint, &validatorSet); err != nil {
			return err
		}

		validators = append(validators, validatorSet...)
	}

	log.Debug().Int("validators", len(validators)).Msg("updating validators")
	return db.SaveValidators(validators)
}
