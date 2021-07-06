package mint

import (
	"context"

	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/desmos-labs/juno/client"
	"github.com/rs/zerolog/log"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/types"
)

// HandleBlock represents a method that is called each time a new block is created
func HandleBlock(block *tmctypes.ResultBlock, mintClient minttypes.QueryClient, db *database.Db) error {
	// Update the params
	go updateParams(block.Block.Height, mintClient, db)

	return nil
}

// updateParams gets the updated params and stores them inside the database
func updateParams(height int64, mintClient minttypes.QueryClient, db *database.Db) {
	log.Debug().Str("module", "mint").Int64("height", height).
		Msg("updating params")

	res, err := mintClient.Params(
		context.Background(),
		&minttypes.QueryParamsRequest{},
		client.GetHeightRequestHeader(height),
	)
	if err != nil {
		log.Error().Str("module", "mint").Err(err).
			Int64("height", height).
			Msg("error while getting params")
		return
	}

	err = db.SaveMintParams(types.NewMintParams(res.Params, height))
	if err != nil {
		log.Error().Str("module", "mint").Err(err).
			Int64("height", height).
			Msg("error while saving params")
		return
	}
}
