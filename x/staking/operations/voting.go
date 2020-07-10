package operations

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/x/slashing"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	sdk "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/desmos-labs/juno/parse/client"
	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/x/staking/types"
	"github.com/rs/zerolog/log"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"

)

//this fetch all 
func UpdateValidatorVotingPower(cp client.ClientProxy, db database.BigDipperDb)error{
	log.Debug().
		Str("module", "staking").
		Str("operation", "uptime").
		Msg("getting validators uptime")

	// Get the staking parameters
	var params []tmctypes.ValidatorInfo
	height, err := cp.QueryLCD("/validators/result/validators",param)
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