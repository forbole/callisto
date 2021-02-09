package staking

import (
	"context"
	"encoding/hex"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	juno "github.com/desmos-labs/juno/types"

	"github.com/forbole/bdjuno/x/utils"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/rs/zerolog/log"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/x/staking/common"
	"github.com/forbole/bdjuno/x/staking/types"
)

// HandleBlock represents a method that is called each time a new block is created
func HandleBlock(
	block *tmctypes.ResultBlock, vals *tmctypes.ResultValidators,
	stakingClient stakingtypes.QueryClient,
	cdc codec.Marshaler, db *database.BigDipperDb,
) error {
	// Update the validators
	validators, err := updateValidators(block.Block.Height, stakingClient, cdc, db)
	if err != nil {
		return err
	}

	// Update the delegations
	err = common.UpdateValidatorsDelegations(block.Block.Height, validators, stakingClient, db)
	if err != nil {
		log.Error().Str("module", "staking").Int64("height", block.Block.Height).
			Err(err).Msg("error while updating validators delegations")
	}

	// Update the unbonding delegations
	err = common.UpdateValidatorsUnbondingDelegations(block.Block.Height, validators, stakingClient, db)
	if err != nil {
		log.Error().Str("module", "staking").Int64("height", block.Block.Height).
			Err(err).Msg("error while updating validators unbonding delegations")
	}

	// Update the redelegations
	err = common.UpdateValidatorsRedelegations(block.Block.Height, validators, stakingClient, db)
	if err != nil {
		log.Error().Str("module", "staking").Int64("height", block.Block.Height).
			Err(err).Msg("error while updating validators redelegations")
	}

	// Update the voting powers
	err = updateValidatorVotingPower(block.Block.Height, vals, db)
	if err != nil {
		log.Error().Str("module", "staking").Int64("height", block.Block.Height).
			Err(err).Msg("error while updating validators voting powers")
	}

	// Update the validators statuses
	err = updateValidatorsStatus(block.Block.Height, validators, cdc, db)
	if err != nil {
		log.Error().Str("module", "staking").Int64("height", block.Block.Height).
			Err(err).Msg("error while updating validators status")
	}

	// Updated the double sign evidences
	err = updateDoubleSignEvidence(block.Block.Height, block.Block.Evidence.Evidence, db)
	if err != nil {
		log.Error().Str("module", "staking").Int64("height", block.Block.Height).
			Err(err).Msg("error while updating double sign evidences")
	}

	// Update the staking pool
	err = updateStakingPool(block.Block.Height, stakingClient, db)
	if err != nil {
		log.Error().Str("module", "staking").Int64("height", block.Block.Height).
			Err(err).Msg("error while updating staking pool")
	}

	return nil
}

// updateValidators updates the list of validators that are present at the given height
func updateValidators(
	height int64, client stakingtypes.QueryClient, cdc codec.Marshaler, db *database.BigDipperDb,
) ([]stakingtypes.Validator, error) {
	log.Debug().Str("module", "staking").Int64("height", height).
		Msg("updating validators")

	validators, err := common.GetValidators(height, client)
	if err != nil {
		return nil, err
	}

	var vals = make([]types.Validator, len(validators))
	for index, val := range validators {
		consAddr, err := common.GetValidatorConsAddr(cdc, val)
		if err != nil {
			return nil, err
		}

		consPubKey, err := common.GetValidatorConsPubKey(cdc, val)
		if err != nil {
			return nil, err
		}

		vals[index] = types.NewValidator(
			consAddr.String(),
			val.OperatorAddress,
			consPubKey.String(),
			sdk.AccAddress(consAddr).String(),
			&val.Commission.MaxChangeRate,
			&val.Commission.MaxRate,
		)
	}

	err = db.SaveValidators(vals)
	if err != nil {
		return nil, err
	}

	return validators, err
}

// updateValidatorsStatus updates all validators' statuses
func updateValidatorsStatus(
	height int64, validators []stakingtypes.Validator, cdc codec.Marshaler, db *database.BigDipperDb,
) error {
	log.Debug().Str("module", "staking").Int64("height", height).
		Msg("updating validators statuses")

	statuses := make([]types.ValidatorStatus, len(validators))
	for index, validator := range validators {
		consAddr, err := common.GetValidatorConsAddr(cdc, validator)
		if err != nil {
			return err
		}

		consPubKey, err := common.GetValidatorConsPubKey(cdc, validator)
		if err != nil {
			return err
		}

		statuses[index] = types.NewValidatorStatus(
			consAddr.String(),
			consPubKey.String(),
			int(validator.GetStatus()),
			validator.IsJailed(),
			height,
		)
	}

	return db.SaveValidatorsStatuses(statuses)
}

// updateValidatorVotingPower fetches and stores into the database all the current validators' voting powers
func updateValidatorVotingPower(
	height int64, vals *tmctypes.ResultValidators, db *database.BigDipperDb,
) error {
	log.Debug().Str("module", "staking").Int64("height", height).
		Msg("updating validators voting powers")

	votingPowers := make([]types.ValidatorVotingPower, len(vals.Validators))
	for index, validator := range vals.Validators {
		consAddr := juno.ConvertValidatorAddressToBech32String(validator.Address)
		if found, _ := db.HasValidator(consAddr); !found {
			continue
		}

		votingPowers[index] = types.NewValidatorVotingPower(consAddr, validator.VotingPower, height)
	}

	return db.SaveValidatorsVotingPowers(votingPowers)
}

// updateDoubleSignEvidence updates the double sign evidence of all validators
func updateDoubleSignEvidence(height int64, evidenceList tmtypes.EvidenceList, db *database.BigDipperDb) error {
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
			return err
		}

	}
	return nil
}

// updateStakingPool reads from the LCD the current staking pool and stores its value inside the database
func updateStakingPool(height int64, stakingClient stakingtypes.QueryClient, db *database.BigDipperDb) error {
	log.Debug().Str("module", "staking").Int64("height", height).
		Msg("updating staking pool")

	res, err := stakingClient.Pool(
		context.Background(),
		&stakingtypes.QueryPoolRequest{},
		utils.GetHeightRequestHeader(height),
	)
	if err != nil {
		return err
	}

	return db.SaveStakingPool(res.Pool, height)
}
