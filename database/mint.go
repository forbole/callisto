package database

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/forbole/bdjuno/types"
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
	return err
}

// SaveMintParams allows to store the given params inside the database
func (db *Db) SaveMintParams(params types.MintParams) error {
	paramsBz, err := json.Marshal(&params.Params)
	if err != nil {
		return err
	}

	stmt := `
INSERT INTO mint_params (params, height) 
VALUES ($1, $2)
ON CONFLICT (one_row_id) DO UPDATE 
    SET params = excluded.params,
        height = excluded.height
WHERE mint_params.height <= excluded.height`
	_, err = db.Sql.Exec(stmt, string(paramsBz), params.Height)
	return err
}

// GetEpochIdentifier returns epoch_identifier param stored in db
func (db *Db) GetEpochIdentifier(height int64) (string, error) {
	stmt := `SELECT params FROM mint_params WHERE height = $1`

	var epochIdentifier = ""
	err := db.Sqlx.Select(&epochIdentifier, stmt, height)
	if err != nil {
		return "", err
	}

	return epochIdentifier, nil
}
