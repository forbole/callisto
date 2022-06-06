package database

import (
	"encoding/json"
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

// UpdatePoolSponsor updates the pool sponsor address
func (db *Db) UpdatePoolSponsor(poolID uint64, sponsor string, sponsorAddress string) error {
	stmt := `UPDATE shield_pool SET sponsor = $1 AND sponsor_address = $2 WHERE pool_id = %3`

	_, err := db.Sql.Exec(stmt, sponsor, sponsorAddress, poolID)
	if err != nil {
		return fmt.Errorf("error while updating shield pool sponsor: %s", err)
	}

	return nil
}

// SaveShieldPurchase allows to save shield purchase for the given height
func (db *Db) SaveShieldPurchase(shield *types.ShieldPurchase) error {
	stmt := `
INSERT INTO shield_purchase (pool_id, from_address, shield, description, height) 
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (from_address) DO UPDATE 
    SET pool_id = excluded.pool_id, 
	shield = excluded.shield, 
	description = excluded.description, 
    height = excluded.height
WHERE shield_purchase.height <= excluded.height`

	_, err := db.Sql.Exec(stmt,
		shield.PoolID,
		shield.FromAddress,
		pq.Array(dbtypes.NewDbCoins(shield.Shield)),
		shield.Description,
		shield.Height,
	)

	if err != nil {
		return fmt.Errorf("error while storing shield purchase: %s", err)
	}

	return nil
}

// SaveShieldPoolParams allows to save shield pool params
func (db *Db) SaveShieldPoolParams(params *types.ShieldPoolParams) error {
	paramsBz, err := json.Marshal(&params.Params)
	if err != nil {
		return fmt.Errorf("error while marshaling shield pool params: %s", err)
	}

	stmt := `
INSERT INTO shield_pool_params (params, height) 
VALUES ($1, $2)
ON CONFLICT (one_row_id) DO UPDATE 
    SET params = excluded.params,
        height = excluded.height
WHERE shield_pool_params.height <= excluded.height`

	_, err = db.Sql.Exec(stmt, string(paramsBz), params.Height)
	if err != nil {
		return fmt.Errorf("error while storing shield pool  params: %s", err)
	}

	return nil
}
