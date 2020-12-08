package types

import "time"

// ________________________________________________

// DelegationRow represents a single delegation table row
type DelegationRow struct {
	ValidatorAddress string  `db:"validator_address"`
	DelegatorAddress string  `db:"delegator_address"`
	Amount           DbCoin  `db:"amount"`
	Shares           float64 `db:"shares"`
}

// NewDelegationRow allows to build a new DelegationRow
func NewDelegationRow(consAddr, delegator string, amount DbCoin, shares float64) DelegationRow {
	return DelegationRow{
		ValidatorAddress: consAddr,
		DelegatorAddress: delegator,
		Amount:           amount,
		Shares:           shares,
	}
}

// Equals tells whether v and w represent the same row
func (v DelegationRow) Equal(w DelegationRow) bool {
	return v.ValidatorAddress == w.ValidatorAddress &&
		v.DelegatorAddress == w.DelegatorAddress &&
		v.Amount.Equal(w.Amount) &&
		v.Shares == w.Shares
}

// DelegationHistoryRow represents a single row of the delegations_history table
type DelegationHistoryRow struct {
	ValidatorAddress string    `db:"validator_address"`
	DelegatorAddress string    `db:"delegator_address"`
	Amount           DbCoin    `db:"amount"`
	Shares           float64   `db:"shares"`
	Height           int64     `db:"height"`
	Timestamp        time.Time `db:"timestamp"`
}

// NewDelegationHistoryRow allows to build a new DelegationHistoryRow
func NewDelegationHistoryRow(
	consAddr, delegator string, amount DbCoin, shares float64, height int64, timestamp time.Time,
) DelegationHistoryRow {
	return DelegationHistoryRow{
		ValidatorAddress: consAddr,
		DelegatorAddress: delegator,
		Amount:           amount,
		Shares:           shares,
		Height:           height,
		Timestamp:        timestamp,
	}
}

// Equals tells whether v and w represent the same row
func (v DelegationHistoryRow) Equal(w DelegationHistoryRow) bool {
	return v.ValidatorAddress == w.ValidatorAddress &&
		v.DelegatorAddress == w.DelegatorAddress &&
		v.Amount.Equal(w.Amount) &&
		v.Shares == w.Shares &&
		v.Height == w.Height &&
		v.Timestamp.Equal(w.Timestamp)
}

// ________________________________________________

// UnbondingDelegationRow represents a single row of the unbonding_delegation table
type UnbondingDelegationRow struct {
	ConsensusAddress    string    `db:"validator_address"`
	DelegatorAddress    string    `db:"delegator_address"`
	Amount              DbCoin    `db:"amount"`
	CompletionTimestamp time.Time `db:"completion_timestamp"`
}

// NewUnbondingDelegationRow allows to build a new UnbondingDelegationRow instance
func NewUnbondingDelegationRow(consAddr, delegator string, amount DbCoin, completionTimestamp time.Time) UnbondingDelegationRow {
	return UnbondingDelegationRow{
		ConsensusAddress:    consAddr,
		DelegatorAddress:    delegator,
		Amount:              amount,
		CompletionTimestamp: completionTimestamp,
	}
}

// Equal tells whether v and w represent the same rows
func (v UnbondingDelegationRow) Equal(w UnbondingDelegationRow) bool {
	return v.ConsensusAddress == w.ConsensusAddress &&
		v.DelegatorAddress == w.DelegatorAddress &&
		v.Amount.Equal(w.Amount) &&
		v.CompletionTimestamp.Equal(w.CompletionTimestamp)
}

// UnbondingDelegationHistoryRow represents a single row of the unbonding_delegation_history table
type UnbondingDelegationHistoryRow struct {
	ConsensusAddress    string    `db:"validator_address"`
	DelegatorAddress    string    `db:"delegator_address"`
	Amount              DbCoin    `db:"amount"`
	CompletionTimestamp time.Time `db:"completion_timestamp"`
	Height              int64     `db:"height"`
	Timestamp           time.Time `db:"timestamp"`
}

