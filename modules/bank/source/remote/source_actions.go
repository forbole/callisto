package remote

import (
	"fmt"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/forbole/juno/v2/node/remote"

	"github.com/forbole/bdjuno/v2/types"
)

// GetBalances implements bankkeeper.Source
func (s Source) GetBalances(addresses []string, height int64) ([]types.AccountBalance, error) {
	header := remote.GetHeightRequestHeader(height)

	var balances []types.AccountBalance
	for _, address := range addresses {
		balRes, err := s.bankClient.AllBalances(s.Ctx, &banktypes.QueryAllBalancesRequest{Address: address}, header)
		if err != nil {
			return nil, fmt.Errorf("error while getting all balances: %s", err)
		}

		balances = append(balances, types.NewAccountBalance(
			address,
			balRes.Balances,
			height,
		))
	}

	return balances, nil
}
