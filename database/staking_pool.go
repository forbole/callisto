package database

import (
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// SaveStakingPool allows to save for the given height the given stakingtypes pool
func (db *BigDipperDb) SaveStakingPool(pool stakingtypes.Pool) error {
	stmt := `DELETE FROM staking_pool WHERE TRUE`
	_, err := db.Sql.Exec(stmt)
	if err != nil {
		return err
	}

	stmt = `INSERT INTO staking_pool (bonded_tokens, not_bonded_tokens) VALUES ($1, $2)`
	_, err = db.Sql.Exec(stmt, pool.BondedTokens.Int64(), pool.NotBondedTokens.Int64())
	return err
}
