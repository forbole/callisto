package common

import (
	"context"
	"sync"

	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/database"
	dbtypes "github.com/forbole/bdjuno/database/types"
	bdistrtypes "github.com/forbole/bdjuno/x/distribution/types"
	"github.com/forbole/bdjuno/x/utils"
)

// UpdateValidatorsCommissionAmounts updates the validators commissions amounts
func UpdateValidatorsCommissionAmounts(height int64, client distrtypes.QueryClient, db *database.BigDipperDb) error {
	log.Debug().Str("module", "distribution").Int64("height", height).Msg("updating validators commissions")

	validators, err := db.GetValidators()
	if err != nil {
		return err
	}

	if len(validators) == 0 {
		// No validators, just skip
		return nil
	}

	// Get all the commissions
	var wg sync.WaitGroup
	var out = make(chan bdistrtypes.ValidatorCommissionAmount)
	for _, validator := range validators {
		wg.Add(1)
		go getValidatorCommission(height, client, validator, out, &wg)
	}

	// We need to call wg.Wait inside another goroutine in order to solve the hanging bug that's described here:
	// https://dev.to/sophiedebenedetto/synchronizing-go-routines-with-channels-and-waitgroups-3ke2
	go func() {
		wg.Wait()
		close(out)
	}()

	var commissions []bdistrtypes.ValidatorCommissionAmount
	for com := range out {
		commissions = append(commissions, com)
	}

	// Store the commissions
	return db.SaveValidatorCommissionAmounts(commissions, height)
}

func getValidatorCommission(
	height int64,
	client distrtypes.QueryClient,
	validator dbtypes.ValidatorData,
	out chan<- bdistrtypes.ValidatorCommissionAmount,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

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

	out <- bdistrtypes.NewValidatorCommissionAmount(validator.ConsAddress, res.Commission.Commission)
}
