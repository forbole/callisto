package database

import (
	"fmt"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"time"

	dbtypes "github.com/forbole/bdjuno/database/types"
)

// SaveStakingParams allows to store the given params into the database
func (db *Db) SaveStakingParams(params stakingtypes.Params) error {
	stmt := `
INSERT INTO staking_params (bond_denom, unbonding_time, max_entries, historical_entries, max_validators) 
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (one_row_id) DO UPDATE 
    SET bond_denom = excluded.bond_denom,
    	unbonding_time = excluded.unbonding_time,
    	max_entries = excluded.max_entries,
    	historical_entries = excluded.historical_entries,
    	max_validators = excluded.max_validators`

	_, err := db.Sql.Exec(stmt,
		params.BondDenom, params.UnbondingTime.Nanoseconds(), params.MaxEntries, params.HistoricalEntries, params.MaxValidators)
	return err
}

// GetStakingParams returns the types.StakingParams instance containing the current params
func (db *Db) GetStakingParams() (*stakingtypes.Params, error) {
	var rows []dbtypes.StakingParamsRow
	stmt := `SELECT * FROM staking_params LIMIT 1`
	err := db.Sqlx.Select(&rows, stmt)
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, fmt.Errorf("no staking params found")
	}

	return &stakingtypes.Params{
		UnbondingTime:     time.Duration(rows[0].UnbondingTime),
		MaxValidators:     rows[0].MaxValidators,
		MaxEntries:        rows[0].MaxEntries,
		HistoricalEntries: rows[0].HistoricalEntries,
		BondDenom:         rows[0].BondName,
	}, nil
}
