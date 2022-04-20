package remote

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/forbole/juno/v3/node/remote"

	bankkeeper "github.com/forbole/bdjuno/v3/modules/bank/source"
	"github.com/forbole/bdjuno/v3/types"
)

var (
	_ bankkeeper.Source = &Source{}
)

type Source struct {
	*remote.Source
	bankClient banktypes.QueryClient
}

// NewSource builds a new Source instance
func NewSource(source *remote.Source, bankClient banktypes.QueryClient) *Source {
	return &Source{
		Source:     source,
		bankClient: bankClient,
	}
}

// GetBalances implements bankkeeper.Source
func (s Source) GetBalances(addresses []string, height int64) ([]types.AccountBalance, error) {
	ctx := remote.GetHeightRequestContext(s.Ctx, height)

	var balances []types.AccountBalance
	for _, address := range addresses {
		balRes, err := s.bankClient.AllBalances(ctx, &banktypes.QueryAllBalancesRequest{Address: address})
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

// GetSupply implements bankkeeper.Source
func (s Source) GetSupply(height int64) (sdk.Coins, error) {
	res, err := s.bankClient.TotalSupply(remote.GetHeightRequestContext(s.Ctx, height), &banktypes.QueryTotalSupplyRequest{})
	if err != nil {
		return nil, fmt.Errorf("error while getting total supply: %s", err)
	}

	return res.Supply, nil
}
