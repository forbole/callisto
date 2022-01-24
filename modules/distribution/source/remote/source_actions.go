package remote

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/forbole/juno/v2/node/remote"
)

// DelegatorTotalRewards implements distrsource.Source
func (s Source) DelegatorTotalRewards(delegator string, height int64) ([]distrtypes.DelegationDelegatorReward, error) {
	res, err := s.distrClient.DelegationTotalRewards(
		s.Ctx,
		&distrtypes.QueryDelegationTotalRewardsRequest{DelegatorAddress: delegator},
		remote.GetHeightRequestHeader(height),
	)
	if err != nil {
		return nil, fmt.Errorf("error while getting delegation total rewards for for delegator %s at height %v: %s", delegator, height, err)
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
