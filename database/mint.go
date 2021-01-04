package database

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SaveInflation allows to store the inflation for the given block height as well as timestamp
func (db *BigDipperDb) SaveInflation(inflation sdk.Dec, height int64, timestamp time.Time) error {
	stmt := `INSERT INTO inflation_history (value, height, timestamp) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	_, err := db.Sql.Exec(stmt, inflation.String(), height, timestamp)
	return err
}
