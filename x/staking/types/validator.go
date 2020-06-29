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
}

// NewValidator allows to build a new Validator implementation having the given data
func NewValidator(consAddr sdk.ConsAddress, opAddr sdk.ValAddress, consPubKey crypto.PubKey, description staking.Description) Validator {
	return validator{
		ConsensusAddr: consAddr,
		ConsPubKey:    consPubKey,
		OperatorAddr:  opAddr,
		Description:  description,
		}
	}


// validator allows to easily implement the Validator interface
//unexported
type validator struct {
	ConsensusAddr sdk.ConsAddress
	ConsPubKey    crypto.PubKey
	OperatorAddr  sdk.ValAddress
	Description   staking.Description
}

// GetConsAddr implements the Validator interface
func (v validator) GetConsAddr() sdk.ConsAddress {
	return v.ConsensusAddr
}

// GetConsPubKey implements the Validator interface
func (v validator) GetConsPubKey() crypto.PubKey {
	return v.ConsPubKey
}

// GetOperator implements the Validator interface
func (v validator) GetOperator() sdk.ValAddress {
	return v.OperatorAddr
}

func (v validator) GetMoniker() string {
	return v.Description.Moniker
}
func (v validator) GetIdentity() string {
	return v.Description.Identity
}
func (v validator) GetWebsite() string {
	return v.Description.Website
}
func (v validator) GetSecurityContact() string {
	return v.Description.SecurityContact
}
func (v validator) GetDetails() string {
	return v.Description.Details
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
