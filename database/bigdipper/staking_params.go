package bigdipper

import (
	"fmt"

	dbtypes "github.com/forbole/bdjuno/database/types"
	"github.com/forbole/bdjuno/types"
)

// SaveStakingParams allows to store the given params into the database
func (db *Db) SaveStakingParams(params types.StakingParams) error {
	stmt := `
INSERT INTO staking_params (bond_denom) 
VALUES ($1)
ON CONFLICT (one_row_id) DO UPDATE 
    SET bond_denom = excluded.bond_denom`

	_, err := db.Sql.Exec(stmt, params.BondName)
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

	return &types.StakingParams{
		BondName: rows[0].BondName,
	}, nil
}
