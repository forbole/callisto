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
	GetSelfDelegateAddress() sdk.AccAddress
	GetMaxChangeRate() *sdk.Dec
	GetMaxRate() *sdk.Dec
}

// validator allows to easily implement the Validator interface
type validator struct {
	ConsensusAddr       sdk.ConsAddress
	ConsPubKey          crypto.PubKey
	OperatorAddr        sdk.ValAddress
	SelfDelegateAddress sdk.AccAddress
	MaxChangeRate       *sdk.Dec
	MaxRate             *sdk.Dec
}

// NewValidator allows to build a new Validator implementation having the given data
func NewValidator(
	consAddr sdk.ConsAddress, opAddr sdk.ValAddress, consPubKey crypto.PubKey,
	selfDelegateAddress sdk.AccAddress, maxChangeRate *sdk.Dec,
	maxRate *sdk.Dec,
) Validator {
	return validator{
		ConsensusAddr:       consAddr,
		ConsPubKey:          consPubKey,
		OperatorAddr:        opAddr,
		SelfDelegateAddress: selfDelegateAddress,
		MaxChangeRate:       maxChangeRate,
		MaxRate:             maxRate,
	}
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

func (v validator) GetSelfDelegateAddress() sdk.AccAddress {
	return v.SelfDelegateAddress
}

//Equals return the equality of two validator
func (v validator) Equals(w validator) bool {
	return v.ConsensusAddr.Equals(w.ConsensusAddr) &&
		v.ConsPubKey.Equals(w.ConsPubKey) &&
		v.OperatorAddr.Equals(w.OperatorAddr)
}

func (v validator) GetMaxChangeRate() *sdk.Dec {
	return v.MaxChangeRate
}

func (v validator) GetMaxRate() *sdk.Dec {
	return v.MaxRate
}

// _________________________________________________________

// ValidatorDescription contains the description of a validator
// and timestamp do the description get changed
type ValidatorDescription struct {
	OperatorAddress sdk.ValAddress
	Description     staking.Description
	Timestamp       time.Time
	Height          int64
}

// NewValidatorDescription return a new ValidatorDescription object
func NewValidatorDescription(opAddr sdk.ValAddress, description staking.Description, height int64, timestamp time.Time,
) ValidatorDescription {
	return ValidatorDescription{
		OperatorAddress: opAddr,
		Description:     description,
		Timestamp:       timestamp,
		Height:          height,
	}
}

// Equal tells whether v and w contain the same data
func (v ValidatorDescription) Equals(w ValidatorDescription) bool {
	return v.OperatorAddress.Equals(w.OperatorAddress) &&
		v.Description == w.Description &&
		v.Timestamp.Equal(w.Timestamp) &&
		v.Height == w.Height
}

// _________________________________________________________

// ValidatorUptime contains the uptime information of a single
// validator for a specific height and point in time
type ValidatorUptime struct {
	ValidatorAddress    sdk.ConsAddress
	SignedBlocksWindow  int64
	MissedBlocksCounter int64
	Height              int64
	Timestamp           time.Time
}

// NewValidatorUptime allows to build a new ValidatorUptime instance
func NewValidatorUptime(valAddr sdk.ConsAddress, signedBlocWindow, missedBlocksCounter, height int64, timestamp time.Time) ValidatorUptime {
	return ValidatorUptime{
		ValidatorAddress:    valAddr,
		SignedBlocksWindow:  signedBlocWindow,
		MissedBlocksCounter: missedBlocksCounter,
		Height:              height,
		Timestamp:           timestamp,
	}
}

// Equal tells whether v and w contain the same data
func (v ValidatorUptime) Equal(w ValidatorUptime) bool {
	return v.ValidatorAddress.Equals(w.ValidatorAddress) &&
		v.SignedBlocksWindow == w.SignedBlocksWindow &&
		v.MissedBlocksCounter == w.MissedBlocksCounter &&
		v.Height == w.Height &&
		v.Timestamp.Equal(w.Timestamp)
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
type ValidatorCommission struct {
	ValAddress        sdk.ValAddress
	Commission        *sdk.Dec
	MinSelfDelegation *sdk.Int
	Height            int64
	Timestamp         time.Time
}

// NewValidatorCommission return a new validator commission instance
func NewValidatorCommission(
	valAddress sdk.ValAddress, rate *sdk.Dec, minSelfDelegation *sdk.Int, height int64, timestamp time.Time,
) ValidatorCommission {
	return ValidatorCommission{
		ValAddress:        valAddress,
		Commission:        rate,
		MinSelfDelegation: minSelfDelegation,
		Height:            height,
		Timestamp:         timestamp,
	}
}

//Equals return the equality of two validatorCommission
func (v ValidatorCommission) Equals(w ValidatorCommission) bool {
	return v.ValAddress.Equals(w.ValAddress) &&
		v.Commission == w.Commission &&
		v.MinSelfDelegation == w.MinSelfDelegation &&
		v.Height == w.Height &&
		v.Timestamp.Equal(w.Timestamp)
}

//--------------------------------------------

// ValidatorVotingPower represents the voting power of a validator at a specific block height
type ValidatorVotingPower struct {
	ConsensusAddress sdk.ConsAddress
	VotingPower      int64
	Height           int64
	Timestamp        time.Time
}

// NewValidatorVotingPower creates a new ValidatorVotingPower
func NewValidatorVotingPower(address sdk.ConsAddress, votingPower int64, height int64, timestamp time.Time) ValidatorVotingPower {
	return ValidatorVotingPower{
		ConsensusAddress: address,
		VotingPower:      votingPower,
		Height:           height,
		Timestamp:        timestamp,
	}
}

// Equals tells whether v and w are equals
func (v ValidatorVotingPower) Equals(w ValidatorVotingPower) bool {
	return v.ConsensusAddress.Equals(w.ConsensusAddress) &&
		v.VotingPower == w.VotingPower &&
		v.Height == w.Height &&
		v.Timestamp.Equal(w.Timestamp)
}

//--------------------------------------------------------
// ValidatorStatus represent status and jailed state for validator in specific height an timestamp
type ValidatorStatus struct {
	ConsensusAddress sdk.ConsAddress
	Status           int
	Jailed           bool
	Height           int64
	Timestamp        time.Time
}

// NewValidatorVotingPower creates a new ValidatorVotingPower
func NewValidatorStatus(address sdk.ConsAddress, status int, jailed bool, height int64, timestamp time.Time) ValidatorStatus {
	return ValidatorStatus{
		ConsensusAddress: address,
		Status:           status,
		Jailed:           jailed,
		Height:           height,
		Timestamp:        timestamp,
	}
}

// Equals tells whether v and w are equals
func (v ValidatorStatus) Equals(w ValidatorStatus) bool {
	return v.ConsensusAddress.Equals(w.ConsensusAddress) &&
		v.Jailed == w.Jailed &&
		v.Status == w.Status &&
		v.Height == w.Height &&
		v.Timestamp.Equal(w.Timestamp)
}

//---------------------------------------------------------------

// DoubleSignEvidence represent a double sign evidence on each tendermint block
type DoubleSignEvidence struct {
	Pubkey string
	VoteA  DoubleSignVote
	VoteB  DoubleSignVote
}

// NewDoubleSignEvidence return a new DoubleSignEvidence object
func NewDoubleSignEvidence(
	pubkey string,
	voteA DoubleSignVote,
	voteB DoubleSignVote,
) DoubleSignEvidence {
	return DoubleSignEvidence{
		Pubkey: pubkey,
		VoteA:  voteA,
		VoteB:  voteB,
	}
}

// Equals tells whether v and w contain the same data
func (w DoubleSignEvidence) Equals(v DoubleSignEvidence) bool {
	return w.Pubkey == v.Pubkey &&
		w.VoteA.Equals(v.VoteA) &&
		w.VoteB.Equals(v.VoteB)
}

// DoubleSignVote represents a double vote which is included inside a DoubleSignEvidence
type DoubleSignVote struct {
	Type             int
	Height           int64
	Round            int
	BlockID          string
	Timestamp        time.Time
	ValidatorAddress string
	ValidatorIndex   int
	Signature        string
}

// NewDoubleSignVote allows to create a new DoubleSignVote instance
func NewDoubleSignVote(
	roundType int,
	height int64,
	round int,
	blockID string,
	timestamp time.Time,
	validatorAddress string,
	validatorIndex int,
	signature string,
) DoubleSignVote {
	return DoubleSignVote{
		Type:             roundType,
		Height:           height,
		Round:            round,
		BlockID:          blockID,
		Timestamp:        timestamp,
		ValidatorAddress: validatorAddress,
		ValidatorIndex:   validatorIndex,
		Signature:        signature,
	}
}

func (w DoubleSignVote) Equals(v DoubleSignVote) bool {
	return w.Type == v.Type &&
		w.Height == v.Height &&
		w.Round == v.Round &&
		w.BlockID == v.BlockID &&
		w.Timestamp.Equal(v.Timestamp) &&
		w.ValidatorAddress == v.ValidatorAddress &&
		w.ValidatorIndex == v.ValidatorIndex &&
		w.Signature == v.Signature

}
