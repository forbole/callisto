package common

import (
	"context"

	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/database"
	dbtypes "github.com/forbole/bdjuno/database/types"
	bdistrtypes "github.com/forbole/bdjuno/x/distribution/types"
	"github.com/forbole/bdjuno/x/utils"
)

// UpdateValidatorsCommissionAmounts updates the validators commissions amounts
func UpdateValidatorsCommissionAmounts(height int64, client distrtypes.QueryClient, db *database.BigDipperDb) {
	log.Debug().Str("module", "distribution").Int64("height", height).
		Msg("updating validators commissions")

	validators, err := db.GetValidators()
	if err != nil {
		log.Error().Str("module", "distribution").Err(err).Int64("height", height).
			Msg("error while getting validators")
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

func getValidatorCommission(
	height int64, client distrtypes.QueryClient, validator dbtypes.ValidatorData, db *database.BigDipperDb,
) {
	res, err := client.ValidatorCommission(
		context.Background(),
		&distrtypes.QueryValidatorCommissionRequest{ValidatorAddress: validator.ValAddress},
		utils.GetHeightRequestHeader(height),
	)
	if err != nil {
		log.Error().Str("module", "distribution").Err(err).Int64("height", height).
			Str("validator", validator.ValAddress).Msg("error while getting validator commission")
		return
	}

	err = db.SaveValidatorCommissionAmount(bdistrtypes.NewValidatorCommissionAmount(
		validator.ConsAddress,
		res.Commission.Commission,
		height,
	))
	if err != nil {
		log.Error().Str("module", "distribution").Err(err).Int64("height", height).
			Str("validator", validator.ValAddress).Msg("error while saving validator commission amount")
	}
}
