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

/*
//getValidatorUptime
func (db BigDipperDb) GetValidatorUptime(validator sdk.ConsAddress, height int64) (float64, error) {
	if found, _ := db.HasValidator(validator.String()); !found {
		return 0, nil
	}
	statement := `SELECT signed_blocks_window,missed_blocks_counter FROM validator_uptime WHERE validator_address = $1 and height = $2`
	rows, err := db.Sql.Query(statement, validator.String(), height)
	if err != nil {
		return _, err
	}
	var (
		signed_blocks_window  int64
		missed_blocks_counter int64
	)
	defer rows.Close()
	//supposed there is only one row(height is unique?)
	for rows.Next() {
		if err := rows.Scan(&signed_blocks_window, &missed_blocks_counter); err != nil {
			return _, err
		}
	}
	uptime := (float64(signed_blocks_window) - float64(missed_blocks_counter)) / float64(signed_blocks_window)

	return uptime, err
}
*/
type ValidatorCommission struct {
	ValidatorAddress sdk.ConsAddress `db:"validator_address"`
	Commission       sdk.Dec         `db:"new_commission"`
}

func (db BigDipperDb) SaveVaildatorComission(v ValidatorCommission) error {
	if found, _ := db.HasValidator(v.ValidatorAddress.String()); !found {
		return nil
	}
	statement := `INSERT INTO validator_commission (validatorAddress,commissions,timestamp) VALUES ($1,$2,$3)`
	_, err := db.Sql.Exec(statement,
		v.ValidatorAddress.String(), v.Commission.String(), time.Now().UTC())
	return err
}

//get array of
