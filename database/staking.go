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

// SaveValidatorUptime stores into the database the given validator uptime information
func (db BigDipperDb) SaveValidatorUptime(uptime bstaking.ValidatorUptime) error {
	if found, _ := db.HasValidator(uptime.ValidatorAddress.String()); !found {
		// Validator does not exist, return simply
		return nil
	}

	statement := `INSERT INTO validator_uptime (height, validator_address, signed_blocks_window, missed_blocks_counter)
				  VALUES ($1, $2, $3, $4)`
	_, err := db.Sql.Exec(statement,
		uptime.Height, uptime.ValidatorAddress.String(), uptime.SignedBlocksWindow, uptime.MissedBlocksCounter)

	if err != nil {
		return err
	}
	return nil
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

func (db BigDipperDb) SaveValidatorCommissions(validators []dbtypes.ValidatorCommission) error {
	query := `INSERT INTO validator_commission(validator_address,timestamp,commission,min_self_delegation,height) VALUES`
	var param []interface{}
	for i,validator := range(validators){
	vi:=i*5
		query += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d),",vi+1,vi+2,vi+3,vi+4,vi+5)
		param=append(param,validator.ValidatorAddress,validator.Timestamp.UTC,validator.Commission,
								validator.MinSelfDelegation,validator.Height)
	}
	query = query[:len(query)-1] // Remove trailing ","
	query += " ON CONFLICT DO NOTHING"
	_, err := db.Sql.Exec(query,param...)
	if err != nil {
		return err
	}
	return nil
}

func (db BigDipperDb) SaveValidatorInfo(validators []dbtypes.ValidatorCommission) error{
	query := `INSERT INTO validator_info(consensus_address,operator_address,moniker,identity,website,securityContact, details) VALUES`
	var param []interface
	for i,validator := range(validators){
	vi:=i*5
		query += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d),",vi+1,vi+2,vi+3,vi+4,vi+5)
		param=append(param,validator.ValidatorAddress,validator.Timestamp.UTC,validator.Commission,
								validator.MinSelfDelegation,validator.Height)
	}
	query = query[:len(query)-1] // Remove trailing ","
	query += " ON CONFLICT DO NOTHING"
	_, err := db.Sql.Exec(query,param...)
	if err != nil {
		return err
	}
	return nil
}

func (db BigDipperDb) SaveEditValidator(validator sdk.ValAddress, commissionRate int64, minSelfDelegation int64,
	description staking.Description, time time.Time, height int64) error {

	if found, _ := db.HasValidator(validator.String()); !found {
		return nil
	}
	//query the latest entry and see if the validator detail changed
	query := `SELECT commission,min_self_delegation,details,identity,moniker,website,securityContact
				FROM validator_commission INNER JOIN validator_info 
				ON  validator_commission.validator_address = validator_info.consensus_address
				WHERE timestamp = (
					SELECT MAX(timestamp) 
					FROM validator_commission
					WHERE validator_address = $1
				) and validator_address = $2 ;`
	var c int64
	var m int64
	var d [5]string
	discriptionArray := [5]string{description.Details, description.Identity, description.Moniker,
		description.Website, description.SecurityContact}

	rows, err1 := db.Sql.Query(query, validator.String(), validator.String())
	if err1 != nil {
		return err1
	}
	for rows.Next() {
		err := rows.Scan(&c, &m, &(d[0]), &(d[1]), &(d[2]), &(d[3]), &(d[4]))
		if err != nil {
			return err
		}
		if commissionRate == c || minSelfDelegation == m {
			//commission rate changed(insert)
			statement := `INSERT INTO validator_commission (validator_address,commissions,min_self_delegtion,height,timestamp) VALUES ($1,$2,$3,$4,$5);`
			_, err := db.Sql.Exec(statement,
				validator.String(), commissionRate, minSelfDelegation, height, time)
			if err != nil {
				return err
			}
		}
		if d == discriptionArray {
			//discription change(update)
			descriptionStatement := `UPDATE validator_info 
				SET Moniker     =$1,
				Identity        =$2,
				Website         =$3,
				SecurityContact =$4,
				Details         =$5
				WHERE consensus_address = $6;`
			_, err := db.Sql.Exec(descriptionStatement, d[0], d[1], d[2], d[3], d[4], validator.String())
			if err != nil {
				return err
			}
		}
	}

	return nil
}

//check the new commission is the same as the one before
func (db BigDipperDb) GetCommission(validator sdk.ValAddress) (ValidatorCommission, error) {

	var result []dbtypes.ValidatorCommission
	if found, _ := db.HasValidator(validator.String()); !found {
		return nil, nil
	}
	//query the latest entry and see if the validator detail changed
	query := `SELECT commission,min_self_delegation
				FROM validator_commission
				WHERE timestamp = (
					SELECT MAX(timestamp) 
					FROM validator_commission
					WHERE validator_address = $1
				) and validator_address = $2 ;`

	var c int64
	var m int64
	rows, err1 := db.Sql.Query(query, validator.String(), validator.String())
	if err := db.Sqlx.Select(&result, stmt, valAddress.String().valAddress.String()); err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("no validator with validator address %s could be found", valAddress.String())
	}
	return result[0], nil
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
