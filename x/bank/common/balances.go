package common

import (
	"context"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/x/bank/types"
	"github.com/forbole/bdjuno/x/utils"
	"github.com/rs/zerolog/log"
)

// UpdateBalance updates the balance for the user at the given height
func UpdateBalance(address string, height int64, client banktypes.QueryClient, db *database.BigDipperDb) error {
	// Get the balances
	header := utils.GetHeightRequestHeader(height)
	res, err := client.AllBalances(
		context.Background(),
		&banktypes.QueryAllBalancesRequest{
			Address: address,
		},
		header,
	)
	if err != nil {
		return err
	}

	// Save the balances
	return db.SaveAccountBalances([]types.AccountBalance{
		types.NewAccountBalance(
			address,
			res.Balances,
			height,
		),
	})
}

// RefreshBalance returns a function that when called refreshes the balance of the user having the given address
func RefreshBalance(address string, client banktypes.QueryClient, db *database.BigDipperDb) func() {
	return func() {
		height, err := db.GetLastBlockHeight()
		if err != nil {
			log.Error().Err(err).Str("module", "bank").
				Str("operation", "refresh balance").Msg("error while getting latest block height")
			return
		}

		err = UpdateBalance(address, height, client, db)
		if err != nil {
			log.Error().Err(err).Str("module", "bank").
				Str("operation", "refresh balance").Msg("error while updating balance")
		}
	}
}
