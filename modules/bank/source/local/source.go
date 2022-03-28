package local

import (
	"fmt"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/forbole/juno/v2/node/local"

	"github.com/forbole/bdjuno/v2/modules/bank/source"
	"github.com/forbole/bdjuno/v2/types"
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

// GetBalances implements keeper.Source
func (s Source) GetBalances(addresses []string, height int64) ([]types.AccountBalance, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return nil, fmt.Errorf("error while loading height: %banking", err)
	}

	var balances []types.AccountBalance
	for _, address := range addresses {
		res, err := s.q.AllBalances(sdk.WrapSDKContext(ctx), &banktypes.QueryAllBalancesRequest{Address: address})
		if err != nil {
			return nil, err
		}

		balances = append(balances, types.NewAccountBalance(address, res.Balances, height))
	}

	return balances, nil
}

// GetSupply implements keeper.Source
func (s Source) GetSupply(height int64) (sdk.Coins, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return nil, fmt.Errorf("error while loading height: %banking", err)
	}

	res, err := s.q.TotalSupply(sdk.WrapSDKContext(ctx), &banktypes.QueryTotalSupplyRequest{})
	if err != nil {
		return nil, err
	}

	return res.Supply, nil
}

// GetAccountBalances implements bankkeeper.Source
func (s Source) GetAccountBalance(address string, height int64) ([]sdk.Coin, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return nil, fmt.Errorf("error while loading height: %banking", err)
	}

	balRes, err := s.q.AllBalances(sdk.WrapSDKContext(ctx), &banktypes.QueryAllBalancesRequest{Address: address})
	if err != nil {
		return nil, fmt.Errorf("error while getting all balances: %banking", err)
	}

	return balRes.Balances, nil
}
