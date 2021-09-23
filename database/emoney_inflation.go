package database

import (
	"encoding/json"

	"github.com/forbole/bdjuno/types"
)

// SaveEmoneyInflation allows to store the emoney inflation (scheduler = per day)
func (db *Db) SaveEmoneyInflation(emoneyInflation types.EmoneyInflation) error {

	inflationBz, err := json.Marshal(&emoneyInflation.InflationAssets)
	if err != nil {
		return err
	}

	stmt := `
INSERT INTO emoney_inflation (inflation, last_applied_time, last_applied_height, height)
VALUES ($1, $2, $3, $4)
ON CONFLICT (one_row_id) DO UPDATE
    SET inflation = excluded.inflation,
		last_applied_time = excluded.last_applied_time,
		last_applied_height = excluded.last_applied_height,
		height = excluded.height
WHERE emoney_inflation.height <= excluded.height`
	_, err = db.Sql.Exec(stmt, string(inflationBz), emoneyInflation.LastAppliedTime, emoneyInflation.LastAppliedHeight, emoneyInflation.Height)
	if err != nil {
		return err
	}
	return nil
}
