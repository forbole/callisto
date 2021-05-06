package bank

import (
	"context"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/rs/zerolog/log"

	utils2 "github.com/forbole/bdjuno/modules/common/utils"
	"github.com/forbole/bdjuno/types"
)

// UpdateBalances updates the balances of the accounts having the given addresses,
// taking the data at the provided height
func UpdateBalances(addresses []string, height int64, client banktypes.QueryClient, db DB) error {
	log.Debug().Str("module", "bank").Msg("updating balances")
	header := utils2.GetHeightRequestHeader(height)

	var balances []types.AccountBalance

	for _, address := range addresses {
		balRes, err := client.AllBalances(
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
