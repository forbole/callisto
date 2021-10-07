package types

// StakingParamsRow represents a single row inside the staking_params table
type StakingParamsRow struct {
	OneRowID bool   `db:"one_row_id"`
	Params   string `db:"params"`
	Height   int64  `db:"height"`
}
