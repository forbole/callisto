package types

import (
	"database/sql"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ValidatorData contains all the data of a single validator.
// It implements types.Validator interface
type ValidatorData struct {
	ConsAddress         string `db:"consensus_address"`
	ValAddress          string `db:"operator_address"`
	ConsPubKey          string `db:"consensus_pubkey"`
	SelfDelegateAddress string `db:"self_delegate_address"`
	MaxRate             string `db:"max_rate"`
	MaxChangeRate       string `db:"max_change_rate"`
	Height              int64  `db:"height"`
}

// NewValidatorData allows to build a new ValidatorData
func NewValidatorData(
	consAddress, valAddress, consPubKey, selfDelegateAddress, maxRate, maxChangeRate string, height int64,
) ValidatorData {
	return ValidatorData{
		ConsAddress:         consAddress,
		ValAddress:          valAddress,
		ConsPubKey:          consPubKey,
		SelfDelegateAddress: selfDelegateAddress,
		MaxRate:             maxRate,
		MaxChangeRate:       maxChangeRate,
		Height:              height,
	}
}

// GetConsAddr implements types.Validator
func (v ValidatorData) GetConsAddr() string {
	return v.ConsAddress
}

// GetConsPubKey implements types.Validator
func (v ValidatorData) GetConsPubKey() string {
	return v.ConsPubKey
}

// GetOperator implements types.Validator
func (v ValidatorData) GetOperator() string {
	return v.ValAddress
}

// GetSelfDelegateAddress implements types.Validator
func (v ValidatorData) GetSelfDelegateAddress() string {
	return v.SelfDelegateAddress
}

// GetMaxChangeRate implements types.Validator
func (v ValidatorData) GetMaxChangeRate() *sdk.Dec {
	n, err := strconv.ParseInt(v.MaxChangeRate, 10, 64)
	if err != nil {
		panic(err)
	}
	result := sdk.NewDec(n)
	return &result
}

// GetMaxRate implements types.Validator
func (v ValidatorData) GetMaxRate() *sdk.Dec {
	n, err := strconv.ParseInt(v.MaxRate, 10, 64)
	if err != nil {
		panic(err)
	}
	result := sdk.NewDec(n)
	return &result
}

// GetHeight implements types.Validator
func (v ValidatorData) GetHeight() int64 {
	return v.Height
}

// ________________________________________________

// ValidatorRow represents a single row of the validator table
type ValidatorRow struct {
	ConsAddress string `db:"consensus_address"`
	ConsPubKey  string `db:"consensus_pubkey"`
}

// NewValidatorRow returns a new ValidatorRow
func NewValidatorRow(consAddress, consPubKey string) ValidatorRow {
	return ValidatorRow{
		ConsAddress: consAddress,
		ConsPubKey:  consPubKey,
	}
}

// Equal tells whether v and w contain the same data
func (v ValidatorRow) Equal(w ValidatorRow) bool {
	return v.ConsAddress == w.ConsAddress &&
		v.ConsPubKey == w.ConsPubKey
}

// ________________________________________________

// ValidatorInfoRow represents a single row of the validator_info table
type ValidatorInfoRow struct {
	ConsAddress         string `db:"consensus_address"`
	ValAddress          string `db:"operator_address"`
	SelfDelegateAddress string `db:"self_delegate_address"`
	MaxChangeRate       string `db:"max_change_rate"`
	MaxRate             string `db:"max_rate"`
	Height              int64  `db:"height"`
}

// NewValidatorInfoRow allows to build a new ValidatorInfoRow
func NewValidatorInfoRow(
	consAddress, valAddress, selfDelegateAddress, maxRate, maxChangeRate string, height int64,
) ValidatorInfoRow {
	return ValidatorInfoRow{
		ConsAddress:         consAddress,
		ValAddress:          valAddress,
		SelfDelegateAddress: selfDelegateAddress,
		MaxChangeRate:       maxChangeRate,
		MaxRate:             maxRate,
		Height:              height,
	}
}

// Equal tells whether v and w represent the same rows
func (v ValidatorInfoRow) Equal(w ValidatorInfoRow) bool {
	return v.ConsAddress == w.ConsAddress &&
		v.ValAddress == w.ValAddress &&
		v.SelfDelegateAddress == w.SelfDelegateAddress &&
		v.MaxRate == w.MaxRate &&
		v.MaxChangeRate == w.MaxChangeRate &&
		v.Height == w.Height
}

// --------------------------------------------------------------------------------------------------------------------

// ValidatorDescriptionRow represent a row in validator_description
type ValidatorDescriptionRow struct {
	ValAddress      string         `db:"validator_address"`
	Moniker         sql.NullString `db:"moniker"`
	Identity        sql.NullString `db:"identity"`
	AvatarURL       sql.NullString `db:"avatar_url"`
	Website         sql.NullString `db:"website"`
	SecurityContact sql.NullString `db:"security_contact"`
	Details         sql.NullString `db:"details"`
	Height          int64          `db:"height"`
}

// NewValidatorDescriptionRow return a row representing data structure in validator_description
func NewValidatorDescriptionRow(
	valAddress, moniker, identity, avatarURL, website, securityContact, details string, height int64,
) ValidatorDescriptionRow {
	return ValidatorDescriptionRow{
		ValAddress:      valAddress,
		Moniker:         ToNullString(moniker),
		Identity:        ToNullString(identity),
		AvatarURL:       ToNullString(avatarURL),
		Website:         ToNullString(website),
		SecurityContact: ToNullString(securityContact),
		Details:         ToNullString(details),
		Height:          height,
	}
}

// Equals return true if two ValidatorDescriptionRow are equal
func (w ValidatorDescriptionRow) Equals(v ValidatorDescriptionRow) bool {
	return v.ValAddress == w.ValAddress &&
		v.Moniker == w.Moniker &&
		v.Identity == w.Identity &&
		v.Website == w.Website &&
		v.SecurityContact == w.SecurityContact &&
		v.Details == w.Details &&
		v.Height == w.Height
}

// ________________________________________________

// ValidatorCommissionRow represents a single row of the validator_commission database table
type ValidatorCommissionRow struct {
	OperatorAddress   string         `db:"validator_address"`
	Commission        sql.NullString `db:"commission"`
	MinSelfDelegation sql.NullString `db:"min_self_delegation"`
	Height            int64          `db:"height"`
}

// NewValidatorCommissionRow allows to easily build a new ValidatorCommissionRow instance
func NewValidatorCommissionRow(
	operatorAddress string, commission string, minSelfDelegation string, height int64,
) ValidatorCommissionRow {
	return ValidatorCommissionRow{
		OperatorAddress:   operatorAddress,
		Commission:        ToNullString(commission),
		MinSelfDelegation: ToNullString(minSelfDelegation),
		Height:            height,
	}
}

// Equal tells whether v and w represent the same rows
func (v ValidatorCommissionRow) Equal(w ValidatorCommissionRow) bool {
	return v.OperatorAddress == w.OperatorAddress &&
		v.Commission == w.Commission &&
		v.MinSelfDelegation == w.MinSelfDelegation &&
		v.Height == w.Height
}

// ________________________________________________

// ValidatorVotingPowerRow represents a single row of the validator_voting_power database table
type ValidatorVotingPowerRow struct {
	ValidatorAddress string `db:"validator_address"`
	VotingPower      int64  `db:"voting_power"`
	Height           int64  `db:"height"`
}

// NewValidatorVotingPowerRow allows to easily build a new ValidatorVotingPowerRow instance
func NewValidatorVotingPowerRow(address string, votingPower int64, height int64) ValidatorVotingPowerRow {
	return ValidatorVotingPowerRow{
		ValidatorAddress: address,
		VotingPower:      votingPower,
		Height:           height,
	}
}

// Equal tells whether v and w represent the same rows
func (v ValidatorVotingPowerRow) Equal(w ValidatorVotingPowerRow) bool {
	return v.ValidatorAddress == w.ValidatorAddress &&
		v.VotingPower == w.VotingPower &&
		v.Height == w.Height
}

// ________________________________________________

// ValidatorStatusRow represents a single row of the validator_status table
type ValidatorStatusRow struct {
	Status      int    `db:"status"`
	Jailed      bool   `db:"jailed"`
	ConsAddress string `db:"validator_address"`
	Height      int64  `db:"height"`
}

// NewValidatorStatusRow builds a new ValidatorStatusRow
func NewValidatorStatusRow(status int, jailed bool, consAddess string, height int64) ValidatorStatusRow {
	return ValidatorStatusRow{
		Status:      status,
		Jailed:      jailed,
		ConsAddress: consAddess,
		Height:      height,
	}
}

// Equal tells whether v and w contain the same data
func (v ValidatorStatusRow) Equal(w ValidatorStatusRow) bool {
	return v.Status == w.Status &&
		v.Jailed == w.Jailed &&
		v.ConsAddress == w.ConsAddress &&
		v.Height == w.Height
}

//--------------------------------------------------------

// DoubleSignVoteRow represents a single row of the double_sign_vote table
type DoubleSignVoteRow struct {
	ID               int64  `db:"id"`
	VoteType         int    `db:"type"`
	Height           int64  `db:"height"`
	Round            int    `db:"round"`
	BlockID          string `db:"block_id"`
	ValidatorAddress string `db:"validator_address"`
	ValidatorIndex   int    `db:"validator_index"`
	Signature        string `db:"signature"`
}

// NewDoubleSignVoteRow allows to build a new NewDoubleSignVoteRow
func NewDoubleSignVoteRow(
	id int64,
	voteType int,
	height int64,
	round int,
	blockID string,
	validatorAddress string,
	validatorIndex int,
	signature string,
) DoubleSignVoteRow {
	return DoubleSignVoteRow{
		ID:               id,
		VoteType:         voteType,
		Height:           height,
		Round:            round,
		BlockID:          blockID,
		ValidatorAddress: validatorAddress,
		ValidatorIndex:   validatorIndex,
		Signature:        signature,
	}
}

// Equal tells whether v and w represent the same rows
func (v DoubleSignVoteRow) Equal(w DoubleSignVoteRow) bool {
	return v.ID == w.ID &&
		v.VoteType == w.VoteType &&
		v.Height == w.Height &&
		v.Round == w.Round &&
		v.BlockID == w.BlockID &&
		v.ValidatorAddress == w.ValidatorAddress &&
		v.ValidatorIndex == w.ValidatorIndex &&
		v.Signature == w.Signature
}

//--------------------------------------------------------

// DoubleSignEvidenceRow represents a single row of the double_sign_evidence table
type DoubleSignEvidenceRow struct {
	Height  int64 `db:"height"`
	VoteAID int64 `db:"vote_a_id"`
	VoteBID int64 `db:"vote_b_id"`
}

// NewDoubleSignEvidenceRow allows to build a new NewDoubleSignEvidenceRow
func NewDoubleSignEvidenceRow(height int64, voteAID int64, voteBID int64) DoubleSignEvidenceRow {
	return DoubleSignEvidenceRow{
		Height:  height,
		VoteAID: voteAID,
		VoteBID: voteBID,
	}
}

// Equal tells whether v and w represent the same rows
func (v DoubleSignEvidenceRow) Equal(w DoubleSignEvidenceRow) bool {
	return v.VoteAID == w.VoteAID &&
		v.VoteBID == w.VoteBID &&
		v.Height == w.Height
}
