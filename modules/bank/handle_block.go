package bank

import (
	"context"

	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/modules/utils"

	"github.com/rs/zerolog/log"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// HandleBlock allows to handle a block properly
func HandleBlock(block *tmctypes.ResultBlock, bankClient banktypes.QueryClient, db *database.Db) error {
	err := updateSupply(block.Block.Height, bankClient, db)
	if err != nil {
		log.Error().Str("module", "bank").Int64("height", block.Block.Height).
			Err(err).Msg("error while updating supply")
	}

	return nil
}

// updateSupply updates the supply of all the tokens for the given height
func updateSupply(height int64, bankClient banktypes.QueryClient, db *database.Db) error {
	log.Debug().Str("module", "bank").Int64("height", height).
		Msg("updating supply")

	res, err := bankClient.TotalSupply(
		context.Background(),
		&banktypes.QueryTotalSupplyRequest{},
		utils.GetHeightRequestHeader(height),
	)
	if err != nil {
		return err
	}

	return db.SaveSupply(res.Supply, height)
}
