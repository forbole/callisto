package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/tendermint/tendermint/crypto"
)

// Validator represents a single validator.
// This is defined as an interface so that we can use the SDK types
// as well as database types properly.
type Validator interface {
	GetConsAddr() string
	GetConsPubKey() crypto.PubKey
	GetOperator() string
	GetSelfDelegateAddress() string
	GetMaxChangeRate() *sdk.Dec
	GetMaxRate() *sdk.Dec
}

// validator allows to easily implement the Validator interface
type validator struct {
	ConsensusAddr       string
	ConsPubKey          crypto.PubKey
	OperatorAddr        string
	SelfDelegateAddress string
	MaxChangeRate       *sdk.Dec
	MaxRate             *sdk.Dec
}

// NewValidator allows to build a new Validator implementation having the given data
func NewValidator(
	consAddr string, opAddr string, consPubKey crypto.PubKey,
	selfDelegateAddress string, maxChangeRate *sdk.Dec,
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
func (v validator) GetConsAddr() string {
	return v.ConsensusAddr
}

// GetConsPubKey implements the Validator interface
func (v validator) GetConsPubKey() crypto.PubKey {
	return v.ConsPubKey
}

func (v validator) GetOperator() string {
	return v.OperatorAddr
}

func (v validator) GetSelfDelegateAddress() string {
	return v.SelfDelegateAddress
}

//Equals return the equality of two validator
func (v validator) Equals(w validator) bool {
	return v.ConsensusAddr == w.ConsensusAddr &&
		v.ConsPubKey.Equals(w.ConsPubKey) &&
		v.OperatorAddr == w.OperatorAddr
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
	OperatorAddress string
	Description     staking.Description
	Height          int64
}

// NewValidatorDescription return a new ValidatorDescription object
func NewValidatorDescription(opAddr string, description staking.Description, height int64,
) ValidatorDescription {
	return ValidatorDescription{
		OperatorAddress: opAddr,
		Description:     description,
		Height:          height,
	}
}

// Equal tells whether v and w contain the same data
func (v ValidatorDescription) Equals(w ValidatorDescription) bool {
	return v.OperatorAddress == w.OperatorAddress &&
		v.Description == w.Description &&
		v.Height == w.Height
}

// _________________________________________________________

// ValidatorUptime contains the uptime information of a single
// validator for a specific height and point in time
type ValidatorUptime struct {
	ValidatorAddress    string
	SignedBlocksWindow  int64
	MissedBlocksCounter int64
	Height              int64
}

// NewValidatorUptime allows to build a new ValidatorUptime instance
func NewValidatorUptime(valAddr string, signedBlocWindow, missedBlocksCounter, height int64) ValidatorUptime {
	return ValidatorUptime{
		ValidatorAddress:    valAddr,
		SignedBlocksWindow:  signedBlocWindow,
		MissedBlocksCounter: missedBlocksCounter,
		Height:              height,
	}
}

// Equal tells whether v and w contain the same data
func (v ValidatorUptime) Equal(w ValidatorUptime) bool {
	return v.ValidatorAddress == w.ValidatorAddress &&
		v.SignedBlocksWindow == w.SignedBlocksWindow &&
		v.MissedBlocksCounter == w.MissedBlocksCounter &&
		v.Height == w.Height
}

// _________________________________________________________

// ValidatorDelegations contains both a validator delegations as
// well as its unbonding delegations
type ValidatorDelegations struct {
	ConsAddress          string
	Delegations          staking.Delegations
	UnbondingDelegations staking.UnbondingDelegations
	Height               int64
}

//-----------------------------------------------------

//ValidatorCommission allow to build a validator commission instance
type ValidatorCommission struct {
	ValAddress        string
	Commission        *sdk.Dec
	MinSelfDelegation *sdk.Int
	Height            int64
}

// NewValidatorCommission return a new validator commission instance
func NewValidatorCommission(
	valAddress string, rate *sdk.Dec, minSelfDelegation *sdk.Int, height int64,
) ValidatorCommission {
	return ValidatorCommission{
		ValAddress:        valAddress,
		Commission:        rate,
		MinSelfDelegation: minSelfDelegation,
		Height:            height,
	}
}

//Equals return the equality of two validatorCommission
func (v ValidatorCommission) Equals(w ValidatorCommission) bool {
	return v.ValAddress == w.ValAddress &&
		v.Commission == w.Commission &&
		v.MinSelfDelegation == w.MinSelfDelegation &&
		v.Height == w.Height
}

//--------------------------------------------

// ValidatorVotingPower represents the voting power of a validator at a specific block height
type ValidatorVotingPower struct {
	ConsensusAddress string
	VotingPower      int64
	Height           int64
}

// NewValidatorVotingPower creates a new ValidatorVotingPower
func NewValidatorVotingPower(address string, votingPower int64, height int64) ValidatorVotingPower {
	return ValidatorVotingPower{
		ConsensusAddress: address,
		VotingPower:      votingPower,
		Height:           height,
	}
}

// Equals tells whether v and w are equals
func (v ValidatorVotingPower) Equals(w ValidatorVotingPower) bool {
	return v.ConsensusAddress == w.ConsensusAddress &&
		v.VotingPower == w.VotingPower &&
		v.Height == w.Height
}

//--------------------------------------------------------
// ValidatorStatus represent status and jailed state for validator in specific height an timestamp
type ValidatorStatus struct {
	ConsensusAddress string
	Status           int
	Jailed           bool
	Height           int64
}

// NewValidatorVotingPower creates a new ValidatorVotingPower
func NewValidatorStatus(address string, status int, jailed bool, height int64) ValidatorStatus {
	return ValidatorStatus{
		ConsensusAddress: address,
		Status:           status,
		Jailed:           jailed,
		Height:           height,
	}
}

// Equals tells whether v and w are equals
func (v ValidatorStatus) Equals(w ValidatorStatus) bool {
	return v.ConsensusAddress == w.ConsensusAddress &&
		v.Jailed == w.Jailed &&
		v.Status == w.Status &&
		v.Height == w.Height
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
	validatorAddress string,
	validatorIndex int,
	signature string,
) DoubleSignVote {
	return DoubleSignVote{
		Type:             roundType,
		Height:           height,
		Round:            round,
		BlockID:          blockID,
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
		w.ValidatorAddress == v.ValidatorAddress &&
		w.ValidatorIndex == v.ValidatorIndex &&
		w.Signature == v.Signature

}
