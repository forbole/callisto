package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// ValidatorCommissionAmount represents the commission amount for a specific validator
type ValidatorCommissionAmount struct {
	ValidatorConsAddr string
	Amount            []sdk.DecCoin
	Height            int64
}

// NewValidatorCommissionAmount allows to build a new ValidatorCommissionAmount instance
func NewValidatorCommissionAmount(valConsAddr string, amount sdk.DecCoins, height int64) ValidatorCommissionAmount {
	return ValidatorCommissionAmount{
		ValidatorConsAddr: valConsAddr,
		Amount:            amount,
		Height:            height,
	}
}

// -------------------------------------------------------------------------------------------------------------------

// DelegatorRewardAmount contains the data of a delegator commission amount
type DelegatorRewardAmount struct {
	ValidatorConsAddr string
	DelegatorAddress  string
	WithdrawAddress   string
	Amount            []sdk.DecCoin
	Height            int64
}

// NewDelegatorRewardAmount allows to build a new DelegatorRewardAmount instance
func NewDelegatorRewardAmount(
	delegator, valConsAddr, withdrawAddress string, amount sdk.DecCoins, height int64,
) DelegatorRewardAmount {
	return DelegatorRewardAmount{
		ValidatorConsAddr: valConsAddr,
		DelegatorAddress:  delegator,
		WithdrawAddress:   withdrawAddress,
		Amount:            amount,
		Height:            height,
	}
}
