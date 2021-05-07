package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// StakingParams contains all the parameters related to the staking module
type StakingParams struct {
	BondName string
}

// NewStakingParams allows to build a new StakingParams
func NewStakingParams(bondDenom string) StakingParams {
	return StakingParams{
		BondName: bondDenom,
	}
}

// --------------------------------------------------------------------------------------------------------------------

// Validator represents a single validator.
// This is defined as an interface so that we can use the SDK types
// as well as database types properly.
type Validator interface {
	GetConsAddr() string
	GetConsPubKey() string
	GetOperator() string
	GetSelfDelegateAddress() string
	GetMaxChangeRate() *sdk.Dec
	GetMaxRate() *sdk.Dec
	GetHeight() int64
}

// validator allows to easily implement the Validator interface
type validator struct {
	ConsensusAddr       string
	ConsPubKey          string
	OperatorAddr        string
	SelfDelegateAddress string
	MaxChangeRate       *sdk.Dec
	MaxRate             *sdk.Dec
	Height              int64
}

// NewValidator allows to build a new Validator implementation having the given data
func NewValidator(
	consAddr string, opAddr string, consPubKey string,
	selfDelegateAddress string, maxChangeRate *sdk.Dec,
	maxRate *sdk.Dec, height int64,
) Validator {
	return validator{
		ConsensusAddr:       consAddr,
		ConsPubKey:          consPubKey,
		OperatorAddr:        opAddr,
		SelfDelegateAddress: selfDelegateAddress,
		MaxChangeRate:       maxChangeRate,
		MaxRate:             maxRate,
		Height:              height,
	}
}

// GetConsAddr implements the Validator interface
func (v validator) GetConsAddr() string {
	return v.ConsensusAddr
}

// GetConsPubKey implements the Validator interface
func (v validator) GetConsPubKey() string {
	return v.ConsPubKey
}

func (v validator) GetOperator() string {
	return v.OperatorAddr
}

func (v validator) GetSelfDelegateAddress() string {
	return v.SelfDelegateAddress
}

func (v validator) GetMaxChangeRate() *sdk.Dec {
	return v.MaxChangeRate
}

func (v validator) GetMaxRate() *sdk.Dec {
	return v.MaxRate
}

func (v validator) GetHeight() int64 {
	return v.Height
}

// -----------------------------------------------------------------------------------------------------------------

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
