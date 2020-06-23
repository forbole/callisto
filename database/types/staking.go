package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
)

// StakingPoolRow represents a single row inside the staking_pool table
type StakingPoolRow struct {
	BondedTokens    int64     `db:"bonded_tokens"`
	NotBondedTokens int64     `db:"not_bonded_tokens"`
	Height          int64     `db:"height"`
	Timestamp       time.Time `db:"timestamp"`
}

// NewStakingPoolRow allows to easily create a new StakingPoolRow
func NewStakingPoolRow(bondedTokens, notBondedTokens int64, height int64, timestamp time.Time) StakingPoolRow {
	return StakingPoolRow{
		BondedTokens:    bondedTokens,
		NotBondedTokens: notBondedTokens,
		Height:          height,
		Timestamp:       timestamp,
	}
}

// Equal allows to tells whether r and as represent the same rows
func (r StakingPoolRow) Equal(s StakingPoolRow) bool {
	return r.BondedTokens == s.BondedTokens &&
		r.NotBondedTokens == s.NotBondedTokens &&
		r.Height == s.Height &&
		r.Timestamp.Equal(s.Timestamp)
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
	ConsAddress string `db:"consensus_address"`
	ValAddress  string `db:"operator_address"`
}

// NewValidatorInfoRow allows to build a new ValidatorInfoRow
func NewValidatorInfoRow(consAddress, valAddress string) ValidatorInfoRow {
	return ValidatorInfoRow{
		ConsAddress: consAddress,
		ValAddress:  valAddress,
	}
}

// Equal tells whether v and w represent the same rows
func (v ValidatorInfoRow) Equal(w ValidatorInfoRow) bool {
	return v.ConsAddress == w.ConsAddress &&
		v.ValAddress == w.ValAddress
}

// ________________________________________________

// ValidatorData contains all the data of a single validator.
// It implements types.Validator interface
type ValidatorData struct {
	ConsAddress string `db:"consensus_address"`
	ValAddress  string `db:"operator_address"`
	ConsPubKey  string `db:"consensus_pubkey"`
}

// NewValidatorData allows to build a new ValidatorData
func NewValidatorData(consAddress, valAddress, consPubKey string) ValidatorData {
	return ValidatorData{
		ConsAddress: consAddress,
		ValAddress:  valAddress,
		ConsPubKey:  consPubKey,
	}
}

func (v ValidatorData) GetConsAddr() sdk.ConsAddress {
	addr, err := sdk.ConsAddressFromBech32(v.ConsAddress)
	if err != nil {
		panic(err)
	}

	return addr
}

func (v ValidatorData) GetConsPubKey() crypto.PubKey {
	return sdk.MustGetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, v.ConsPubKey)
}

func (v ValidatorData) GetOperator() sdk.ValAddress {
	addr, err := sdk.ValAddressFromBech32(v.ValAddress)
	if err != nil {
		panic(err)
	}

	return addr
}

// ________________________________________________

// ValidatorUptimeRow represents a single row of the validator_uptime table
type ValidatorUptimeRow struct {
	ConsAddr           string `db:"validator_address"`
	Height             int64  `db:"height"`
	SignedBlockWindow  int64  `db:"signed_blocks_window"`
	MissedBlockCounter int64  `db:"missed_blocks_counter"`
}

// NewValidatorUptimeRow allows to build a new ValidatorUptimeRow
func NewValidatorUptimeRow(consAddr string, signedBlocWindow, missedBlocksCounter, height int64) ValidatorUptimeRow {
	return ValidatorUptimeRow{
		ConsAddr:           consAddr,
		SignedBlockWindow:  signedBlocWindow,
		MissedBlockCounter: missedBlocksCounter,
		Height:             height,
	}
}

// Equal tells whether v and w contain the same data
func (v ValidatorUptimeRow) Equal(w ValidatorUptimeRow) bool {
	return v.ConsAddr == w.ConsAddr &&
		v.Height == w.Height &&
		v.SignedBlockWindow == w.SignedBlockWindow &&
		v.MissedBlockCounter == w.MissedBlockCounter
}

// ________________________________________________

