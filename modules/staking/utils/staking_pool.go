package utils

import (
	"context"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/desmos-labs/juno/client"

	"github.com/forbole/bdjuno/types"
)

func GetStakingPool(height int64, stakingClient stakingtypes.QueryClient) (*types.Pool, error) {
	res, err := stakingClient.Pool(
		context.Background(),
		&stakingtypes.QueryPoolRequest{},
		client.GetHeightRequestHeader(height),
	)
	if err != nil {
		return nil, err
	}

	return types.NewPool(res.Pool.BondedTokens, res.Pool.NotBondedTokens, height), nil
}
