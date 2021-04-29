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
}

// NewDelegation creates a new Delegation instance containing
// the given data
func NewDelegation(delegator string, validatorAddress string, amount sdk.Coin, shares string) Delegation {
	return Delegation{
		DelegatorAddress: delegator,
		ValidatorAddress: validatorAddress,
		Amount:           amount,
		Shares:           shares,
	}
}

// -----------------------------------------------------------------------------------------------------------------

// DelegationUpdateData contains the data needed to update a delegation
type DelegationUpdateData struct {
	Delegator string
	Validator string
	Amount    sdk.Coin
}

func NewDelegationUpdateData(delegator, validator string, amount sdk.Coin) DelegationUpdateData {
	return DelegationUpdateData{
		Delegator: delegator,
		Validator: validator,
		Amount:    amount,
	}
}

// DelegationDeleteData contains the data needed to delete a delegation
type DelegationDeleteData struct {
	Delegator string
	Validator string
}

func NewDelegationDeleteData(delegator, validator string) DelegationDeleteData {
	return DelegationDeleteData{
		Delegator: delegator,
		Validator: validator,
	}
}

// -----------------------------------------------------------------------------------------------------------------

// UnbondingDelegation represents a single unbonding delegation
type UnbondingDelegation struct {
	DelegatorAddress    string
	ValidatorAddress    string
	Amount              sdk.Coin
	CompletionTimestamp time.Time
}

// NewUnbondingDelegation allows to create a new UnbondingDelegation instance
func NewUnbondingDelegation(
	delegator string, validator string, amount sdk.Coin, completionTimestamp time.Time,
) UnbondingDelegation {
	return UnbondingDelegation{
		DelegatorAddress:    delegator,
		ValidatorAddress:    validator,
		Amount:              amount,
		CompletionTimestamp: completionTimestamp,
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
}

// NewRedelegation build a new Redelegation object
func NewRedelegation(
	delegator string, srcValidator, dstValidator string, amount sdk.Coin, completionTime time.Time,
) Redelegation {
	return Redelegation{
		DelegatorAddress: delegator,
		SrcValidator:     srcValidator,
		DstValidator:     dstValidator,
		Amount:           amount,
		CompletionTime:   completionTime,
	}
}
