package database

import (
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
	stmt := `
INSERT INTO mint_params (mint_denom, inflation_rate_change, inflation_min, inflation_max, goal_bonded, blocks_per_year, height) 
VALUES ($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT (one_row_id) DO UPDATE 
    SET mint_denom = excluded.mint_denom,
    	inflation_rate_change = excluded.inflation_rate_change, 
    	inflation_min = excluded.inflation_min, 
    	inflation_max = excluded.inflation_max,
    	goal_bonded = excluded.goal_bonded,
    	blocks_per_year = excluded.blocks_per_year,
        height = excluded.height
WHERE mint_params.height <= excluded.height`
	_, err := db.Sql.Exec(stmt, params.MintDenom,
		params.InflationRateChange.String(), params.InflationMin.String(), params.InflationMax.String(),
		params.GoalBonded.String(), params.BlocksPerYear, params.Height)
	return err
}
