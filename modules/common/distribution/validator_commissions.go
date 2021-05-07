package distribution

import (
	"context"

	"github.com/forbole/bdjuno/modules/common/utils"
	"github.com/forbole/bdjuno/types"

	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/rs/zerolog/log"
)

// UpdateValidatorsCommissionAmounts updates the validators commissions amounts
func UpdateValidatorsCommissionAmounts(height int64, client distrtypes.QueryClient, db DB) {
	log.Debug().Str("module", "distribution").Int64("height", height).
		Msg("updating validators commissions")

	validators, err := db.GetValidatorsInfo()
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

func getValidatorCommission(height int64, client distrtypes.QueryClient, validator types.ValidatorInfo, db DB) {
	res, err := client.ValidatorCommission(
		context.Background(),
		&distrtypes.QueryValidatorCommissionRequest{ValidatorAddress: validator.ValidatorOperAddr},
		utils.GetHeightRequestHeader(height),
	)
	if err != nil {
		log.Error().Str("module", "distribution").Err(err).Int64("height", height).
			Str("validator", validator.ValidatorOperAddr).Msg("error while getting validator commission")
		return
	}

	err = db.SaveValidatorCommissionAmount(types.NewValidatorCommissionAmount(
		validator.ValidatorOperAddr,
		validator.ValidatorSelfDelegateAddr,
		res.Commission.Commission,
		height,
	))
	if err != nil {
		log.Error().Str("module", "distribution").Err(err).Int64("height", height).
			Str("validator", validator.ValidatorOperAddr).Msg("error while saving validator commission amount")
	}
}
