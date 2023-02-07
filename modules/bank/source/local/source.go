package local

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/types/query"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/forbole/juno/v4/node/local"

	"github.com/forbole/bdjuno/v3/modules/bank/source"
	"github.com/forbole/bdjuno/v3/modules/pricefeed"
	"github.com/forbole/bdjuno/v3/types"
)

var (
	_ source.Source = &Source{}
)

// Source represents the implementation of the bank keeper that works on a local node
type Source struct {
	*local.Source
	q banktypes.QueryServer
}

// NewSource builds a new Source instance
func NewSource(source *local.Source, bk banktypes.QueryServer) *Source {
	return &Source{
		Source: source,
		q:      bk,
	}
}

// GetBalances implements bankkeeper.Source
func (s Source) GetBalances(addresses []string, height int64) ([]types.NativeTokenAmount, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return nil, fmt.Errorf("error while loading height: %s", err)
	}

	var balances []types.NativeTokenAmount
	for _, address := range addresses {
		balRes, err := s.q.Balance(
			sdk.WrapSDKContext(ctx),
			&banktypes.QueryBalanceRequest{
				Address: address,
				Denom:   pricefeed.GetDenom(),
			})
		if err != nil {
			return nil, fmt.Errorf("error while getting balance: %s", err)
		}

		balances = append(balances, types.NewNativeTokenAmount(
			address,
			balRes.Balance.Amount,
			height,
		))
	}

	return balances, nil
}

// GetSupply implements bankkeeper.Source
func (s Source) GetSupply(height int64) (sdk.Coins, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return nil, fmt.Errorf("error while loading height: %s", err)
	}

	var coins []sdk.Coin
	var nextKey []byte
	var stop = false
	for !stop {
		res, err := s.q.TotalSupply(
			sdk.WrapSDKContext(ctx),
			&banktypes.QueryTotalSupplyRequest{
				Pagination: &query.PageRequest{
					Key:   nextKey,
					Limit: 100, // Query 100 supplies at time
				},
			})
		if err != nil {
			return nil, fmt.Errorf("error while getting total supply: %s", err)
		}

		nextKey = res.Pagination.NextKey
		stop = len(res.Pagination.NextKey) == 0
		coins = append(coins, res.Supply...)
	}

	return coins, nil
}

// GetAccountBalances implements bankkeeper.Source
func (s Source) GetAccountBalance(address string, height int64) ([]sdk.Coin, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return nil, fmt.Errorf("error while loading height: %s", err)
	}

	balRes, err := s.q.AllBalances(sdk.WrapSDKContext(ctx), &banktypes.QueryAllBalancesRequest{Address: address})
	if err != nil {
		return nil, fmt.Errorf("error while getting all balances: %s", err)
	}

	return balRes.Balances, nil
}
