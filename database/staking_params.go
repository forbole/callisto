package database

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	dbtypes "github.com/forbole/bdjuno/database/types"
	"github.com/forbole/bdjuno/types"
)

// SaveStakingParams allows to store the given params into the database
func (db *Db) SaveStakingParams(params types.StakingParams) error {
	stmt := `
INSERT INTO staking_params (bond_denom, unbonding_time, max_entries, historical_entries, max_validators, min_commission_rate, height) 
VALUES ($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT (one_row_id) DO UPDATE 
    SET bond_denom = excluded.bond_denom,
    	unbonding_time = excluded.unbonding_time,
    	max_entries = excluded.max_entries,
    	historical_entries = excluded.historical_entries,
    	max_validators = excluded.max_validators,
        min_commission_rate = excluded.min_commission_rate,
        height = excluded.height
WHERE staking_params.height <= excluded.height`

	_, err := db.Sql.Exec(stmt,
		params.BondDenom, params.UnbondingTime.Nanoseconds(), params.MaxEntries,
		params.HistoricalEntries, params.MaxValidators, params.MinCommissionRate.String(), params.Height)
	return err
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

	commissionRate, err := sdk.NewDecFromStr(rows[0].MinCommissionRate)
	if err != nil {
		return nil, err
	}

	return &types.StakingParams{
		Params: stakingtypes.NewParams(
			time.Duration(rows[0].UnbondingTime),
			rows[0].MaxValidators,
			rows[0].MaxEntries,
			rows[0].HistoricalEntries,
			rows[0].BondName,
			commissionRate,
		),
		Height: rows[0].Height,
	}, nil
}
