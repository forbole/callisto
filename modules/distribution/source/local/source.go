package local

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/forbole/juno/v3/node/local"

	distrsource "github.com/forbole/bdjuno/v3/modules/distribution/source"
)

var (
	_ distrsource.Source = &Source{}
)

// Source implements distrsource.Source reading the data from a local node
type Source struct {
	*local.Source
	q distrtypes.QueryServer
}

func NewSource(source *local.Source, keeper distrtypes.QueryServer) *Source {
	return &Source{
		Source: source,
		q:      keeper,
	}
}

// ValidatorCommission implements distrsource.Source
func (s Source) ValidatorCommission(valOperAddr string, height int64) (sdk.DecCoins, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return nil, fmt.Errorf("error while loading height: %s", err)
	}

	res, err := s.q.ValidatorCommission(
		sdk.WrapSDKContext(ctx),
		&distrtypes.QueryValidatorCommissionRequest{ValidatorAddress: valOperAddr},
	)
	if err != nil {
		return nil, err
	}

	return res.Commission.Commission, nil
}

// DelegatorTotalRewards implements distrsource.Source
func (s Source) DelegatorTotalRewards(delegator string, height int64) ([]distrtypes.DelegationDelegatorReward, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return nil, fmt.Errorf("error while loading height: %s", err)
	}

	res, err := s.q.DelegationTotalRewards(
		sdk.WrapSDKContext(ctx),
		&distrtypes.QueryDelegationTotalRewardsRequest{DelegatorAddress: delegator},
	)
	if err != nil {
		return nil, err
	}

	return res.Rewards, nil
}

// DelegatorWithdrawAddress implements distrsource.Source
func (s Source) DelegatorWithdrawAddress(delegator string, height int64) (string, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return "", fmt.Errorf("error while loading height: %s", err)
	}

	res, err := s.q.DelegatorWithdrawAddress(
		sdk.WrapSDKContext(ctx),
		&distrtypes.QueryDelegatorWithdrawAddressRequest{DelegatorAddress: delegator},
	)
	if err != nil {
		return "", err
	}

	return res.WithdrawAddress, nil
}

// CommunityPool implements distrsource.Source
func (s Source) CommunityPool(height int64) (sdk.DecCoins, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return nil, fmt.Errorf("error while loading height: %s", err)
	}

	res, err := s.q.CommunityPool(sdk.WrapSDKContext(ctx), &distrtypes.QueryCommunityPoolRequest{})
	if err != nil {
		return nil, err
	}

	return res.Pool, nil
}

// Params implements distrsource.Source
func (s Source) Params(height int64) (distrtypes.Params, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return distrtypes.Params{}, fmt.Errorf("error while loading height: %s", err)
	}

	res, err := s.q.Params(sdk.WrapSDKContext(ctx), &distrtypes.QueryParamsRequest{})
	if err != nil {
		return distrtypes.Params{}, err
	}

	return res.Params, nil
}
