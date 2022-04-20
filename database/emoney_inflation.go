package database

import (
	"encoding/json"

	"github.com/forbole/bdjuno/v3/types"
)

// SaveEMoneyInflation allows to store the eMoney inflation (scheduler = per day)
func (db *Db) SaveEMoneyInflation(eMoneyInflation types.EMoneyInflation) error {

	inflationBz, err := json.Marshal(&eMoneyInflation.InflationAssets)
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
	_, err = db.Sql.Exec(stmt, string(inflationBz), eMoneyInflation.LastAppliedTime, eMoneyInflation.LastAppliedHeight, eMoneyInflation.Height)
	if err != nil {
		return err
	}
	return nil
}
