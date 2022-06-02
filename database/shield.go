package database

import (
	"fmt"

	dbtypes "github.com/forbole/bdjuno/v3/database/types"
	"github.com/forbole/bdjuno/v3/types"
	"github.com/lib/pq"
)

// SaveShieldPool allows to save for the given height the given shieldtypes pool
func (db *Db) SaveShieldPool(pool *types.ShieldPool) error {
	stmt := `
INSERT INTO shield_pool (pool_id, from_address, shield, native_deposit, foreign_deposit, sponsor, sponsor_address, description, shield_limit, pause, height) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
ON CONFLICT (pool_id) DO UPDATE 
    SET from_address = excluded.from_address, 
	shield = excluded.shield, 
	native_deposit = excluded.native_deposit, 
	foreign_deposit = excluded.foreign_deposit, 
	description = excluded.description, 
	shield_limit = excluded.shield_limit, 
    height = excluded.height
WHERE shield_pool.height <= excluded.height`

	_, err := db.Sql.Exec(stmt,
		pool.PoolID,
		pool.FromAddress,
		pq.Array(dbtypes.NewDbCoins(pool.Shield)),
		pq.Array(dbtypes.NewDbCoins(pool.NativeDeposit)),
		pq.Array(dbtypes.NewDbCoins(pool.ForeignDeposit)),
		pool.Sponsor,
		pool.SponsorAddr,
		pool.Description,
		pool.ShieldLimit.Int64(),
		pool.Pause,
		pool.Height,
	)

	if err != nil {
		return fmt.Errorf("error while storing shield pool: %s", err)
	}

	return nil
}

// UpdatePoolPauseStatus updates the pool pause status
func (db *Db) UpdatePoolPauseStatus(poolID uint64, pause bool) error {
	stmt := `UPDATE shield_pool SET pause = $1 WHERE pool_id = %2`

	_, err := db.Sql.Exec(stmt, pause, poolID)
	if err != nil {
		return fmt.Errorf("error while updating shield pool pause status: %s", err)
	}

	return nil
}
