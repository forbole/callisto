package types

type ShieldPool struct {
	PoolID         int64    `db:"pool_id"`
	FromAddress    string   `db:"from_address"`
	Shield         *DbCoins `db:"shield"`
	NativeDeposit  *DbCoins `db:"native_deposit"`
	ForeignDeposit *DbCoins `db:"foreign_deposit"`
	Sponsor        string   `db:"sponsor"`
	SponsorAddr    string   `db:"sponser_address"`
	Description    string   `db:"description"`
	ShieldLimit    string   `db:"shield_limit"`
	Pause          bool     `db:"pause"`
	Height         int64    `db:"height"`
}
