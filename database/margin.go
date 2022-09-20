package database

import (
	"encoding/json"
	"fmt"

	"github.com/forbole/bdjuno/v3/types"
)

// SaveMarginParams allows to store the given params inside the database
func (db *Db) SaveMarginParams(params *types.MarginParams) error {
	paramsBz, err := json.Marshal(&params.Params)
	if err != nil {
		return fmt.Errorf("error while marshaling margin params: %s", err)
	}

	stmt := `
INSERT INTO margin_params (params, height) 
VALUES ($1, $2)
ON CONFLICT (one_row_id) DO UPDATE 
    SET params = excluded.params,
        height = excluded.height
WHERE margin_params.height <= excluded.height`

	_, err = db.Sql.Exec(stmt, string(paramsBz), params.Height)
	if err != nil {
		return fmt.Errorf("error while storing margin params: %s", err)
	}

	return nil
}

// SaveMarginEvent allows to store the given x/margin events inside the database
func (db *Db) SaveMarginEvent(event types.MarginEvent) error {

	stmt := `
INSERT INTO margin_events (transaction_hash, index, type, value, involved_accounts_addresses, height) 
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT DO NOTHING`

	_, err := db.Sql.Exec(stmt, event.TxHash, event.Index, event.MsgType, event.Value, event.Addressess, event.Height)
	if err != nil {
		return fmt.Errorf("error while storing margin event: %s", err)
	}

	return nil
}
