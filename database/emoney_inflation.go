package database

import (
	"encoding/json"
	"fmt"

	"github.com/e-money/em-ledger/x/inflation/types"
)

// SaveInflation allows to store the inflation for the given block height as well as timestamp
func (db *Db) SaveEmoneyInflation(state types.InflationState) error {

	inflationBz, err := json.Marshal(&state.InflationAssets)
	if err != nil {
		return err
	}

	fmt.Println("string(inflationBz)")
	fmt.Println(string(inflationBz))

	stmt := `
INSERT INTO emoney_inflation (inflation, last_applied_time, last_applied_height)
VALUES ($1, $2, $3)
ON CONFLICT (one_row_id) DO UPDATE
    SET inflation = excluded.inflation,
		last_applied_time = excluded.last_applied_time,
		last_applied_height = excluded.last_applied_height
WHERE emoney_inflation.last_applied_height <= excluded.last_applied_height`
	_, err = db.Sql.Exec(stmt, string(inflationBz), state.LastAppliedTime, state.LastAppliedHeight.Int64())
	if err != nil {
		return err
	}
	return nil
}
