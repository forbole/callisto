package types

import (
	"time"
)

// ________________________________________________

// DelegationRow represents a single delegation table row
type DelegationRow struct {
	ID               string `db:"id"`
	ValidatorAddress string `db:"validator_address"`
	DelegatorAddress string `db:"delegator_address"`
	Amount           DbCoin `db:"amount"`
	Height           int64  `db:"height"`
}

// NewDelegationRow allows to build a new DelegationRow
func NewDelegationRow(delegator, consAddr string, amount DbCoin, height int64) DelegationRow {
	return DelegationRow{
		ValidatorAddress: consAddr,
		DelegatorAddress: delegator,
		Amount:           amount,
		Height:           height,
	}
}

// Equal tells whether v and w represent the same row
func (v DelegationRow) Equal(w DelegationRow) bool {
	return v.ValidatorAddress == w.ValidatorAddress &&
		v.DelegatorAddress == w.DelegatorAddress &&
		v.Amount.Equal(w.Amount) &&
		v.Height == w.Height
}

// ________________________________________________

// UnbondingDelegationRow represents a single row of the unbonding_delegation table
type UnbondingDelegationRow struct {
	ConsensusAddress    string    `db:"validator_address"`
	DelegatorAddress    string    `db:"delegator_address"`
	Amount              DbCoin    `db:"amount"`
	CompletionTimestamp time.Time `db:"completion_timestamp"`
	Height              int64     `db:"height"`
}

// NewUnbondingDelegationRow allows to build a new UnbondingDelegationRow instance
func NewUnbondingDelegationRow(
	delegator, consAddr string, amount DbCoin, completionTimestamp time.Time, height int64,
) UnbondingDelegationRow {
	return UnbondingDelegationRow{
		ConsensusAddress:    consAddr,
		DelegatorAddress:    delegator,
		Amount:              amount,
		CompletionTimestamp: completionTimestamp,
		Height:              height,
	}
}

// Equal tells whether v and w represent the same rows
func (v UnbondingDelegationRow) Equal(w UnbondingDelegationRow) bool {
	return v.ConsensusAddress == w.ConsensusAddress &&
		v.DelegatorAddress == w.DelegatorAddress &&
		v.Amount.Equal(w.Amount) &&
		v.CompletionTimestamp.Equal(w.CompletionTimestamp) &&
		v.Height == w.Height
}

// ________________________________________________

// RedelegationRow represents a single row of the redelegation database table
type RedelegationRow struct {
	DelegatorAddress    string    `db:"delegator_address"`
	SrcValidatorAddress string    `db:"src_validator_address"`
	DstValidatorAddress string    `db:"dst_validator_address"`
	Amount              DbCoin    `db:"amount"`
	CompletionTime      time.Time `db:"completion_time"`
	Height              int64     `db:"height"`
}

// NewRedelegationRow allows to easily build a new RedelegationRow instance
func NewRedelegationRow(
	delegator, srcConsAddr, dstConsAddr string, amount DbCoin, completionTime time.Time, height int64,
) RedelegationRow {
	return RedelegationRow{
		DelegatorAddress:    delegator,
		SrcValidatorAddress: srcConsAddr,
		DstValidatorAddress: dstConsAddr,
		Amount:              amount,
		CompletionTime:      completionTime,
		Height:              height,
	}
}

// Equal tells whether v and w represent the same database rows
func (v RedelegationRow) Equal(w RedelegationRow) bool {
	return v.DelegatorAddress == w.DelegatorAddress &&
		v.SrcValidatorAddress == w.SrcValidatorAddress &&
		v.DstValidatorAddress == w.DstValidatorAddress &&
		v.Amount.Equal(w.Amount) &&
		v.CompletionTime.Equal(w.CompletionTime) &&
		v.Height == w.Height
}