// ValidatorDelegationRow represents a single validator_delegation table row
type ValidatorDelegationRow struct {
	ConsensusAddress string    `db:"consensus_address"`
	DelegatorAddress string    `db:"delegator_address"`
	Amount           DbCoin    `db:"amount"`
	Height           int64     `db:"height"`
	Timestamp        time.Time `db:"timestamp"`
}

// NewValidatorDelegationRow allows to build a new ValidatorDelegationRow
func NewValidatorDelegationRow(
	consAddr, delegator string, amount DbCoin,
	height int64, timestamp time.Time,
) ValidatorDelegationRow {
	return ValidatorDelegationRow{
		ConsensusAddress: consAddr,
		DelegatorAddress: delegator,
		Amount:           amount,
		Height:           height,
		Timestamp:        timestamp,
	}
}

// Equals tells whether v and w represent the same row
func (v ValidatorDelegationRow) Equal(w ValidatorDelegationRow) bool {
	return v.ConsensusAddress == w.ConsensusAddress &&
		v.DelegatorAddress == w.DelegatorAddress &&
		v.Amount.Equal(w.Amount) &&
		v.Height == w.Height &&
		v.Timestamp.Equal(w.Timestamp)
}

// ________________________________________________

// ValidatorUnbondingDelegationRow represents a single row inside the
// validator_unbonding_delegation table
type ValidatorUnbondingDelegationRow struct {
	ConsensusAddress    string    `db:"consensus_address"`
	DelegatorAddress    string    `db:"delegator_address"`
	Amount              DbCoin    `db:"amount"`
	CompletionTimestamp time.Time `db:"completion_timestamp"`
	Height              int64     `db:"height"`
	Timestamp           time.Time `db:"timestamp"`
}

// NewValidatorUnbondingDelegationRow allows to build a new
// ValidatorUnbondingDelegationRow instance
func NewValidatorUnbondingDelegationRow(
	consAddr, delegator string, amount DbCoin, completionTimestamp time.Time,
	height int64, timestamp time.Time,
) ValidatorUnbondingDelegationRow {
	return ValidatorUnbondingDelegationRow{
		ConsensusAddress:    consAddr,
		DelegatorAddress:    delegator,
		Amount:              amount,
		CompletionTimestamp: completionTimestamp,
		Height:              height,
		Timestamp:           timestamp,
	}
}

// Equal tells whether v and w represent the same rows
func (v ValidatorUnbondingDelegationRow) Equal(w ValidatorUnbondingDelegationRow) bool {
	return v.ConsensusAddress == w.ConsensusAddress &&
		v.DelegatorAddress == w.DelegatorAddress &&
		v.Amount.Equal(w.Amount) &&
		v.CompletionTimestamp.Equal(w.CompletionTimestamp) &&
		v.Height == w.Height &&
		v.Timestamp.Equal(w.Timestamp)
}

// ________________________________________________

// ValidatorReDelegationRow represents a single row of the
// validator_redelegation database table
type ValidatorReDelegationRow struct {
	DelegatorAddress    string    `db:"delegator_address"`
	SrcValidatorAddress string    `db:"src_validator_address"`
	DstValidatorAddress string    `db:"dst_validator_address"`
	Amount              DbCoin    `db:"amount"`
	Height              int64     `db:"height"`
	CompletionTime      time.Time `db:"completion_time"`
}

// NewValidatorReDelegationRow allows to easily build a new
// ValidatorReDelegationRow instance
func NewValidatorReDelegationRow(
	delegator, srcConsAddr, dstConsAddr string, amount DbCoin,
	height int64, completionTime time.Time,
) ValidatorReDelegationRow {
	return ValidatorReDelegationRow{
		DelegatorAddress:    delegator,
		SrcValidatorAddress: srcConsAddr,
		DstValidatorAddress: dstConsAddr,
		Amount:              amount,
		Height:              height,
		CompletionTime:      completionTime,
	}
}

// Equals tells whether v and w represent the same database rows
func (v ValidatorReDelegationRow) Equal(w ValidatorReDelegationRow) bool {
	return v.DelegatorAddress == w.DelegatorAddress &&
		v.SrcValidatorAddress == w.SrcValidatorAddress &&
		v.DstValidatorAddress == w.DstValidatorAddress &&
		v.Amount.Equal(w.Amount) &&
		v.Height == w.Height &&
		v.CompletionTime.Equal(w.CompletionTime)
}
