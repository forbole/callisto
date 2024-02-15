package database

import (
	"encoding/json"
	"fmt"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	dbtypes "github.com/forbole/callisto/v4/database/types"
	"github.com/forbole/callisto/v4/types"
)

// SaveStakingParams allows to store the given params into the database
func (db *Db) SaveStakingParams(params types.StakingParams) error {
	paramsBz, err := json.Marshal(&params.Params)
	if err != nil {
		return fmt.Errorf("error while marshaling staking params: %s", err)
	}

	stmt := `
INSERT INTO staking_params (params, height) 
VALUES ($1, $2)
ON CONFLICT (one_row_id) DO UPDATE 
    SET params = excluded.params,
        height = excluded.height
WHERE staking_params.height <= excluded.height`

	_, err = db.SQL.Exec(stmt, string(paramsBz), params.Height)
	if err != nil {
		return fmt.Errorf("error while storing staking params: %s", err)
	}

	return nil
}

// GetStakingParams returns the types.StakingParams instance containing the current params
func (db *Db) GetStakingParams() (*types.StakingParams, error) {
	var rows []dbtypes.StakingParamsRow
	stmt := `SELECT * FROM staking_params LIMIT 1`
	err := db.Sqlx.Select(&rows, stmt)
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, fmt.Errorf("no staking params found")
	}

	var stakingParams stakingtypes.Params
	err = json.Unmarshal([]byte(rows[0].Params), &stakingParams)
	if err != nil {
		return nil, err
	}

	return &types.StakingParams{
		Params: stakingParams,
		Height: rows[0].Height,
	}, nil
}
