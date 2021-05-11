package staking

import (
	"context"
	"encoding/hex"

	"github.com/forbole/bdjuno/modules/common/staking"

	"github.com/cosmos/cosmos-sdk/codec"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	juno "github.com/desmos-labs/juno/types"

	bigdipperdb "github.com/forbole/bdjuno/database/bigdipper"
	bstakingtypes "github.com/forbole/bdjuno/modules/bigdipper/staking/types"
	utils2 "github.com/forbole/bdjuno/modules/common/utils"

	"github.com/rs/zerolog/log"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

// HandleBlock represents a method that is called each time a new block is created
func HandleBlock(
	block *tmctypes.ResultBlock, vals *tmctypes.ResultValidators,
	stakingClient stakingtypes.QueryClient,
	cdc codec.Marshaler, db *bigdipperdb.Db,
) error {
	// Update the validators
	validators, err := staking.UpdateValidators(block.Block.Height, stakingClient, cdc, db)
	if err != nil {
		return err
	}

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

// updateValidatorsStatus updates all validators' statuses
func updateValidatorsStatus(
	height int64, validators []stakingtypes.Validator, cdc codec.Marshaler, db *bigdipperdb.Db,
) {
	log.Debug().Str("module", "staking").Int64("height", height).
		Msg("updating validators statuses")

	statuses := make([]bstakingtypes.ValidatorStatus, len(validators))
	for index, validator := range validators {
		consAddr, err := staking.GetValidatorConsAddr(cdc, validator)
		if err != nil {
			log.Error().Str("module", "staking").Err(err).
				Int64("height", height).
				Str("address", validator.OperatorAddress).
				Msg("error while getting validator consensus address")
			return
		}

		consPubKey, err := staking.GetValidatorConsPubKey(cdc, validator)
		if err != nil {
			log.Error().Str("module", "staking").Err(err).
				Int64("height", height).
				Str("address", validator.OperatorAddress).
				Msg("error while getting validator consensus public key")
			return
		}

		statuses[index] = bstakingtypes.NewValidatorStatus(
			consAddr.String(),
			consPubKey.String(),
			int(validator.GetStatus()),
			validator.IsJailed(),
			height,
		)
	}

	err := db.SaveValidatorsStatuses(statuses)
	if err != nil {
		log.Error().Str("module", "staking").Err(err).Int64("height", height).
			Msg("error while saving validators statuses")
	}
}

// updateValidatorVotingPower fetches and stores into the database all the current validators' voting powers
func updateValidatorVotingPower(
	height int64, vals *tmctypes.ResultValidators, db *bigdipperdb.Db,
) {
	log.Debug().Str("module", "staking").Int64("height", height).
		Msg("updating validators voting powers")

	votingPowers := make([]bstakingtypes.ValidatorVotingPower, len(vals.Validators))
	for index, validator := range vals.Validators {
		consAddr := juno.ConvertValidatorAddressToBech32String(validator.Address)
		if found, _ := db.HasValidator(consAddr); !found {
			continue
		}

		votingPowers[index] = bstakingtypes.NewValidatorVotingPower(consAddr, validator.VotingPower, height)
	}

	err := db.SaveValidatorsVotingPowers(votingPowers)
	if err != nil {
		log.Error().Str("module", "staking").Err(err).Int64("height", height).
			Msg("error while saving validators voting powers")
	}
}

// updateDoubleSignEvidence updates the double sign evidence of all validators
func updateDoubleSignEvidence(height int64, evidenceList tmtypes.EvidenceList, db *bigdipperdb.Db) {
	log.Debug().Str("module", "staking").Int64("height", height).
		Msg("updating double sign evidence")

	for _, ev := range evidenceList {
		dve, ok := ev.(*tmtypes.DuplicateVoteEvidence)
		if !ok {
			continue
		}

		evidence := bstakingtypes.NewDoubleSignEvidence(
			height,
			bstakingtypes.NewDoubleSignVote(
				int(dve.VoteA.Type),
				dve.VoteA.Height,
				dve.VoteA.Round,
				dve.VoteA.BlockID.String(),
				juno.ConvertValidatorAddressToBech32String(dve.VoteA.ValidatorAddress),
				dve.VoteA.ValidatorIndex,
				hex.EncodeToString(dve.VoteA.Signature),
			),
			bstakingtypes.NewDoubleSignVote(
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
func updateStakingPool(height int64, stakingClient stakingtypes.QueryClient, db *bigdipperdb.Db) {
	log.Debug().Str("module", "staking").Int64("height", height).
		Msg("updating staking pool")

	res, err := stakingClient.Pool(
		context.Background(),
		&stakingtypes.QueryPoolRequest{},
		utils2.GetHeightRequestHeader(height),
	)
	if err != nil {
		log.Error().Str("module", "staking").Err(err).Int64("height", height).
			Msg("error while getting staking pool")
		return
	}

	err = db.SaveStakingPool(bstakingtypes.NewPool(res.Pool.BondedTokens, res.Pool.NotBondedTokens, height))
	if err != nil {
		log.Error().Str("module", "staking").Err(err).Int64("height", height).
			Msg("error while saving staking pool")
	}
}
