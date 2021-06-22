package slashing

import (
	"context"

	"github.com/desmos-labs/juno/client"

	slashingutils "github.com/forbole/bdjuno/modules/slashing/utils"

	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/types"

	"github.com/rs/zerolog/log"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"

	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
)

// HandleBlock represents a method that is called each time a new block is created
func HandleBlock(block *tmctypes.ResultBlock, slashingClient slashingtypes.QueryClient, db *database.Db) error {
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
func updateSigningInfo(height int64, slashingClient slashingtypes.QueryClient, db *database.Db) error {
	log.Debug().Str("module", "slashing").Int64("height", height).
		Msg("updating signing info")

	signingInfos, err := slashingutils.GetSigningInfos(height, slashingClient)
	if err != nil {
		return err
	}

	return db.SaveValidatorsSigningInfos(signingInfos)
}

// updateSlashingParams gets the slashing params for the given height, and stores them inside the database
func updateSlashingParams(height int64, slashingClient slashingtypes.QueryClient, db *database.Db) error {
	log.Debug().Str("module", "slashing").Int64("height", height).
		Msg("updating slashing params")

	res, err := slashingClient.Params(
		context.Background(),
		&slashingtypes.QueryParamsRequest{},
		client.GetHeightRequestHeader(height),
	)
	if err != nil {
		return err
	}

	return db.SaveSlashingParams(types.NewSlashingParams(res.Params, height))
}
