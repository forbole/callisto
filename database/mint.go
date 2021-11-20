package database

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	dbtypes "github.com/forbole/bdjuno/v2/database/types"
	"github.com/forbole/bdjuno/v2/types"
)

// SaveInflation allows to store the inflation for the given block height as well as timestamp
func (db *Db) SaveInflation(inflation sdk.Dec, height int64) error {
	stmt := `
INSERT INTO inflation (value, height) 
VALUES ($1, $2) 
ON CONFLICT (one_row_id) DO UPDATE 
    SET value = excluded.value, 
        height = excluded.height 
WHERE inflation.height <= excluded.height`

	_, err := db.Sql.Exec(stmt, inflation.String(), height)
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

// GetMintParams returns the most recent mint parameters
func (db *Db) GetMintParams() (*types.MintParams, error) {
	var rows []dbtypes.MintParamsRow
	err := db.Sqlx.Select(&rows, `SELECT * FROM mint_params`)
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, nil
	}

	row := rows[0]

	var mintParams minttypes.Params
	err = json.Unmarshal([]byte(row.Params), &mintParams)
	if err != nil {
		return nil, err
	}

	return types.NewMintParams(mintParams, row.Height), nil
}
