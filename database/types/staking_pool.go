package types

// StakingPoolRow represents a single row inside the staking_pool table
type StakingPoolRow struct {
	OneRowID        bool  `db:"one_row_id"`
	BondedTokens    int64 `db:"bonded_tokens"`
	NotBondedTokens int64 `db:"not_bonded_tokens"`
	UnbondingTokens int64 `db:"unbonding_tokens"`
	Height          int64 `db:"height"`
}

// NewStakingPoolRow allows to easily create a new StakingPoolRow
func NewStakingPoolRow(bondedTokens, notBondedTokens, unbondingTokens int64, height int64) StakingPoolRow {
	return StakingPoolRow{
		OneRowID:        true,
		BondedTokens:    bondedTokens,
		NotBondedTokens: notBondedTokens,
		UnbondingTokens: unbondingTokens,
		Height:          height,
	}
}

// Equal allows to tells whether r and as represent the same rows
func (r StakingPoolRow) Equal(s StakingPoolRow) bool {
	return r.BondedTokens == s.BondedTokens &&
		r.NotBondedTokens == s.NotBondedTokens &&
		r.UnbondingTokens == s.UnbondingTokens &&
		r.Height == s.Height
}
