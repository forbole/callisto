package distribution

import (
	"context"

	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/rs/zerolog/log"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/forbole/bdjuno/database"
	bdistrtypes "github.com/forbole/bdjuno/x/distribution/types"
	"github.com/forbole/bdjuno/x/utils"
)

// HandleBlock represents a method that is called each time a new block is created
func HandleBlock(block *tmctypes.ResultBlock, client distrtypes.QueryClient, db *database.BigDipperDb) error {

	// Update the validator commissions
	err := updateValidatorsCommissionAmounts(block.Block.Height, client, db)
	if err != nil {
		log.Error().Str("module", "distribution").Int64("height", block.Block.Height).
			Err(err).Msg("error while updating validators commissions")
	}

	// Update the delegators commissions amounts
	err = updateDelegatorsCommissionsAmounts(block.Block.Height, client, db)
	if err != nil {
		log.Error().Str("module", "distribution").Int64("height", block.Block.Height).
			Err(err).Msg("error while updating delegators commissions amounts")
	}

	return nil
}

// updateValidatorsCommissionAmounts updates the validators commissions amounts
func updateValidatorsCommissionAmounts(height int64, client distrtypes.QueryClient, db *database.BigDipperDb) error {
	log.Debug().Str("module", "distribution").Int64("height", height).Msg("updating validators commissions")

	validators, err := db.GetValidators()
	if err != nil {
		return err
	}

	heightHeader := utils.GetHeightRequestHeader(height)

	// Get all the commissions
	var commissions = make([]bdistrtypes.ValidatorCommissionAmount, len(validators))
	for index, validator := range validators {
		res, err := client.ValidatorCommission(
			context.Background(),
			&distrtypes.QueryValidatorCommissionRequest{ValidatorAddress: validator.ValAddress},
			heightHeader,
		)
		if err != nil {
			return err
		}

		commissions[index] = bdistrtypes.NewValidatorCommissionAmount(
			validator.ConsAddress,
			res.Commission.Commission,
		)
	}

	// Store the commissions
	return db.SaveValidatorCommissionAmounts(commissions, height)
}

// updateDelegatorsCommissionsAmounts updates the delegators commission amounts
func updateDelegatorsCommissionsAmounts(height int64, client distrtypes.QueryClient, db *database.BigDipperDb) error {
	log.Debug().Str("module", "distribution").Int64("height", height).Msg("updating delegators commissions")

	// Get the delegators
	delegators, err := db.GetDelegatorsForHeight(height)
	if err != nil {
		return err
	}

	if len(delegators) == 0 {
		log.Debug().Str("module", "distribution").Int64("height", height).
			Msg("no delegations found, make sure you are calling this module after the staking module")
		return nil
	}

	header := utils.GetHeightRequestHeader(height)

	// Get the rewards
	var rewards []bdistrtypes.DelegatorCommissionAmount
	for _, delegator := range delegators {
		res, err := client.DelegationTotalRewards(
			context.Background(),
			&distrtypes.QueryDelegationTotalRewardsRequest{DelegatorAddress: delegator},
			header,
		)
		if err != nil {
			return err
		}

		for _, reward := range res.Rewards {
			rewards = append(rewards, bdistrtypes.NewDelegatorCommissionAmount(
				reward.ValidatorAddress,
				delegator,
				reward.Reward,
			))
		}
	}

	// Save the rewards
	return db.SaveDelegatorsCommissionAmounts(rewards, height)
}
