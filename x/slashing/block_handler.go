package slashing

import (
	"context"

	"github.com/forbole/bdjuno/x/utils"

	"github.com/rs/zerolog/log"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"

	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"

	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/x/slashing/types"
)

// HandleBlock represents a method that is called each time a new block is created
func HandleBlock(block *tmctypes.ResultBlock, slashingClient slashingtypes.QueryClient, db *database.BigDipperDb) error {
	// Update the signing infos
	err := updateSigningInfo(block.Block.Height, slashingClient, db)
	if err != nil {
		log.Error().Str("module", "slashing").Int64("height", block.Block.Height).
			Err(err).Msg("error while updating signing info")
	}

	err = updateSlashingParams(block.Block.Height, slashingClient, db)
	if err != nil {
		log.Error().Str("module", "slashing").Int64("height", block.Block.Height).
			Err(err).Msg("error while updating params")
	}

	return nil
}

// updateSigningInfo reads from the LCD the current staking pool and stores its value inside the database
func updateSigningInfo(height int64, slashingClient slashingtypes.QueryClient, db *database.BigDipperDb) error {
	log.Debug().Str("module", "slashing").Int64("height", height).
		Str("operation", "signing info").Msg("getting signing info")

	res, err := slashingClient.SigningInfos(
		context.Background(),
		&slashingtypes.QuerySigningInfosRequest{},
		utils.GetHeightRequestHeader(height),
	)
	if err != nil {
		return err
	}

	log.Debug().Str("module", "slashing").Int64("height", height).
		Str("operation", "signing info").Msg("saving signing info")

	infos := make([]types.ValidatorSigningInfo, len(res.Info))
	for index, info := range res.Info {
		infos[index] = types.NewValidatorSigningInfo(
			info.Address,
			info.StartHeight,
			info.IndexOffset,
			info.JailedUntil,
			info.Tombstoned,
			info.MissedBlocksCounter,
			height,
		)
	}

	return db.SaveValidatorsSigningInfos(infos)
}

// updateSlashingParams gets the slashing params for the given height, and stores them inside the database
func updateSlashingParams(height int64, slashingClient slashingtypes.QueryClient, db *database.BigDipperDb) error {
	log.Debug().Str("module", "slashing").Int64("height", height).
		Str("operation", "params").Msg("getting slashing params")

	res, err := slashingClient.Params(
		context.Background(),
		&slashingtypes.QueryParamsRequest{},
		utils.GetHeightRequestHeader(height),
	)
	if err != nil {
		return err
	}

	return db.SaveSlashingParams(res.Params, height)
}
