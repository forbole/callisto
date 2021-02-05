package distribution

import (
	"context"

	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/database"
	bdistrtypes "github.com/forbole/bdjuno/x/distribution/types"
	"github.com/forbole/bdjuno/x/utils"
)

// updateCommunityPool fetch total amount of coins in the system from RPC and store it into database
func updateCommunityPool(height int64, distrClient distrtypes.QueryClient, db *database.BigDipperDb) error {
	log.Debug().Str("module", "distribution").Int64("height", height).Msg("getting community pool")

	res, err := distrClient.CommunityPool(context.Background(), &distrtypes.QueryCommunityPoolRequest{})
	if err != nil {
		return err
	}

	// Store the signing infos into the database
	return db.SaveCommunityPool(res.Pool, height)
}

// updateValidatorsCommissionAmounts updates the validators commissions amounts
func updateValidatorsCommissionAmounts(height int64, client distrtypes.QueryClient, db *database.BigDipperDb) error {
	log.Debug().Str("module", "distribution").Int64("height", height).Msg("updating validators commissions")

	validators, err := db.GetValidators()
	if err != nil {
		return err
	}

	if len(validators) == 0 {
		// No validators, just skip
		return nil
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
