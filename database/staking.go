package database

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	dbtypes "github.com/forbole/bdjuno/database/types"
	bstaking "github.com/forbole/bdjuno/x/staking/types"
)

// SaveStakingPool allows to save for the given height the given staking pool
func (db BigDipperDb) SaveStakingPool(height int64, pool staking.Pool) error {
	statement := `INSERT INTO staking_pool (timestamp, height, bonded_tokens, not_bonded_tokens) 
				  VALUES ($1, $2, $3, $4)`
	_, err := db.Sql.Exec(statement,
		time.Now().UTC(), height, pool.BondedTokens.Int64(), pool.NotBondedTokens.Int64())
	return err
}

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

// SaveValidatorsDelegations stores into the database the given validator delegations information.
// It assumes that for each delegation information provided, the associated validator data
// have already been saved inside the database properly.
func (db BigDipperDb) SaveValidatorsDelegations(
	delegationsInfo []bstaking.ValidatorDelegations, height int64, timestamp time.Time,
) error {
	var delegations []staking.DelegationResponse
	var unbondingDelegations []staking.UnbondingDelegation

	for _, info := range delegationsInfo {
		delegations = append(delegations, info.Delegations...)
		unbondingDelegations = append(unbondingDelegations, info.UnbondingDelegations...)
	}

	// Save the delegations
	err := db.saveValidatorDelegations(height, timestamp, delegations)
	if err != nil {
		return err
	}

	// Save the unbonding delegations
	return db.saveUnbondingDelegations(height, timestamp, unbondingDelegations)
}

func (db BigDipperDb) saveValidatorDelegations(
	height int64, timestamp time.Time, delegations []staking.DelegationResponse,
) error {
	accountsQuery := `INSERT INTO account (address) VALUES `
	var accountsParams []interface{}

	delegationsQuery := `INSERT INTO validator_delegations 
						 (consensus_address, delegator_address, shares, balance, height, timestamp) VALUES `
	var delegationsParams []interface{}

	for i, delegation := range delegations {
		a1 := i * 1
		d1 := i * 6 // Starting position for the delegations query

		accountsQuery += fmt.Sprintf("($%d),", a1+1)
		accountsParams = append(accountsParams, delegation.DelegatorAddress.String())

		var result []dbtypes.ValidatorInfoRow
		stmt := "SELECT consensus_address FROM validator_info WHERE operator_address = $1"
		err := db.sqlx.Select(&result, stmt, delegation.ValidatorAddress.String())
		if err != nil {
			return err
		}

		balance := dbtypes.NewDbCoin(delegation.Balance)
		coin, err := balance.Value()
		if err != nil {
			return err
		}

		delegationsQuery += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d),", d1+1, d1+2, d1+3, d1+4, d1+5, d1+6)
		delegationsParams = append(delegationsParams,
			result[0].ConsAddress, delegation.DelegatorAddress.String(), delegation.Shares.String(),
			coin, height, timestamp)
	}

	// Insert the accounts
	accountsQuery = accountsQuery[:len(accountsQuery)-1] // Remove the trailing ","
	accountsQuery += " ON CONFLICT DO NOTHING"
	_, err := db.Sql.Exec(accountsQuery, accountsParams...)
	if err != nil {
		return err
	}

	// Insert the delegations
	delegationsQuery = delegationsQuery[:len(delegationsQuery)-1] // Remove the trailing ","
	delegationsQuery += " ON CONFLICT DO NOTHING"
	_, err = db.Sql.Exec(delegationsQuery, delegationsParams...)
	return err
}

func (db BigDipperDb) saveUnbondingDelegations(
	height int64, timestamp time.Time, delegations []staking.UnbondingDelegation,
) error {
	// If the delegations are empty just return
	if len(delegations) == 0 {
		return nil
	}

	accountsQuery := `INSERT INTO account (address) VALUES `
	var accountsParams []interface{}

	udQuery := `INSERT INTO validator_unbonding_delegations 
			    (consensus_address, delegator_address, initial_balance, balance, creation_height, height, timestamp) VALUES `
	var delegationsParams []interface{}

	for i, delegation := range delegations {
		var result []dbtypes.ValidatorInfoRow
		stmt := "SELECT consensus_address FROM validator_info WHERE operator_address = $1"
		err := db.sqlx.Select(&result, stmt, delegation.ValidatorAddress.String())
		if err != nil {
			return err
		}

		for j, entry := range delegation.Entries {
			a1 := (i * 1) + j // Starting position for the account query
			d1 := (i * 7) + j // Starting position for the delegations query

			accountsQuery += fmt.Sprintf("($%d),", a1+1)
			accountsParams = append(accountsParams, delegation.DelegatorAddress.String())

			udQuery += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,$%d),", d1+1, d1+2, d1+3, d1+4, d1+5, d1+6, d1+7)
			delegationsParams = append(delegationsParams,
				result[0].ConsAddress, delegation.DelegatorAddress.String(),
				entry.InitialBalance.String(), entry.Balance.String(), entry.CreationHeight, height, timestamp)
		}
	}

	// Insert the accounts
	accountsQuery = accountsQuery[:len(accountsQuery)-1] // Remove the trailing ","
	accountsQuery += " ON CONFLICT DO NOTHING"
	_, err := db.Sql.Exec(accountsQuery, accountsParams...)
	if err != nil {
		return err
	}

	// Insert the delegations
	udQuery = udQuery[:len(udQuery)-1] // Remove the trailing ","
	udQuery += " ON CONFLICT DO NOTHING"
	_, err = db.Sql.Exec(udQuery, delegationsParams...)
	return err
}
