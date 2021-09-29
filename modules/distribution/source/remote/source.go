package remote

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/desmos-labs/juno/node/remote"

	distrsource "github.com/forbole/bdjuno/modules/distribution/source"
)

var (
	_ distrsource.Source = &Source{}
)

// Source implements distrsource.Source querying the data from a remote node
type Source struct {
	*remote.Source
	distrClient distrtypes.QueryClient
}

// NewSource returns a new Source instace
func NewSource(source *remote.Source, distrClient distrtypes.QueryClient) *Source {
	return &Source{
		Source:      source,
		distrClient: distrClient,
	}
}

// ValidatorCommission implements distrsource.Source
func (s Source) ValidatorCommission(valOperAddr string, height int64) (sdk.DecCoins, error) {
	res, err := s.distrClient.ValidatorCommission(
		s.Ctx,
		&distrtypes.QueryValidatorCommissionRequest{ValidatorAddress: valOperAddr},
		remote.GetHeightRequestHeader(height),
	)
	if err != nil {
		return nil, err
	}

	return res.Commission.Commission, nil
}

// DelegatorTotalRewards implements distrsource.Source
func (s Source) DelegatorTotalRewards(delegator string, height int64) ([]distrtypes.DelegationDelegatorReward, error) {
	res, err := s.distrClient.DelegationTotalRewards(
		s.Ctx,
		&distrtypes.QueryDelegationTotalRewardsRequest{DelegatorAddress: delegator},
		remote.GetHeightRequestHeader(height),
	)
	if err != nil {
		return nil, err
	}

	return res.Rewards, nil
}

// DelegatorWithdrawAddress implements distrsource.Source
func (s Source) DelegatorWithdrawAddress(delegator string, height int64) (string, error) {
	res, err := s.distrClient.DelegatorWithdrawAddress(
		s.Ctx,
		&distrtypes.QueryDelegatorWithdrawAddressRequest{DelegatorAddress: delegator},
		remote.GetHeightRequestHeader(height),
	)
	if err != nil {
		return "", err
	}

	return res.WithdrawAddress, nil
}

// CommunityPool implements distrsource.Source
func (s Source) CommunityPool(height int64) (sdk.DecCoins, error) {
	res, err := s.distrClient.CommunityPool(
		s.Ctx,
		&distrtypes.QueryCommunityPoolRequest{},
		remote.GetHeightRequestHeader(height),
	)
	if err != nil {
		return nil, err
	}

	return res.Pool, nil
}

// Params implements distrsource.Source
func (s Source) Params(height int64) (distrtypes.Params, error) {
	res, err := s.distrClient.Params(
		s.Ctx,
		&distrtypes.QueryParamsRequest{},
		remote.GetHeightRequestHeader(height),
	)
	if err != nil {
		return distrtypes.Params{}, err
	}

	return res.Params, nil
}
