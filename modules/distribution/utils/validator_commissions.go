package utils

import (
	"context"

	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/modules/utils"
	"github.com/forbole/bdjuno/types"

	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/rs/zerolog/log"
)

// UpdateValidatorsCommissionAmounts updates the validators commissions amounts
func UpdateValidatorsCommissionAmounts(height int64, client distrtypes.QueryClient, db *database.Db) {
	log.Debug().Str("module", "distribution").
		Int64("height", height).
		Msg("updating validators commissions")

	validators, err := db.GetValidators()
	if err != nil {
		log.Error().Str("module", "distribution").Err(err).
			Int64("height", height).
			Msg("error while getting validators")
		return
	}

	if len(validators) == 0 {
		// No validators, just skip
		return
	}

	// Get all the commissions
	for _, validator := range validators {
		go getValidatorCommission(height, client, validator, db)
	}
}

func getValidatorCommission(height int64, client distrtypes.QueryClient, validator types.Validator, db *database.Db) {
	res, err := client.ValidatorCommission(
		context.Background(),
		&distrtypes.QueryValidatorCommissionRequest{ValidatorAddress: validator.GetOperator()},
		utils.GetHeightRequestHeader(height),
	)
	if err != nil {
		log.Error().Str("module", "distribution").Err(err).
			Int64("height", height).
			Str("validator", validator.GetOperator()).
			Msg("error while getting validator commission")
		return
	}

	delegationAmount := types.NewValidatorCommissionAmount(
		validator.GetOperator(),
		validator.GetSelfDelegateAddress(),
		res.Commission.Commission,
		height,
	)

	err = db.SaveValidatorCommissionAmount(delegationAmount)
	if err != nil {
		log.Error().Str("module", "distribution").Err(err).
			Int64("height", height).
			Str("validator", validator.GetOperator()).
			Msg("error while saving validator commission amounts")
	}
}
