package staking

import (
	"encoding/hex"
	"fmt"
	"time"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/forbole/bdjuno/v2/types"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	juno "github.com/forbole/juno/v2/types"

	"github.com/rs/zerolog/log"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

// HandleBlock implements BlockModule
func (m *Module) HandleBlock(
	block *tmctypes.ResultBlock, res *tmctypes.ResultBlockResults, _ []*juno.Tx, vals *tmctypes.ResultValidators,
) error {
	// Update the validators
	validators, err := m.updateValidators(block.Block.Height)
	if err != nil {
		return fmt.Errorf("error while updating validators: %s", err)
	}

	// Update the voting powers
	go m.updateValidatorVotingPower(block.Block.Height, vals)

	// Update the validators statuses
	go m.updateValidatorsStatus(block.Block.Height, validators)

	// Updated the double sign evidences
	go m.updateDoubleSignEvidence(block.Block.Height, block.Block.Evidence.Evidence)

	// Update the staking pool
	go m.updateStakingPool(block.Block.Height)

	// Update redelegations and unbonding delegations
	go m.updateElapsedDelegations(block.Block.Height, block.Block.Time, res.EndBlockEvents)

	return nil
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
func (m *Module) updateElapsedDelegations(height int64, blockTime time.Time, events []abci.Event) {
	log.Debug().Str("module", "staking").Int64("height", height).
		Msg("updating elapsed redelegations and unbonding delegations")

	// Get past delegators to be refreshed now
	delegators, err := m.db.DeleteDelegatorsToRefresh(height)
	if err != nil {
		log.Error().Str("module", "staking").Err(err).Int64("height", height).
			Msg("error while getting delegators to refresh")
		return
	}

	// Delete the completed entries if there is someone to refresh it for
	if len(delegators) > 0 {
		err = m.db.DeleteCompletedRedelegations(blockTime)
		if err != nil {
			log.Error().Str("module", "staking").Err(err).Int64("height", height).
				Msg("error while deleting completed redelegations")
			return
		}

		err = m.db.DeleteCompletedUnbondingDelegations(blockTime)
		if err != nil {
			log.Error().Str("module", "staking").Err(err).Int64("height", height).
				Msg("error while deleting completed unbonding delegations")
			return
		}
	}

	// Update the delegations and balances of all the delegators
	for _, delegator := range delegators {
		err = m.refreshDelegatorDelegations(height, delegator)
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

	// Get all the events that identify a completed unbonding delegation or redelegations
	var completedEvents []abci.Event
	completedEvents = append(completedEvents, juno.FindEventsByType(events, stakingtypes.EventTypeCompleteUnbonding)...)
	completedEvents = append(completedEvents, juno.FindEventsByType(events, stakingtypes.EventTypeCompleteRedelegation)...)

	// Get the address of all the delegators to be refreshed
	var delegatorsToRefresh []string
	for _, event := range completedEvents {
		attr, err := juno.FindAttributeByKey(event, stakingtypes.AttributeKeyDelegator)
		if err != nil {
			log.Error().Str("module", "staking").Err(err).Int64("height", height).
				Msgf("error while getting %s attribute", stakingtypes.AttributeKeyDelegator)
			return
		}
		delegatorsToRefresh = append(delegatorsToRefresh, string(attr.Value))
	}

	// Store the delegators to refresh at the next height
	err = m.db.SaveDelegatorsToRefresh(height, delegatorsToRefresh)
	if err != nil {
		log.Error().Str("module", "staking").Err(err).Int64("height", height).
			Msg("error while saving delegators to refresh")
		return
	}
}
