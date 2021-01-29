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

// DelegatorCommissionAmount contains the data of a delegator commission amount
type DelegatorCommissionAmount struct {
	ValidatorAddress string
	DelegatorAddress string
	Amount           []sdk.DecCoin
}

// NewDelegatorCommissionAmount allows to build a new DelegatorCommissionAmount instance
func NewDelegatorCommissionAmount(validator, delegator string, amount sdk.DecCoins) DelegatorCommissionAmount {
	return DelegatorCommissionAmount{
		ValidatorAddress: validator,
		DelegatorAddress: delegator,
		Amount:           amount,
	}
}
