package utils

import (
	"context"

	"github.com/desmos-labs/juno/client"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/types"
)

// RefreshBalances updates the balances of the accounts having the given addresses,
// taking the data at the provided height
func RefreshBalances(height int64, addresses []string, bankClient banktypes.QueryClient, db *database.Db) error {
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

// GetUserBalance returns the current balance for the specified user
func GetUserBalance(height int64, address string, bankClient banktypes.QueryClient) (types.AccountBalance, error) {
	header := client.GetHeightRequestHeader(height)
	balRes, err := bankClient.AllBalances(
		context.Background(),
		&banktypes.QueryAllBalancesRequest{Address: address},
		header,
	)
	if err != nil {
		return types.AccountBalance{}, err
	}

	return types.NewAccountBalance(
		address,
		balRes.Balances,
		height,
	), nil
}
