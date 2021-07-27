package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Delegation represents a single delegation made from a delegator
// to a specific validator at a specific height (and timestamp)
// containing a given amount of tokens
type Delegation struct {
	DelegatorAddress  string
	ValidatorOperAddr string
	Amount            sdk.Coin
	Height            int64
}

// NewDelegation creates a new Delegation instance containing
// the given data
func NewDelegation(delegator string, validatorOperAddr string, amount sdk.Coin, height int64) Delegation {
	return Delegation{
		DelegatorAddress:  delegator,
		ValidatorOperAddr: validatorOperAddr,
		Amount:            amount,
		Height:            height,
	}
}

// -----------------------------------------------------------------------------------------------------------------

// UnbondingDelegation represents a single unbonding delegation
type UnbondingDelegation struct {
	DelegatorAddress    string
	ValidatorOperAddr   string
	Amount              sdk.Coin
	CompletionTimestamp time.Time
	Height              int64
}

// NewUnbondingDelegation allows to create a new UnbondingDelegation instance
func NewUnbondingDelegation(
	delegator string, validatorOperAddr string, amount sdk.Coin, completionTimestamp time.Time, height int64,
) UnbondingDelegation {
	return UnbondingDelegation{
		DelegatorAddress:    delegator,
		ValidatorOperAddr:   validatorOperAddr,
		Amount:              amount,
		CompletionTimestamp: completionTimestamp,
		Height:              height,
	}
}

// Equal returns true iff u and v contain the same data
func (u UnbondingDelegation) Equal(v UnbondingDelegation) bool {
	return u.DelegatorAddress == v.DelegatorAddress &&
		u.ValidatorOperAddr == v.ValidatorOperAddr &&
		u.Amount.IsEqual(v.Amount) &&
		u.CompletionTimestamp.Equal(v.CompletionTimestamp) &&
		u.Height == v.Height
}

// _________________________________________________________

// Redelegation represents a single re-delegations
type Redelegation struct {
	DelegatorAddress string
	SrcValidator     string
	DstValidator     string
	Amount           sdk.Coin
	CompletionTime   time.Time
	Height           int64
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
		Height:           height,
	}
}

// Equal returns true iff r and s contain the same data
func (r Redelegation) Equal(s Redelegation) bool {
	return r.DelegatorAddress == s.DelegatorAddress &&
		r.SrcValidator == s.SrcValidator &&
		r.DstValidator == s.DstValidator &&
		r.Amount.IsEqual(s.Amount) &&
		r.CompletionTime.Equal(s.CompletionTime) &&
		r.Height == s.Height
}
