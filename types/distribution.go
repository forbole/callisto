package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// ValidatorInfo contains the information about a single validator
type ValidatorInfo struct {
	ValidatorOperAddr         string
	ValidatorSelfDelegateAddr string
}

// NewValidatorInfo returns a new ValidatorInfo instance
func NewValidatorInfo(validatorOperAddr string, validatorSelfDelegateAddr string) ValidatorInfo {
	return ValidatorInfo{
		ValidatorOperAddr:         validatorOperAddr,
		ValidatorSelfDelegateAddr: validatorSelfDelegateAddr,
	}
}

// -------------------------------------------------------------------------------------------------------------------

// ValidatorCommissionAmount represents the commission amount for a specific validator
type ValidatorCommissionAmount struct {
	ValidatorOperAddr         string
	ValidatorSelfDelegateAddr string
	Amount                    []sdk.DecCoin
	Height                    int64
}

// NewValidatorCommissionAmount allows to build a new ValidatorCommissionAmount instance
func NewValidatorCommissionAmount(
	valOperAddr, valSelfDelegateAddress string, amount sdk.DecCoins, height int64,
) ValidatorCommissionAmount {
	return ValidatorCommissionAmount{
		ValidatorOperAddr:         valOperAddr,
		ValidatorSelfDelegateAddr: valSelfDelegateAddress,
		Amount:                    amount,
		Height:                    height,
	}
}

// -------------------------------------------------------------------------------------------------------------------

// DelegatorRewardAmount contains the data of a delegator commission amount
type DelegatorRewardAmount struct {
	ValidatorOperAddr string
	DelegatorAddress  string
	WithdrawAddress   string
	Amount            []sdk.DecCoin
	Height            int64
}

// NewDelegatorRewardAmount allows to build a new DelegatorRewardAmount instance
func NewDelegatorRewardAmount(
	delegator, valOperAddr, withdrawAddress string, amount sdk.DecCoins, height int64,
) DelegatorRewardAmount {
	return DelegatorRewardAmount{
		ValidatorOperAddr: valOperAddr,
		DelegatorAddress:  delegator,
		WithdrawAddress:   withdrawAddress,
		Amount:            amount,
		Height:            height,
	}
}
