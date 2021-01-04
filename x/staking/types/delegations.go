package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Delegation represents a single delegation made from a delegator
// to a specific validator at a specific height (and timestamp)
// containing a given amount of tokens
type Delegation struct {
	DelegatorAddress string
	ValidatorAddress string
	Amount           sdk.Coin
	Shares           string
	Height           int64
}

// NewDelegation creates a new Delegation instance containing
// the given data
func NewDelegation(
	delegator string, validatorAddress string, amount sdk.Coin, shares string, height int64,
) Delegation {
	return Delegation{
		DelegatorAddress: delegator,
		ValidatorAddress: validatorAddress,
		Amount:           amount,
		Shares:           shares,
		Height:           height,
	}
}

// _________________________________________________________

// UnbondingDelegation represents a single unbonding delegation
type UnbondingDelegation struct {
	DelegatorAddress    string
	ValidatorAddress    string
	Amount              sdk.Coin
	CompletionTimestamp time.Time
	Height              int64
}

// NewUnbondingDelegation allows to create a new UnbondingDelegation instance
func NewUnbondingDelegation(
	delegator string, validator string, amount sdk.Coin, completionTimestamp time.Time,
	height int64,
) UnbondingDelegation {
	return UnbondingDelegation{
		DelegatorAddress:    delegator,
		ValidatorAddress:    validator,
		Amount:              amount,
		CompletionTimestamp: completionTimestamp,
		Height:              height,
	}
}

// _________________________________________________________

// Redelegation represents a single re-delegations
type Redelegation struct {
	DelegatorAddress string
	SrcValidator     string
	DstValidator     string
	Amount           sdk.Coin
	CompletionTime   time.Time
	CreationHeight   int64
}

// NewRedelegation build a new Redelegation object
func NewRedelegation(
	delegator string, srcValidator, dstValidator string, amount sdk.Coin, completionTime time.Time, height int64,
) Redelegation {
	return Redelegation{
		DelegatorAddress: delegator,
		SrcValidator:     srcValidator,
		DstValidator:     dstValidator,
		Amount:           amount,
		CompletionTime:   completionTime,
		CreationHeight:   height,
	}
}

//DelegationShare save the self delegation ratio on that instance
type DelegationShare struct {
	ValidatorAddress string
	DelegatorAddress string
	Shares           float64
	Height           int64
	Timestamp        time.Time
}

//NewDelegationShare get a new instance of modify self Delegation
func NewDelegationShare(
	validatorAddress string, delegatorAddress string, shares float64,
	height int64, timestamp time.Time,
) DelegationShare {
	return DelegationShare{
		ValidatorAddress: validatorAddress,
		DelegatorAddress: delegatorAddress,
		Shares:           shares,
		Height:           height,
		Timestamp:        timestamp,
	}
}
