package types

import "time"

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
