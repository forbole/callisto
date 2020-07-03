package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/tendermint/tendermint/crypto"
)

// Validator represents a single validator.
// This is defined as an interface so that we can use the SDK types
// as well as database types properly.
type Validator interface {
	GetConsAddr() sdk.ConsAddress
	GetConsPubKey() crypto.PubKey
	GetOperator() sdk.ValAddress
	GetDescription() staking.Description
}

// NewValidator allows to build a new Validator implementation having the given data
func NewValidator(consAddr sdk.ConsAddress, opAddr sdk.ValAddress, consPubKey crypto.PubKey,description staking.Description) Validator {
	return validator{
		ConsensusAddr: consAddr,
		ConsPubKey:    consPubKey,
		OperatorAddr:  opAddr,
		Description:   description,
	}
}

// validator allows to easily implement the Validator interface
//unexported
type validator struct {
	ConsensusAddr sdk.ConsAddress
	ConsPubKey    crypto.PubKey
	OperatorAddr  sdk.ValAddress
	Description    staking.Description
}

// GetConsAddr implements the Validator interface
func (v validator) GetConsAddr() sdk.ConsAddress {
	return v.ConsensusAddr
}

// GetConsPubKey implements the Validator interface
func (v validator) GetConsPubKey() crypto.PubKey {
	return v.ConsPubKey
}

func (v validator) GetOperator() sdk.ValAddress {
	return v.OperatorAddr
}

func (v validator) GetDescription() staking.Description{
	return v.Description
}

//Equals return the equality of two validator
func (v validator) Equals(w validator) bool{
	return v.ConsensusAddr.Equals(w.ConsensusAddr)&&
	v.ConsPubKey.Equals(w.ConsPubKey)&&
	v.OperatorAddr.Equals(w.OperatorAddr)&&
	v.Description==w.Description
} 

// _________________________________________________________

// ValidatorUptime contains the uptime information of a single
// validator for a specific height and point in time
type ValidatorUptime struct {
	ValidatorAddress    sdk.ConsAddress
	SignedBlocksWindow  int64
	MissedBlocksCounter int64
	Height              int64
}

// NewValidatorUptime allows to build a new ValidatorUptime instance
func NewValidatorUptime(valAddr sdk.ConsAddress, signedBlocWindow, missedBlocksCounter, height int64) ValidatorUptime {
	return ValidatorUptime{
		ValidatorAddress:    valAddr,
		SignedBlocksWindow:  signedBlocWindow,
		MissedBlocksCounter: missedBlocksCounter,
		Height:              height,
	}
}

// Equal tells whether v and w represent the same uptime
func (v ValidatorUptime) Equal(w ValidatorUptime) bool {
	return v.ValidatorAddress.Equals(w.ValidatorAddress) &&
		v.SignedBlocksWindow == w.SignedBlocksWindow &&
		v.MissedBlocksCounter == w.MissedBlocksCounter &&
		v.Height == w.Height
}

// _________________________________________________________

// ValidatorDelegations contains both a validator delegations as
// well as its unbonding delegations
type ValidatorDelegations struct {
	ConsAddress          sdk.ConsAddress
	Delegations          staking.Delegations
	UnbondingDelegations staking.UnbondingDelegations
	Height               int64
	Timestamp            time.Time
}

//-----------------------------------------------------

//ValidatorCommission allow to build a validator commission instance
type ValidatorCommission struct{
	ValAddress sdk.ValAddress
	Commission 			int64
	MinSelfDelegation int64
	Height               int64
	Timestamp            time.Time
}

//NewValidatorCommission return a new validator commission instance
func NewValidatorCommission(ValAddress sdk.ValAddress,Rate int64,
	MinSelfDelegation int64,Height int64,Timestamp time.Time)(ValidatorCommission){
	return ValidatorCommission{
		ValAddress    :ValAddress,
		Commission 		:Rate,
		MinSelfDelegation:MinSelfDelegation,
        Height        :Height,
        Timestamp     :Timestamp,
	}
}

//Equals return the equality of two validatorCommission
func (v ValidatorCommission) Equals(w ValidatorCommission) bool{
	return v.ValAddress.Equals(w.ValAddress) &&
	v.Commission == w.Commission&&
	v.MinSelfDelegation == w.MinSelfDelegation&&
	v.Height == w.Height&&
	v.Timestamp.Equal(w.Timestamp)
}

//--------------------------------------------

