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

// validator allows to easily implement the Validator interface
type validator struct {
	ConsensusAddr sdk.ConsAddress
	ConsPubKey    crypto.PubKey
	OperatorAddr  sdk.ValAddress
}

func (v validator) GetConsAddr() sdk.ConsAddress {
	return v.ConsensusAddr
}

func (v validator) GetConsPubKey() crypto.PubKey {
	return v.ConsPubKey
}

func (v validator) GetOperator() sdk.ValAddress {
	return v.OperatorAddr
}

func NewValidator(consAddr sdk.ConsAddress, opAddr sdk.ValAddress, consPubKey crypto.PubKey) Validator {
	return validator{
		ConsensusAddr: consAddr,
		ConsPubKey:    consPubKey,
		OperatorAddr:  opAddr,
	}
}

// ValidatorUptime contains the uptime information of a single
// validator for a specific height and point in time
type ValidatorUptime struct {
	ValidatorAddress    sdk.ConsAddress
	Height              int64
	SignedBlocksWindow  int64
	MissedBlocksCounter int64
}

// ValidatorDelegations contains both a validator delegations as
// well as its unbonding delegations
type ValidatorDelegations struct {
	ConsAddress          sdk.ConsAddress
	Delegations          staking.Delegations
	UnbondingDelegations staking.UnbondingDelegations
	Height               int64
	Timestamp            time.Time
}
