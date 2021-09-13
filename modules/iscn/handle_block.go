package iscn

import (
	"context"

	"github.com/desmos-labs/juno/client"
	iscntypes "github.com/likecoin/likechain/x/iscn/types"
	"github.com/rs/zerolog/log"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/types"
)

// HandleBlock represents a method that is called each time a new block is created
func HandleBlock(block *tmctypes.ResultBlock, iscnClient iscntypes.QueryClient, db *database.Db) error {
	// Update the params
	go updateParams(block.Block.Height, iscnClient, db)
	return nil
}

// updateParams gets the updated iscn params and stores them inside the database
func updateParams(height int64, iscnClient iscntypes.QueryClient, db *database.Db) {
	log.Debug().Str("module", "iscn").Int64("height", height).
		Msg("updating iscn params")

	res, err := iscnClient.Params(
		context.Background(),
		&iscntypes.QueryParamsRequest{},
		client.GetHeightRequestHeader(height),
	)
	if err != nil {
		log.Error().Str("module", "iscn").Err(err).
			Int64("height", height).
			Msg("error while getting iscn params")
		return
	}

	err = db.SaveIscnParams(types.NewIscnParams(res.Params, height))
	if err != nil {
		log.Error().Str("module", "iscn").Err(err).
			Int64("height", height).
			Msg("error while saving iscn params")
		return
	}
}
