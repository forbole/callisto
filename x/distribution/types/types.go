package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// ValidatorCommissionAmount represents the commission amount for a specific validator
type ValidatorCommissionAmount struct {
	ValidatorAddress string
	Amount           []sdk.DecCoin
}

// NewValidatorCommissionAmount allows to build a new ValidatorCommissionAmount instance
func NewValidatorCommissionAmount(address string, amount sdk.DecCoins) ValidatorCommissionAmount {
	return ValidatorCommissionAmount{
		ValidatorAddress: address,
		Amount:           amount,
	}
}

// DelegatorRewardAmount contains the data of a delegator commission amount
type DelegatorRewardAmount struct {
	ValidatorAddress string
	DelegatorAddress string
	Amount           []sdk.DecCoin
}

// NewDelegatorRewardAmount allows to build a new DelegatorRewardAmount instance
func NewDelegatorRewardAmount(validator, delegator string, amount sdk.DecCoins) DelegatorRewardAmount {
	return DelegatorRewardAmount{
		ValidatorAddress: validator,
		DelegatorAddress: delegator,
		Amount:           amount,
	}
}
