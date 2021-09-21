package staking

import (
	"context"
	"encoding/hex"
	"fmt"
	"time"

	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	bankutils "github.com/forbole/bdjuno/modules/bank/utils"
	historyutils "github.com/forbole/bdjuno/modules/history/utils"
	"github.com/forbole/bdjuno/modules/utils"

	"github.com/desmos-labs/juno/client"

	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/types"

	"github.com/cosmos/cosmos-sdk/codec"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	juno "github.com/desmos-labs/juno/types"

	"github.com/rs/zerolog/log"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"

	stakingutils "github.com/forbole/bdjuno/modules/staking/utils"
)

// HandleBlock represents a method that is called each time a new block is created
func HandleBlock(
	cfg juno.Config, block *tmctypes.ResultBlock, vals *tmctypes.ResultValidators,
	stakingClient stakingtypes.QueryClient, bankClient banktypes.QueryClient, distrClient distrtypes.QueryClient,
	cdc codec.Marshaler, db *database.Db,
) error {
	// Update the validators
	validators, err := stakingutils.UpdateValidators(block.Block.Height, stakingClient, cdc, db)
	if err != nil {
		return fmt.Errorf("error while updating validators: %s", err)
	}

	// Get the params
	go updateParams(block.Block.Height, stakingClient, db)

	// Update the voting powers
	go updateValidatorVotingPower(block.Block.Height, vals, db)

	// Update the validators statuses
	go updateValidatorsStatus(block.Block.Height, validators, cdc, db)

	// Updated the double sign evidences
	go updateDoubleSignEvidence(block.Block.Height, block.Block.Evidence.Evidence, db)

	// Update the staking pool
	go updateStakingPool(block.Block.Height, stakingClient, db)

	// Update redelegations and unbonding delegations
	go updateElapsedDelegations(cfg, block.Block.Height, block.Block.Time, stakingClient, bankClient, distrClient, db)

	return nil
}

// updateParams gets the updated params and stores them inside the database
func updateParams(height int64, stakingClient stakingtypes.QueryClient, db *database.Db) {
	log.Debug().Str("module", "staking").Int64("height", height).
		Msg("updating params")

	res, err := stakingClient.Params(
		context.Background(),
		&stakingtypes.QueryParamsRequest{},
		client.GetHeightRequestHeader(height),
	)
	if err != nil {
		log.Error().Str("module", "staking").Err(err).
			Int64("height", height).
			Msg("error while getting params")
		return
	}

	err = db.SaveStakingParams(types.NewStakingParams(res.Params, height))
	if err != nil {
		log.Error().Str("module", "staking").Err(err).
			Int64("height", height).
			Msg("error while saving params")
		return
	}
}

// updateValidatorsStatus updates all validators' statuses
func updateValidatorsStatus(height int64, validators []stakingtypes.Validator, cdc codec.Marshaler, db *database.Db) {
	log.Debug().Str("module", "staking").Int64("height", height).
		Msg("updating validators statuses")

	statuses, err := stakingutils.GetValidatorsStatuses(height, validators, cdc)
	if err != nil {
		log.Error().Str("module", "staking").Err(err).
			Int64("height", height).
			Send()
		return
	}

	err = db.SaveValidatorsStatuses(statuses)
	if err != nil {
		log.Error().Str("module", "staking").Err(err).
			Int64("height", height).
			Msg("error while saving validators statuses")
	}
}

// updateValidatorVotingPower fetches and stores into the database all the current validators' voting powers
func updateValidatorVotingPower(height int64, vals *tmctypes.ResultValidators, db *database.Db) {
	log.Debug().Str("module", "staking").Int64("height", height).
		Msg("updating validators voting powers")

	votingPowers := stakingutils.GetValidatorsVotingPowers(height, vals, db)

	err := db.SaveValidatorsVotingPowers(votingPowers)
	if err != nil {
		log.Error().Str("module", "staking").Err(err).Int64("height", height).
			Msg("error while saving validators voting powers")
	}
}

