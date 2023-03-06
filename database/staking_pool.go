package database

import (
	"fmt"

	"github.com/forbole/bdjuno/v3/types"
)

// SaveStakingPool allows to save for the given height the given stakingtypes pool
func (db *Db) SaveStakingPool(pool *types.Pool) error {
	stmt := `
INSERT INTO staking_pool (bonded_tokens, not_bonded_tokens, height) 
VALUES ($1, $2, $3)
ON CONFLICT (one_row_id) DO UPDATE 
    SET bonded_tokens = excluded.bonded_tokens, 
        not_bonded_tokens = excluded.not_bonded_tokens, 
        height = excluded.height
WHERE staking_pool.height <= excluded.height`

	_, err := db.SQL.Exec(stmt, pool.BondedTokens.String(), pool.NotBondedTokens.String(), pool.Height)
	if err != nil {
		return fmt.Errorf("error while storing staking pool: %s", err)
	}

	return nil
}
