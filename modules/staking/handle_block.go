package staking

import (
	"context"
	"encoding/hex"

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
	block *tmctypes.ResultBlock, vals *tmctypes.ResultValidators,
	stakingClient stakingtypes.QueryClient,
	cdc codec.Marshaler, db *database.Db,
) error {
	// Update the validators
	validators, err := stakingutils.UpdateValidators(block.Block.Height, stakingClient, cdc, db)
	if err != nil {
		return err
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
