package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Delegation represents a single delegation made from a delegator
// to a specific validator at a specific height (and timestamp)
// containing a given amount of tokens
type Delegation struct {
	DelegatorAddress sdk.AccAddress
	ValidatorAddress sdk.ValAddress
	Amount           sdk.Coin
	Height           int64
	Timestamp        time.Time
}

// NewDelegation creates a new Delegation instance containing
// the given data
func NewDelegation(
	delegator sdk.AccAddress, validatorAddress sdk.ValAddress, amount sdk.Coin,
	height int64, timestamp time.Time,
) Delegation {
	return Delegation{
		DelegatorAddress: delegator,
		ValidatorAddress: validatorAddress,
		Amount:           amount,
		Height:           height,
		Timestamp:        timestamp,
	}
}

// _________________________________________________________

// UnbondingDelegation represents a single unbonding delegation
type UnbondingDelegation struct {
	DelegatorAddress    sdk.AccAddress
	ValidatorAddress    sdk.ValAddress
	Amount              sdk.Coin
	CompletionTimestamp time.Time
	Height              int64
	Timestamp           time.Time
}

// NewUnbondingDelegation allows to create a new UnbondingDelegation instance
func NewUnbondingDelegation(
	delegator sdk.AccAddress, validator sdk.ValAddress, amount sdk.Coin, completionTimestamp time.Time,
	height int64, timestamp time.Time,
) UnbondingDelegation {
	return UnbondingDelegation{
		DelegatorAddress:    delegator,
		ValidatorAddress:    validator,
		Amount:              amount,
		CompletionTimestamp: completionTimestamp,
		Height:              height,
		Timestamp:           timestamp,
	}
}

// _________________________________________________________

// Redelegation represents a single re-delegations
type Redelegation struct {
	DelegatorAddress sdk.AccAddress
	SrcValidator     sdk.ValAddress
	DstValidator     sdk.ValAddress
	Amount           sdk.Coin
	CreationHeight   int64
	CompletionTime   time.Time
}

// NewRedelegation build a new Redelegation object
func NewRedelegation(
	delegator sdk.AccAddress, srcValidator, dstValidator sdk.ValAddress, amount sdk.Coin,
	completionTime time.Time, height int64,
) Redelegation {
	return Redelegation{
		DelegatorAddress: delegator,
		SrcValidator:     srcValidator,
		DstValidator:     dstValidator,
		Amount:           amount,
		CreationHeight:   height,
		CompletionTime:   completionTime,
	}
}

//DelegationShare save the self delegation ratio on that instance
type DelegationShare struct{
	ValidatorAddress sdk.ValAddress
	DelegatorAddress sdk.AccAddress
	Shares           int64
	Height           int64
	Timestamp        time.Time
}

//NewDelegationShare get a new instance of modifly self Delegation
func NewDelegationShare (ValidatorAddress sdk.ValAddress,
	DelegatorAddress sdk.AccAddress,
	Shares           int64,
	Height           int64,
	Timestamp        time.Time) DelegationShare{
		return DelegationShare{
			ValidatorAddress:	ValidatorAddress,
			DelegatorAddress:   DelegatorAddress,
			Shares          :	Shares        ,
			Height          :	Height        ,
			Timestamp       :	Timestamp     ,
		}
	}

