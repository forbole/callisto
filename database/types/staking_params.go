package types

// StakingParamsRow represents a single row inside the staking_params table
type StakingParamsRow struct {
	OneRowID bool   `db:"one_row_id"`
	BondName string `db:"bond_denom"`
}

// NewStakingParamsRow allows to build a new StakingParamsRow object
func NewStakingParamsRow(bondName string) StakingParamsRow {
	return StakingParamsRow{
		OneRowID: true,
		BondName: bondName,
	}
}

// Equal tells whether r and s contain the same data
func (r StakingParamsRow) Equal(s StakingParamsRow) bool {
	return r.BondName == s.BondName
}
