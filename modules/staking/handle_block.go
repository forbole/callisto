package staking

import (
	"encoding/hex"
	"fmt"

	"github.com/forbole/bdjuno/v3/types"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	juno "github.com/forbole/juno/v3/types"

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

	// Get the voting powers
	votingPowers, err := m.GetValidatorsVotingPowers(height, vals)
	if err != nil {
		log.Error().Str("module", "staking").Err(err).Int64("height", height).
			Msg("error while getting validators voting powers")
		return
	}

	// Save all the voting powers
	err = m.db.SaveValidatorsVotingPowers(votingPowers)
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
