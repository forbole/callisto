package database

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	dbtypes "github.com/forbole/bdjuno/database/types"
	bstaking "github.com/forbole/bdjuno/x/staking/types"
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


// GetAccounts returns all the accounts that are currently stored inside the database.
func (db BigDipperDb) GetValidators() ([]bstaking.Validator, error) {
	sqlStmt := `SELECT DISTINCT ON (validator.consensus_address)
				validator.consensus_address, validator.consensus_pubkey, validator_info.operator_address 
				FROM validator 
				INNER JOIN validator_info 
				ON validator.consensus_address = validator_info.consensus_address`

	var rows []dbtypes.ValidatorInfoRow
	err := db.sqlx.Select(&rows, sqlStmt)
	if err != nil {
		return nil, err
	}

	validators := make([]bstaking.Validator, len(rows))
	for index, row := range rows {
		consAddress, err := sdk.ConsAddressFromBech32(row.ConsAddress)
		if err != nil {
			return nil, err
		}

		consPubKey, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, row.ConsPubKey)
		if err != nil {
			return nil, err
		}

		valAddress, err := sdk.ValAddressFromBech32(row.ValAddress)
		if err != nil {
			return nil, err
		}

		validators[index] = dbtypes.ValidatorData{
			ConsAddress: consAddress,
			ValAddress:  valAddress,
			ConsPubKey:  consPubKey,
		}
	}

	return validators, nil
}

// SaveValidators allows the bulk saving of a list of validators
func (db BigDipperDb) SaveValidators(validators []bstaking.Validator) error {
	validatorQuery := `INSERT INTO validator (consensus_address, consensus_pubkey) VALUES `
	var validatorParams []interface{}

	validatorInfoQuery := `INSERT INTO validator_info (consensus_address, operator_address) VALUES `
	var validatorInfoParams []interface{}

	for i, validator := range validators {
		v1 := i * 2 // Starting position for validator params
		i1 := i * 2 // Starting position for validator info params

		publicKey, err := sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeConsPub, validator.GetConsPubKey())
		if err != nil {
			return err
		}

		validatorQuery += fmt.Sprintf("($%d,$%d),", v1+1, v1+2)
		validatorParams = append(validatorParams,
			validator.GetConsAddr().String(), publicKey)

		validatorInfoQuery += fmt.Sprintf("($%d,$%d),", i1+1, i1+2)
		validatorInfoParams = append(validatorInfoParams,
			validator.GetConsAddr().String(), validator.GetOperator().String())
	}

	validatorQuery = validatorQuery[:len(validatorQuery)-1] // Remove trailing ","
	validatorQuery += " ON CONFLICT DO NOTHING"
	_, err := db.Sql.Exec(validatorQuery, validatorParams...)
	if err != nil {
		return err
	}

	validatorInfoQuery = validatorInfoQuery[:len(validatorInfoQuery)-1] // Remove the trailing ","
	validatorInfoQuery += " ON CONFLICT DO NOTHING"
	_, err = db.Sql.Exec(validatorInfoQuery, validatorInfoParams...)
	return err
}

// SaveValidatorUptime stores into the database the given validator uptime information.
// It assumes that for each uptime information provided, the associated validator data
// have already been saved inside the database properly.
func (db BigDipperDb) SaveValidatorUptime(uptime bstaking.ValidatorUptime) error {
	statement := `INSERT INTO validator_uptime (height, validator_address, signed_blocks_window, missed_blocks_counter)
				  VALUES ($1, $2, $3, $4)`
	_, err := db.Sql.Exec(statement,
		uptime.Height, uptime.ValidatorAddress.String(), uptime.SignedBlocksWindow, uptime.MissedBlocksCounter)
	return err
}
