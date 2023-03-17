package database

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	dbtypes "github.com/forbole/bdjuno/v3/database/types"
	"github.com/forbole/bdjuno/v3/types"
	"github.com/lib/pq"
)

// SaveShieldPool allows to save for the given height the given shieldtypes pool
func (db *Db) SaveShieldPool(pool *types.ShieldPool) error {
	stmt := `
INSERT INTO shield_pool (pool_id, shield, native_service_fees, foreign_service_fees, sponsor, sponsor_address, description, shield_limit, pause, height) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
ON CONFLICT (pool_id) DO UPDATE 
    SET shield = excluded.shield, 
	native_service_fees = excluded.native_service_fees, 
	foreign_service_fees = excluded.foreign_service_fees, 
	description = excluded.description, 
	shield_limit = excluded.shield_limit, 
    height = excluded.height
WHERE shield_pool.height <= excluded.height`

	_, err := db.SQL.Exec(stmt,
		pool.PoolID,
		pool.Shield.String(),
		pq.Array(dbtypes.NewDbCoins(pool.NativeServiceFees)),
		pq.Array(dbtypes.NewDbCoins(pool.ForeignServiceFees)),
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
func (db *Db) UpdatePoolPauseStatus(poolID uint64, pause bool, height int64) error {
	stmt := `UPDATE shield_pool SET pause = $1, height = $2 WHERE pool_id = %3`

	_, err := db.SQL.Exec(stmt, pause, height, poolID)
	if err != nil {
		return fmt.Errorf("error while updating shield pool pause status: %s", err)
	}

	return nil
}

// UpdatePoolSponsor updates the pool sponsor address
func (db *Db) UpdatePoolSponsor(poolID uint64, sponsor string, sponsorAddress string, height int64) error {
	stmt := `UPDATE shield_pool SET sponsor = $1, sponsor_address = $2, height = $3 WHERE pool_id = %4`

	_, err := db.SQL.Exec(stmt, sponsor, sponsorAddress, height, poolID)
	if err != nil {
		return fmt.Errorf("error while updating shield pool sponsor: %s", err)
	}

	return nil
}

// UpdateShieldProviderCollateral updates the shield provider' collateral value
func (db *Db) UpdateShieldProviderCollateral(address string, collateral int64, height int64) error {
	stmt := `UPDATE shield_provider SET collateral = $1, height = $2 WHERE address = $3`

	_, err := db.SQL.Exec(stmt, collateral, height, address)
	if err != nil {
		return fmt.Errorf("error while updating shield provider collateral value: %s", err)
	}

	return nil
}

// WithdrawNativeRewards withdraws the shield provider' native rewards
func (db *Db) WithdrawNativeRewards(address string, height int64) error {
	stmt := `UPDATE shield_provider SET native_rewards = $1, height = $2 WHERE address = $3`

	_, err := db.SQL.Exec(stmt, pq.Array(dbtypes.NewDbDecCoins(sdk.DecCoins{})), height, address)
	if err != nil {
		return fmt.Errorf("error while withdrawing the native rewards: %s", err)
	}

	return nil
}

// WithdrawForeignRewards withdraws the shield provider' foreign rewards
func (db *Db) WithdrawForeignRewards(address string, height int64) error {
	stmt := `UPDATE shield_provider SET foreign_rewards = $1, height = $2 WHERE address = $3`

	_, err := db.SQL.Exec(stmt, pq.Array(dbtypes.NewDbDecCoins(sdk.DecCoins{})), height, address)
	if err != nil {
		return fmt.Errorf("error while withdrawing the foreign rewards: %s", err)
	}

	return nil
}

// UpdateShieldProviderDelegation updates the shield provider' delegation value
func (db *Db) UpdateShieldProviderDelegation(address string, delegation int64, height int64) error {
	stmt := `UPDATE shield_provider SET delegation_bonded = $1, height = $2 WHERE address = $3`

	_, err := db.SQL.Exec(stmt, delegation, height, address)
	if err != nil {
		return fmt.Errorf("error while updating shield provider delegation value: %s", err)
	}

	return nil
}

// GetShieldProviderCollateral returns the shield provider' collateral value
func (db *Db) GetShieldProviderCollateral(address string) (int64, error) {
	var collateral []int64
	stmt := `SELECT collateral from shield_provider WHERE address = $1`

	err := db.Sqlx.Select(&collateral, stmt, address)
	if err != nil {
		return 0, fmt.Errorf("error while getting shield provider collateral value: %s", err)
	}

	if len(collateral) == 0 {
		return 0, fmt.Errorf("cannot find shield provider for address %s", address)
	}

	return collateral[0], nil
}

// GetShieldProviderDelegation returns the shield provider' delegation value
func (db *Db) GetShieldProviderDelegation(address string) (int64, error) {
	var delegation int64
	stmt := `SELECT delegation_bonded from shield_provider WHERE address = $1`

	err := db.Sqlx.Select(&delegation, stmt, address)
	if err != nil {
		return 0, fmt.Errorf("error while getting shield provider delegation value: %s", err)
	}

	return delegation, nil
}

// SaveShieldPurchase allows to save shield purchase for the given height
func (db *Db) SaveShieldPurchase(shield *types.ShieldPurchase) error {
	stmt := `
INSERT INTO shield_purchase (pool_id, purchaser, shield, description, height)
VALUES ($1, $2, $3, $4, $5)`

	_, err := db.SQL.Exec(stmt,
		shield.PoolID,
		shield.FromAddress,
		shield.Shield.String(),
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

	_, err = db.SQL.Exec(stmt, string(paramsBz), params.Height)
	if err != nil {
		return fmt.Errorf("error while storing shield pool params: %s", err)
	}

	return nil
}

// SaveShieldClaimProposalParams allows to save shield claim proposal params
func (db *Db) SaveShieldClaimProposalParams(params *types.ShieldClaimProposalParams) error {
	paramsBz, err := json.Marshal(&params.Params)
	if err != nil {
		return fmt.Errorf("error while marshaling shield claim proposal params: %s", err)
	}

	stmt := `
INSERT INTO shield_claim_proposal_params (params, height) 
VALUES ($1, $2)
ON CONFLICT (one_row_id) DO UPDATE 
    SET params = excluded.params,
        height = excluded.height
WHERE shield_claim_proposal_params.height <= excluded.height`

	_, err = db.SQL.Exec(stmt, string(paramsBz), params.Height)
	if err != nil {
		return fmt.Errorf("error while storing shield claim proposal params: %s", err)
	}

	return nil
}

// SaveShieldProvider allows to save the shield provider for the given height
func (db *Db) SaveShieldProvider(provider *types.ShieldProvider) error {
	stmt := `
INSERT INTO shield_provider (address, collateral, delegation_bonded, native_rewards, 
    foreign_rewards,total_locked, withdrawing, height) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
ON CONFLICT (address) DO UPDATE 
    SET collateral = excluded.collateral, 
	delegation_bonded = excluded.delegation_bonded, 
	native_rewards = excluded.native_rewards, 
	foreign_rewards = excluded.foreign_rewards, 
	total_locked = excluded.total_locked, 
	withdrawing = excluded.withdrawing, 
    height = excluded.height
WHERE shield_provider.height <= excluded.height`

	_, err := db.SQL.Exec(stmt,
		provider.Address,
		provider.Collateral,
		provider.DelegationBonded,
		pq.Array(dbtypes.NewDbDecCoins(provider.NativeRewards)),
		pq.Array(dbtypes.NewDbDecCoins(provider.ForeignRewards)),
		provider.TotalLocked,
		provider.Withdrawing,
		provider.Height,
	)

	if err != nil {
		return fmt.Errorf("error while storing shield provider: %s", err)
	}

	return nil
}

// SaveShieldWithdraw allows to save the shield withdraw for the given height
func (db *Db) SaveShieldWithdraw(withdraw *types.ShieldWithdraw) error {
	stmt := `
INSERT INTO shield_withdraws (address, amount, completion_time, height)
VALUES ($1, $2, $3, $4)`

	_, err := db.SQL.Exec(stmt,
		withdraw.Address,
		withdraw.Amount,
		withdraw.CompletionTime,
		withdraw.Height,
	)

	if err != nil {
		return fmt.Errorf("error while storing shield withdraw: %s", err)
	}

	return nil
}

// SaveShieldStatus allows to save the shield status for the given height
func (db *Db) SaveShieldStatus(status *types.ShieldStatus) error {
	stmt := `
INSERT INTO shield_status (global_staking_pool, current_native_service_fees, current_foreign_service_fees, remaining_native_service_fees,
	remaining_foreign_service_fees, total_collateral, total_shield, total_withdrawing, height) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
ON CONFLICT (one_row_id) DO UPDATE 
    SET global_staking_pool = excluded.global_staking_pool, 
	current_native_service_fees = excluded.current_native_service_fees,
	current_foreign_service_fees = excluded.current_foreign_service_fees,
	remaining_native_service_fees = excluded.remaining_native_service_fees,
	remaining_foreign_service_fees = excluded.remaining_foreign_service_fees, 
	total_collateral = excluded.total_collateral, 
	total_shield = excluded.total_shield, 
	total_withdrawing = excluded.total_withdrawing, 
    height = excluded.height
WHERE shield_status.height <= excluded.height`

	_, err := db.SQL.Exec(stmt,
		status.GobalStakingPool.Int64(),
		pq.Array(dbtypes.NewDbDecCoins(status.CurrentNativeServiceFees)),
		pq.Array(dbtypes.NewDbDecCoins(status.CurrentForeignServiceFees)),
		pq.Array(dbtypes.NewDbDecCoins(status.RemainingNativeServiceFees)),
		pq.Array(dbtypes.NewDbDecCoins(status.RemainingForeignServiceFees)),
		status.TotalCollateral.Int64(),
		status.TotalShield.Int64(),
		status.TotalWithdrawing.Int64(),
		status.Height,
	)

	if err != nil {
		return fmt.Errorf("error while storing shield status: %s", err)
	}

	return nil
}
