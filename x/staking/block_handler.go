package staking

import (
	"context"
	"encoding/hex"

	"github.com/forbole/bdjuno/x/utils"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	jutils "github.com/desmos-labs/juno/db/utils"

	"github.com/rs/zerolog/log"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/x/staking/types"
	stakingutils "github.com/forbole/bdjuno/x/staking/utils"
)

// HandleBlock represents a method that is called each time a new block is created
func HandleBlock(
	block *tmctypes.ResultBlock, vals *tmctypes.ResultValidators,
	stakingClient stakingtypes.QueryClient,
	db *database.BigDipperDb,
) error {
	// Update the staking pool
	err := updateStakingPool(block.Block.Height, stakingClient, db)
	if err != nil {
		log.Error().Str("module", "staking").Int64("height", block.Block.Height).
			Err(err).Msg("error while updating staking pool")
	}

	// Update the voting powers
	err = updateValidatorVotingPower(block.Block.Height, vals, db)
	if err != nil {
		return err
	}

	// Update the delegations
	err = updateValidatorsDelegations(block.Block.Height, stakingClient, db)
	if err != nil {
		log.Error().Str("module", "staking").Int64("height", block.Block.Height).
			Err(err).Msg("error while updating validators delegations")
	}

	// Update the validators statuses
	err = updateValidatorsStatus(block.Block.Height, stakingClient, db)
	if err != nil {
		log.Error().Str("module", "staking").Int64("height", block.Block.Height).
			Err(err).Msg("error while updating validators status")
	}

	// Updated the double sign evidences
	err = updateDoubleSignEvidence(block.Block.Evidence.Evidence, db)
	if err != nil {
		log.Error().Str("module", "staking").Int64("height", block.Block.Height).
			Err(err).Msg("error while updating double sign evidences")
	}

	return nil
}

// updateValidatorVotingPower fetches and stores into the database all the current validators' voting powers
func updateValidatorVotingPower(height int64, vals *tmctypes.ResultValidators, db *database.BigDipperDb) error {
	// Store the signing infos into the database
	log.Debug().Str("module", "staking").Int64("height", height).Msg("saving validators voting percentage")
	for _, validator := range vals.Validators {
		consAddr := jutils.ConvertValidatorAddressToBech32String(validator.Address)
		if found, _ := db.HasValidator(consAddr); !found {
			continue
		}

		err := db.SaveValidatorVotingPower(types.NewValidatorVotingPower(consAddr, validator.VotingPower, height))
		if err != nil {
			return err
		}
	}

	return nil
}

// updateStakingPool reads from the LCD the current staking pool and stores its value inside the database
func updateStakingPool(height int64, stakingClient stakingtypes.QueryClient, db *database.BigDipperDb) error {
	log.Debug().Str("module", "staking").Int64("height", height).
		Str("operation", "staking pool").Msg("getting staking pool")

	res, err := stakingClient.Pool(
		context.Background(),
		&stakingtypes.QueryPoolRequest{},
		utils.GetHeightRequestHeader(height),
	)
	if err != nil {
		return err
	}

	log.Debug().Str("module", "staking").Int64("height", height).
		Str("operation", "staking pool").Msg("saving staking pool")

	return db.SaveStakingPool(res.Pool, height)
}

// updateValidatorsDelegations updates the current validators delegations and stores them inside the database
func updateValidatorsDelegations(height int64, stakingClient stakingtypes.QueryClient, db *database.BigDipperDb) error {
	log.Debug().Str("module", "staking").Int64("height", height).
		Str("operation", "delegations").Msg("getting delegations")

	// Get the params
	params, err := db.GetStakingParams()
	if err != nil {
		return err
	}

	// Get the validators
	validators, err := db.GetValidators()
	if err != nil {
		return err
	}

	for _, validator := range validators {
		// Update the delegations
		delegations, err := stakingutils.GetDelegations(validator.GetOperator(), height, stakingClient)
		if err != nil {
			return err
		}

		err = db.SaveDelegations(delegations)
		if err != nil {
			return err
		}

		// Update the unbonding delegations
		unDels, err := stakingutils.GetUnbondingDelegations(validator.GetOperator(), params.BondName, height, stakingClient)
		if err != nil {
			return err
		}

		err = db.SaveUnbondingDelegations(unDels)
		if err != nil {
			return err
		}
	}

	return nil
}

// updateValidatorsStatus updates all validators' statuses
func updateValidatorsStatus(height int64, stakingClient stakingtypes.QueryClient, db *database.BigDipperDb) error {
	log.Debug().Str("module", "staking").Int64("height", height).
		Str("operation", "validators status").Msg("getting statuses")

	res, err := stakingClient.Validators(
		context.Background(),
		&stakingtypes.QueryValidatorsRequest{},
		utils.GetHeightRequestHeader(height),
	)
	if err != nil {
		return err
	}

	log.Debug().Str("module", "staking").Int64("height", height).
		Str("operation", "validators status").Msg("saving statuses")

	for _, validator := range res.Validators {
		consAddr, err := validator.GetConsAddr()
		if err != nil {
			return err
		}

		err = db.SaveValidatorStatus(types.NewValidatorStatus(
			consAddr.String(),
			int(validator.GetStatus()),
			validator.IsJailed(),
			height,
		))
		if err != nil {
			return err
		}
	}

	return nil
}

// updateDoubleSignEvidence updates the double sign evidence of all validators
func updateDoubleSignEvidence(evidenceList tmtypes.EvidenceList, db *database.BigDipperDb) error {
	for _, ev := range evidenceList {
		dve, ok := ev.(*tmtypes.DuplicateVoteEvidence)
		if !ok {
			continue
		}

		log.Debug().Str("module", "staking").
			Str("operation", "double sign evidence").Msg("saving evidence")

		evidence := types.NewDoubleSignEvidence(
			types.NewDoubleSignVote(
				int(dve.VoteA.Type),
				dve.VoteA.Height,
				dve.VoteA.Round,
				dve.VoteA.BlockID.String(),
				jutils.ConvertValidatorAddressToBech32String(dve.VoteA.ValidatorAddress),
				dve.VoteA.ValidatorIndex,
				hex.EncodeToString(dve.VoteA.Signature),
			),
			types.NewDoubleSignVote(
				int(dve.VoteB.Type),
				dve.VoteB.Height,
				dve.VoteB.Round,
				dve.VoteB.BlockID.String(),
				jutils.ConvertValidatorAddressToBech32String(dve.VoteB.ValidatorAddress),
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
