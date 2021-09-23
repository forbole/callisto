package database

import (
	"encoding/json"

	"github.com/e-money/em-ledger/x/inflation/types"
)

// SaveEmoneyInflation allows to store the emoney inflation (scheduler = per day)
func (db *Db) SaveEmoneyInflation(state types.InflationState, height int64) error {

	inflationBz, err := json.Marshal(&state.InflationAssets)
	if err != nil {
		return err
	}

	stmt := `
INSERT INTO emoney_inflation (inflation, last_applied_time, last_applied_height, height)
VALUES ($1, $2, $3)
ON CONFLICT (one_row_id) DO UPDATE
    SET inflation = excluded.inflation,
		last_applied_time = excluded.last_applied_time,
		last_applied_height = excluded.last_applied_height,
		height = excluded.height
WHERE emoney_inflation.height <= excluded.height`
	_, err = db.Sql.Exec(stmt, string(inflationBz), state.LastAppliedTime, state.LastAppliedHeight.Int64(), height)
	if err != nil {
		return err
	}
	return nil
}
