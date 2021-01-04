package database

import (
	"time"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// SaveStakingPool allows to save for the given height the given stakingtypes pool
func (db *BigDipperDb) SaveStakingPool(pool stakingtypes.Pool, height int64, timestamp time.Time) error {
	stmt := `
INSERT INTO staking_pool_history (timestamp, height, bonded_tokens, not_bonded_tokens) 
VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`

	_, err := db.Sql.Exec(stmt,
		timestamp, height, pool.BondedTokens.Int64(), pool.NotBondedTokens.Int64())
	return err
}
