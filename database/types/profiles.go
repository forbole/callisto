package types

// ProfilesParamsRow represents a single row inside the profiles_params table
type ProfilesParamsRow struct {
	OneRowID bool   `db:"one_row_id"`
	Params   string `db:"params"`
	Height   int64  `db:"height"`
}
