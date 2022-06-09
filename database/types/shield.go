package types

type ShieldPool struct {
	PoolID             int64    `db:"pool_id"`
	FromAddress        string   `db:"from_address"`
	Shield             *DbCoins `db:"shield"`
	NativeServiceFees  *DbCoins `db:"native_service_fees"`
	ForeignServiceFees *DbCoins `db:"foreign_service_fees"`
	Sponsor            string   `db:"sponsor"`
	SponsorAddr        string   `db:"sponser_address"`
	Description        string   `db:"description"`
	ShieldLimit        string   `db:"shield_limit"`
	Pause              bool     `db:"pause"`
	Height             int64    `db:"height"`
}
