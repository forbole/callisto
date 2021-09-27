package local

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/desmos-labs/juno/node/local"
	"github.com/forbole/bdjuno/modules/bank/source"
	"github.com/forbole/bdjuno/types"
)

var (
	_ source.Source = &Source{}
)

// Source represents the implementation of the bank keeper that works on a local node
type Source struct {
	*local.Source
	bankKeeper bankkeeper.BaseKeeper
}

// NewSource builds a new Source instance
func NewSource(source *local.Source, bk bankkeeper.BaseKeeper) *Source {
	return &Source{
		Source:     source,
		bankKeeper: bk,
	}
}

// GetBalances implements keeper.Source
func (k Source) GetBalances(addresses []string, height int64) ([]types.AccountBalance, error) {
	ctx, err := k.LoadHeight(height)
	if err != nil {
		return nil, fmt.Errorf("error while loading height: %s", err)
	}

	var balances []types.AccountBalance
	for _, address := range addresses {
		addr, err := sdk.AccAddressFromBech32(address)
		if err != nil {
			return nil, err
		}

		balance := k.bankKeeper.GetAllBalances(ctx, addr)
		balances = append(balances, types.NewAccountBalance(address, balance, height))
	}

	return balances, nil
}

// GetSupply implements keeper.Source
func (k Source) GetSupply(height int64) (sdk.Coins, error) {
	ctx, err := k.LoadHeight(height)
	if err != nil {
		return nil, fmt.Errorf("error while loading height: %s", err)
	}

	return k.bankKeeper.GetSupply(ctx).GetTotal(), nil
}
