package common

import (
	"context"

	bigdipperdb "github.com/forbole/bdjuno/database/bigdipper"
	utils2 "github.com/forbole/bdjuno/modules/common/utils"
	"github.com/forbole/bdjuno/types"

	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/rs/zerolog/log"

	dbtypes "github.com/forbole/bdjuno/database/bigdipper/types"
)

// UpdateValidatorsCommissionAmounts updates the validators commissions amounts
func UpdateValidatorsCommissionAmounts(height int64, client distrtypes.QueryClient, db *bigdipperdb.Db) {
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
	height int64, client distrtypes.QueryClient, validator dbtypes.ValidatorData, db *bigdipperdb.Db,
) {
	res, err := client.ValidatorCommission(
		context.Background(),
		&distrtypes.QueryValidatorCommissionRequest{ValidatorAddress: validator.ValAddress},
		utils2.GetHeightRequestHeader(height),
	)
	if err != nil {
		log.Error().Str("module", "distribution").Err(err).Int64("height", height).
			Str("validator", validator.ValAddress).Msg("error while getting validator commission")
		return
	}

	err = db.SaveValidatorCommissionAmount(types.NewValidatorCommissionAmount(
		validator.ConsAddress,
		res.Commission.Commission,
		height,
	))
	if err != nil {
		log.Error().Str("module", "distribution").Err(err).Int64("height", height).
			Str("validator", validator.ValAddress).Msg("error while saving validator commission amount")
	}
}
