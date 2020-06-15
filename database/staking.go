package database

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
)

type ValidatorUptime struct {
	Height              int64           `db:"height"`
	ValidatorAddress    sdk.ConsAddress `db:"validator_address"`
	SignedBlocksWindow  int64           `db:"signed_blocks_window"`
	MissedBlocksCounter int64           `db:"missed_blocks_counter"`
}

// SaveStakingPool allows to save for the given height the given staking pool
func (db BigDipperDb) SaveStakingPool(height int64, pool staking.Pool) error {
	statement := `INSERT INTO staking_pool (timestamp, height, bonded_tokens, not_bonded_tokens) 
				  VALUES ($1, $2, $3, $4)`
	_, err := db.Sql.Exec(statement,
		time.Now().UTC(), height, pool.BondedTokens.Int64(), pool.NotBondedTokens.Int64())
	return err
}

// SaveValidatorUptime stores into the database the given validator uptime information
func (db BigDipperDb) SaveValidatorUptime(uptime ValidatorUptime) error {
	if found, _ := db.HasValidator(uptime.ValidatorAddress.String()); !found {
		// Validator does not exist, return simply
		return nil
	}

	statement := `INSERT INTO validator_uptime (height, validator_address, signed_blocks_window, missed_blocks_counter)
				  VALUES ($1, $2, $3, $4)`
	_, err := db.Sql.Exec(statement,
		uptime.Height, uptime.ValidatorAddress.String(), uptime.SignedBlocksWindow, uptime.MissedBlocksCounter)
	return err
}
