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

// NewProfilesParamsRow allows to easily create a new ProfilesParamsRow
func NewProfilesParamsRow(nickNameParams string, dTagParams string, bioParams string, oracleParams string, height int64) ProfilesParamsRow {
	return ProfilesParamsRow{
		OneRowID:       true,
		NickNameParams: nickNameParams,
		DTagParams:     dTagParams,
		BioParams:      bioParams,
		OracleParams:   oracleParams,
		Height:         height,
	}
}

// Equal allows to tells whether a and b represent the same row
func (a ProfilesParamsRow) Equal(b ProfilesParamsRow) bool {
	return a.NickNameParams == b.NickNameParams &&
		a.DTagParams == b.DTagParams &&
		a.BioParams == b.BioParams &&
		a.OracleParams == b.OracleParams &&
		a.Height == b.Height
}
