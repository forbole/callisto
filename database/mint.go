package database

import (
	"encoding/json"
	"fmt"

	creminttypes "github.com/crescent-network/crescent/v2/x/mint/types"
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
	var rows []creminttypes.Params
	err := db.Sqlx.Select(&rows, `SELECT * FROM mint_params`)
	if err != nil {
		return creminttypes.Params{}, fmt.Errorf("error while getting supply: %s", err)
	}

	return rows[0], nil
}