// NewUnbondingDelegationHistoryRow allows to build a new UnbondingDelegationHistoryRow instance
func NewUnbondingDelegationHistoryRow(
	consAddr, delegator string, amount DbCoin, completionTimestamp time.Time, height int64, timestamp time.Time,
) UnbondingDelegationHistoryRow {
	return UnbondingDelegationHistoryRow{
		ConsensusAddress:    consAddr,
		DelegatorAddress:    delegator,
		Amount:              amount,
		CompletionTimestamp: completionTimestamp,
		Height:              height,
		Timestamp:           timestamp,
	}
}

// Equal tells whether v and w represent the same rows
func (v UnbondingDelegationHistoryRow) Equal(w UnbondingDelegationHistoryRow) bool {
	return v.ConsensusAddress == w.ConsensusAddress &&
		v.DelegatorAddress == w.DelegatorAddress &&
		v.Amount.Equal(w.Amount) &&
		v.CompletionTimestamp.Equal(w.CompletionTimestamp) &&
		v.Height == w.Height &&
		v.Timestamp.Equal(w.Timestamp)
}

// ________________________________________________

// ReDelegationRow represents a single row of the redelegation database table
type ReDelegationRow struct {
	DelegatorAddress    string    `db:"delegator_address"`
	SrcValidatorAddress string    `db:"src_validator_address"`
	DstValidatorAddress string    `db:"dst_validator_address"`
	Amount              DbCoin    `db:"amount"`
	CompletionTime      time.Time `db:"completion_time"`
}

// NewReDelegationRow allows to easily build a new ReDelegationRow instance
func NewReDelegationRow(delegator, srcConsAddr, dstConsAddr string, amount DbCoin, completionTime time.Time) ReDelegationRow {
	return ReDelegationRow{
		DelegatorAddress:    delegator,
		SrcValidatorAddress: srcConsAddr,
		DstValidatorAddress: dstConsAddr,
		Amount:              amount,
		CompletionTime:      completionTime,
	}
}

// Equal tells whether v and w represent the same database rows
func (v ReDelegationRow) Equal(w ReDelegationRow) bool {
	return v.DelegatorAddress == w.DelegatorAddress &&
		v.SrcValidatorAddress == w.SrcValidatorAddress &&
		v.DstValidatorAddress == w.DstValidatorAddress &&
		v.Amount.Equal(w.Amount) &&
		v.CompletionTime.Equal(w.CompletionTime)
}

// ReDelegationHistoryRow represents a single row of the redelegation_history database table
type ReDelegationHistoryRow struct {
	DelegatorAddress    string    `db:"delegator_address"`
	SrcValidatorAddress string    `db:"src_validator_address"`
	DstValidatorAddress string    `db:"dst_validator_address"`
	Amount              DbCoin    `db:"amount"`
	CompletionTime      time.Time `db:"completion_time"`
	Height              int64     `db:"height"`
	Timestamp           time.Time `db:"timestamp"`
}

// NewReDelegationHistoryRow allows to easily build a new ReDelegationHistoryRow instance
func NewReDelegationHistoryRow(
	delegator, srcConsAddr, dstConsAddr string, amount DbCoin, completionTime time.Time, height int64, timestamp time.Time,
) ReDelegationHistoryRow {
	return ReDelegationHistoryRow{
		DelegatorAddress:    delegator,
		SrcValidatorAddress: srcConsAddr,
		DstValidatorAddress: dstConsAddr,
		Amount:              amount,
		CompletionTime:      completionTime,
		Height:              height,
		Timestamp:           timestamp,
	}
}

// Equal tells whether v and w represent the same database rows
func (v ReDelegationHistoryRow) Equal(w ReDelegationHistoryRow) bool {
	return v.DelegatorAddress == w.DelegatorAddress &&
		v.SrcValidatorAddress == w.SrcValidatorAddress &&
		v.DstValidatorAddress == w.DstValidatorAddress &&
		v.Amount.Equal(w.Amount) &&
		v.CompletionTime.Equal(w.CompletionTime) &&
		v.Height == w.Height &&
		v.Timestamp.Equal(w.Timestamp)
}