// updateDoubleSignEvidence updates the double sign evidence of all validators
func updateDoubleSignEvidence(height int64, evidenceList tmtypes.EvidenceList, db *database.Db) {
	log.Debug().Str("module", "staking").Int64("height", height).
		Msg("updating double sign evidence")

	for _, ev := range evidenceList {
		dve, ok := ev.(*tmtypes.DuplicateVoteEvidence)
		if !ok {
			continue
		}

		evidence := types.NewDoubleSignEvidence(
			height,
			types.NewDoubleSignVote(
				int(dve.VoteA.Type),
				dve.VoteA.Height,
				dve.VoteA.Round,
				dve.VoteA.BlockID.String(),
				juno.ConvertValidatorAddressToBech32String(dve.VoteA.ValidatorAddress),
				dve.VoteA.ValidatorIndex,
				hex.EncodeToString(dve.VoteA.Signature),
			),
			types.NewDoubleSignVote(
				int(dve.VoteB.Type),
				dve.VoteB.Height,
				dve.VoteB.Round,
				dve.VoteB.BlockID.String(),
				juno.ConvertValidatorAddressToBech32String(dve.VoteB.ValidatorAddress),
				dve.VoteB.ValidatorIndex,
				hex.EncodeToString(dve.VoteB.Signature),
			),
		)

		err := db.SaveDoubleSignEvidence(evidence)
		if err != nil {
			log.Error().Str("module", "staking").Err(err).Int64("height", height).
				Msg("error while saving double sign evidence")
			return
		}

	}
}

// updateStakingPool reads from the LCD the current staking pool and stores its value inside the database
func updateStakingPool(height int64, stakingClient stakingtypes.QueryClient, db *database.Db) {
	log.Debug().Str("module", "staking").Int64("height", height).
		Msg("updating staking pool")

	pool, err := stakingutils.GetStakingPool(height, stakingClient)
	if err != nil {
		log.Error().Str("module", "staking").Err(err).Int64("height", height).
			Msg("error while getting staking pool")
		return
	}

	err = db.SaveStakingPool(pool)
	if err != nil {
		log.Error().Str("module", "staking").Err(err).Int64("height", height).
			Msg("error while saving staking pool")
		return
	}
}

// updateElapsedDelegations updates the redelegations and unbonding delegations that have elapsed
func updateElapsedDelegations(
	cfg juno.Config, height int64, blockTime time.Time,
	stakingClient stakingtypes.QueryClient, bankClient banktypes.QueryClient, distrClient distrtypes.QueryClient,
	db *database.Db,
) {
	log.Debug().Str("module", "staking").Int64("height", height).
		Msg("updating elapsed redelegations and unbonding delegations")

	deletedRedelegations, err := db.DeleteCompletedRedelegations(blockTime)
	if err != nil {
		log.Error().Str("module", "staking").Err(err).Int64("height", height).
			Msg("error while deleting completed redelegations")
		return
	}

	deletedUnbondingDelegations, err := db.DeleteCompletedUnbondingDelegations(blockTime)
	if err != nil {
		log.Error().Str("module", "staking").Err(err).Int64("height", height).
			Msg("error while deleting completed unbonding delegations")
		return
	}

	var delegators = map[string]bool{}

	// Add all the delegators from the redelegations
	for _, redelegation := range deletedRedelegations {
		if _, ok := delegators[redelegation.DelegatorAddress]; !ok {
			delegators[redelegation.DelegatorAddress] = true
		}
	}

	// Add all the delegators from unbonding delegations
	for _, delegation := range deletedUnbondingDelegations {
		if _, ok := delegators[delegation.DelegatorAddress]; !ok {
			delegators[delegation.DelegatorAddress] = true
		}
	}

	// Update the delegations and balances of all the delegators
	for delegator := range delegators {
		err = stakingutils.RefreshDelegations(height, delegator, stakingClient, distrClient, db)
		if err != nil {
			log.Error().Str("module", "staking").Err(err).Int64("height", height).
				Str("delegator", delegator).Msg("error while refreshing the delegations")
			return
		}

		err = bankutils.RefreshBalances(height, []string{delegator}, bankClient, db)
		if err != nil {
			log.Error().Str("module", "staking").Err(err).Int64("height", height).
				Str("delegator", delegator).Msg("error while refreshing the balance")
			return
		}

		if utils.IsModuleEnabled(cfg, types.HistoryModuleName) {
			err = historyutils.UpdateAccountBalanceHistoryWithTime(delegator, blockTime, db)
			if err != nil {
				log.Error().Str("module", "staking").Err(err).Int64("height", height).
					Str("delegator", delegator).Msg("error while updating account balance history")
				return
			}
		}
	}
}
