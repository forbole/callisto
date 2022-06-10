package types

type ShieldPoolRow struct {
	PoolID             int64    `db:"pool_id"`
	Shield             string   `db:"shield"`
	NativeServiceFees  *DbCoins `db:"native_service_fees"`
	ForeignServiceFees *DbCoins `db:"foreign_service_fees"`
	Sponsor            string   `db:"sponsor"`
	SponsorAddr        string   `db:"sponsor_address"`
	Description        string   `db:"description"`
	ShieldLimit        string   `db:"shield_limit"`
	Pause              bool     `db:"pause"`
	Height             int64    `db:"height"`
}

// Equal tells whether v and w represent the same rows
func (v ShieldPoolRow) Equal(w ShieldPoolRow) bool {
	return v.PoolID == w.PoolID &&
		v.Shield == w.Shield &&
		v.NativeServiceFees.Equal(w.NativeServiceFees) &&
		v.ForeignServiceFees.Equal(w.ForeignServiceFees) &&
		v.Sponsor == w.Sponsor &&
		v.SponsorAddr == w.SponsorAddr &&
		v.Description == w.Description &&
		v.ShieldLimit == w.ShieldLimit &&
		v.Pause == w.Pause &&
		v.Height == w.Height
}

// NewShieldPoolRow allows to build a new ShieldPoolRow
func NewShieldPoolRow(
	poolID int64,
	shield string,
	nativeServiceFees DbCoins,
	foreignServiceFees DbCoins,
	sponsor string,
	sponsorAddress string,
	description string,
	shieldLimit string,
	pause bool,
	height int64,
) ShieldPoolRow {
	return ShieldPoolRow{
		PoolID:             poolID,
		Shield:             shield,
		NativeServiceFees:  &nativeServiceFees,
		ForeignServiceFees: &foreignServiceFees,
		Sponsor:            sponsor,
		SponsorAddr:        sponsorAddress,
		Description:        description,
		ShieldLimit:        shieldLimit,
		Pause:              pause,
		Height:             height,
	}
}
