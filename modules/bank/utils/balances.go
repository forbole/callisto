package utils

import (
	"context"

	"github.com/desmos-labs/juno/client"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/types"
)

// RefreshBalance returns a function that when called refreshes the balance of the user having the given address
func RefreshBalance(address string, client banktypes.QueryClient, db *database.Db) func() {
	return func() {
		height, err := db.GetLastBlockHeight()
		if err != nil {
			log.Error().Err(err).Str("module", "bank").
				Str("operation", "refresh balance").Msg("error while getting latest block height")
			return
		}

		err = UpdateBalances([]string{address}, height, client, db)
		if err != nil {
			log.Error().Err(err).Str("module", "bank").
				Str("operation", "refresh balance").Msg("error while updating balance")
		}
	}
}

// UpdateBalances updates the balances of the accounts having the given addresses,
// taking the data at the provided height
func UpdateBalances(addresses []string, height int64, bankClient banktypes.QueryClient, db *database.Db) error {
	log.Debug().Str("module", "bank").Int64("height", height).Msg("updating balances")
	header := client.GetHeightRequestHeader(height)

	var balances []types.AccountBalance
	for _, address := range addresses {
		balRes, err := bankClient.AllBalances(
			context.Background(),
			&banktypes.QueryAllBalancesRequest{Address: address},
			header,
		)
		if err != nil {
			return err
		}

		balances = append(balances, types.NewAccountBalance(
			address,
			balRes.Balances,
			height,
		))
	}

	return db.SaveAccountBalances(balances)
}
