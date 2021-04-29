package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// ValidatorCommissionAmount represents the commission amount for a specific validator
type ValidatorCommissionAmount struct {
	ValidatorConsAddress string
	Amount               []sdk.DecCoin
}

// NewValidatorCommissionAmount allows to build a new ValidatorCommissionAmount instance
func NewValidatorCommissionAmount(valConsAddr string, amount sdk.DecCoins) ValidatorCommissionAmount {
	return ValidatorCommissionAmount{
		ValidatorConsAddress: valConsAddr,
		Amount:               amount,
	}
}

// --------------------------------------------------------------------------------------------------------------------

// DelegatorReward contains the data of a delegator commission amount
type DelegatorReward struct {
	ValidatorConsAddress string
	DelegatorAddress     string
	WithdrawAddress      string
	Amount               []sdk.DecCoin
}

// NewDelegatorRewardAmount allows to build a new DelegatorReward instance
func NewDelegatorRewardAmount(
	valConsAddr, delegator, withdrawAddress string, amount sdk.DecCoins,
) DelegatorReward {
	return DelegatorReward{
		ValidatorConsAddress: valConsAddr,
		DelegatorAddress:     delegator,
		WithdrawAddress:      withdrawAddress,
		Amount:               amount,
	}
}
