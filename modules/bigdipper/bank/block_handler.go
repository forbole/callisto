package bank

import (
	"context"

	"github.com/rs/zerolog/log"

	bigdipperdb "github.com/forbole/bdjuno/database/bigdipper"
	utils2 "github.com/forbole/bdjuno/modules/common/utils"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// HandleBlock allows to handle a block properly
func HandleBlock(block *tmctypes.ResultBlock, bankClient banktypes.QueryClient, db *bigdipperdb.Db) error {
	err := updateSupply(block.Block.Height, bankClient, db)
	if err != nil {
		log.Error().Str("module", "bank").Int64("height", block.Block.Height).
			Err(err).Msg("error while updating supply")
	}

	return nil
}

// updateSupply updates the supply of all the tokens for the given height
func updateSupply(height int64, bankClient banktypes.QueryClient, db *bigdipperdb.Db) error {
	log.Debug().Str("module", "bank").Int64("height", height).
		Msg("updating supply")

	res, err := bankClient.TotalSupply(
		context.Background(),
		&banktypes.QueryTotalSupplyRequest{},
		utils2.GetHeightRequestHeader(height),
	)
	if err != nil {
		return err
	}

	return db.SaveSupply(res.Supply, height)
}
