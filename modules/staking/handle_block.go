package staking

import (
	"encoding/hex"
	"fmt"
	"time"

	"github.com/forbole/bdjuno/types"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	juno "github.com/desmos-labs/juno/types"

	"github.com/rs/zerolog/log"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

// HandleBlock implements BlockModule
func (m *Module) HandleBlock(block *tmctypes.ResultBlock, _ []*juno.Tx, vals *tmctypes.ResultValidators) error {
	// Update the validators
	validators, err := m.updateValidators(block.Block.Height)
	if err != nil {
		return fmt.Errorf("error while updating validators: %s", err)
	}

	// Get the params
	go m.updateParams(block.Block.Height)

	// Update the voting powers
	go m.updateValidatorVotingPower(block.Block.Height, vals)

	// Update the validators statuses
	go m.updateValidatorsStatus(block.Block.Height, validators)

	// Updated the double sign evidences
	go m.updateDoubleSignEvidence(block.Block.Height, block.Block.Evidence.Evidence)

	// Update the staking pool
	go m.updateStakingPool(block.Block.Height)

	// Update redelegations and unbonding delegations
	go m.updateElapsedDelegations(block.Block.Height, block.Block.Time)

	return nil
}

// updateParams gets the updated params and stores them inside the database
func (m *Module) updateParams(height int64) {
	log.Debug().Str("module", "staking").Int64("height", height).
		Msg("updating params")

	params, err := m.source.GetParams(height)
	if err != nil {
		log.Error().Str("module", "staking").Err(err).
			Int64("height", height).
			Msg("error while getting params")
		return
	}

	err = m.db.SaveStakingParams(types.NewStakingParams(params, height))
	if err != nil {
		log.Error().Str("module", "staking").Err(err).
			Int64("height", height).
			Msg("error while saving params")
		return
	}
}

// updateValidatorsStatus updates all validators' statuses
func (m *Module) updateValidatorsStatus(height int64, validators []stakingtypes.Validator) {
	log.Debug().Str("module", "staking").Int64("height", height).
		Msg("updating validators statuses")

	statuses, err := m.GetValidatorsStatuses(height, validators)
	if err != nil {
		log.Error().Str("module", "staking").Err(err).
			Int64("height", height).
			Send()
		return
	}

	err = m.db.SaveValidatorsStatuses(statuses)
	if err != nil {
		log.Error().Str("module", "staking").Err(err).
			Int64("height", height).
			Msg("error while saving validators statuses")
	}
}

// updateValidatorVotingPower fetches and stores into the database all the current validators' voting powers
func (m *Module) updateValidatorVotingPower(height int64, vals *tmctypes.ResultValidators) {
	log.Debug().Str("module", "staking").Int64("height", height).
		Msg("updating validators voting powers")

	votingPowers := m.GetValidatorsVotingPowers(height, vals)

	err := m.db.SaveValidatorsVotingPowers(votingPowers)
	if err != nil {
		log.Error().Str("module", "staking").Err(err).Int64("height", height).
			Msg("error while saving validators voting powers")
	}
}

// updateDoubleSignEvidence updates the double sign evidence of all validators
func (m *Module) updateDoubleSignEvidence(height int64, evidenceList tmtypes.EvidenceList) {
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

		err := m.db.SaveDoubleSignEvidence(evidence)
		if err != nil {
			log.Error().Str("module", "staking").Err(err).Int64("height", height).
				Msg("error while saving double sign evidence")
			return
		}

	}
}

// updateStakingPool reads from the LCD the current staking pool and stores its value inside the database
func (m *Module) updateStakingPool(height int64) {
	log.Debug().Str("module", "staking").Int64("height", height).
		Msg("updating staking pool")

	pool, err := m.GetStakingPool(height)
	if err != nil {
		log.Error().Str("module", "staking").Err(err).Int64("height", height).
			Msg("error while getting staking pool")
		return
	}

	err = m.db.SaveStakingPool(pool)
	if err != nil {
		log.Error().Str("module", "staking").Err(err).Int64("height", height).
			Msg("error while saving staking pool")
		return
	}
}

// updateElapsedDelegations updates the redelegations and unbonding delegations that have elapsed
func (m *Module) updateElapsedDelegations(height int64, blockTime time.Time) {
	log.Debug().Str("module", "staking").Int64("height", height).
		Msg("updating elapsed redelegations and unbonding delegations")

	deletedRedelegations, err := m.db.DeleteCompletedRedelegations(blockTime)
	if err != nil {
		log.Error().Str("module", "staking").Err(err).Int64("height", height).
			Msg("error while deleting completed redelegations")
		return
	}

	deletedUnbondingDelegations, err := m.db.DeleteCompletedUnbondingDelegations(blockTime)
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
		err = m.refreshDelegations(height, delegator)
		if err != nil {
			log.Error().Str("module", "staking").Err(err).Int64("height", height).
				Str("delegator", delegator).Msg("error while refreshing the delegations")
			return
		}

		err = m.bankModule.RefreshBalances(height, []string{delegator})
		if err != nil {
			log.Error().Str("module", "staking").Err(err).Int64("height", height).
				Str("delegator", delegator).Msg("error while refreshing the balance")
			return
		}

		err = m.historyModule.UpdateAccountBalanceHistoryWithTime(delegator, blockTime)
		if err != nil {
			log.Error().Str("module", "staking").Err(err).Int64("height", height).
				Str("delegator", delegator).Msg("error while updating account balance history")
			return
		}

	}
}
