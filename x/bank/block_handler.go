package bank

import (
	"context"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/forbole/bdjuno/database"
)

// HandleBlock allows to handle a block properly
func HandleBlock(block *tmctypes.ResultBlock, bankClient banktypes.QueryClient, db *database.BigDipperDb) error {
	err := updateSupply(block.Block.Height, bankClient, db)
	if err != nil {
		return err
	}

	return nil
}

// updateSupply updates the supply of all the tokens for the given height
func updateSupply(height int64, bankClient banktypes.QueryClient, db *database.BigDipperDb) error {
	res, err := bankClient.TotalSupply(context.Background(), &banktypes.QueryTotalSupplyRequest{})
	if err != nil {
		return err
	}

	return db.SaveSupplyToken(res.Supply, height)
}
