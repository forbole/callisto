package database

import (
	"encoding/json"
	"fmt"

	creminttypes "github.com/crescent-network/crescent/v3/x/mint/types"
	dbtypes "github.com/forbole/bdjuno/v3/database/types"
	"github.com/forbole/bdjuno/v3/types"
)

// SaveInflation allows to store the inflation for the given block height as well as timestamp
func (db *Db) SaveInflation(inflation string, height int64) error {
	stmt := `
INSERT INTO inflation (value, height) 
VALUES ($1, $2) 
ON CONFLICT (one_row_id) DO UPDATE 
    SET value = excluded.value, 
        height = excluded.height 
WHERE inflation.height <= excluded.height`

	_, err := db.Sql.Exec(stmt, inflation, height)
	if err != nil {
		return fmt.Errorf("error while storing inflation: %s", err)
	}

	return nil
}

// SaveMintParams allows to store the given params inside the database
func (db *Db) SaveMintParams(params *types.MintParams) error {
	paramsBz, err := json.Marshal(&params.Params)
	if err != nil {
		return fmt.Errorf("error while marshaling mint params: %s", err)
	}

	stmt := `
INSERT INTO mint_params (params, height) 
VALUES ($1, $2)
ON CONFLICT (one_row_id) DO UPDATE 
    SET params = excluded.params,
        height = excluded.height
WHERE mint_params.height <= excluded.height`

	_, err = db.Sql.Exec(stmt, string(paramsBz), params.Height)
	if err != nil {
		return fmt.Errorf("error while storing mint params: %s", err)
	}

	return nil
}

// GetMintParams allows to get the current mint params
func (db *Db) GetMintParams() (creminttypes.Params, error) {
	var rows []dbtypes.MintParamsRow
	err := db.Sqlx.Select(&rows, `SELECT * FROM mint_params`)
	if err != nil {
		return creminttypes.Params{}, fmt.Errorf("error while getting mint params: %s", err)
	}

	if len(rows) == 0 {
		return creminttypes.Params{}, fmt.Errorf("no mint params stored")
	}

	var params creminttypes.Params
	err = json.Unmarshal([]byte(rows[0].Params), &params)
	if err != nil {
		return creminttypes.Params{}, fmt.Errorf("error while unmarshaling mint params: %s", err)
	}

	return params, nil
}
