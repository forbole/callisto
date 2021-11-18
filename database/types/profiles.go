package types

// ProfilesParamsRow represents a single row inside the profiles_params table
type ProfilesParamsRow struct {
	OneRowID       bool   `db:"one_row_id"`
	NickNameParams string `db:"nickname_params"`
	DTagParams     string `db:"d_tag_params"`
	BioParams      string `db:"bio_params"`
	OracleParams   string `db:"oracle_params"`
	Height         int64  `db:"height"`
}
