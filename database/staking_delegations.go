package database

import (
	"database/sql"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/forbole/bdjuno/v2/types"

	dbtypes "github.com/forbole/bdjuno/v2/database/types"
	dbutils "github.com/forbole/bdjuno/v2/database/utils"
)

// SaveDelegations stores inside the database the given delegations data.
// It assumes that the validators addresses are already present inside
// the proper database table.
// TIP: To store the validators data call SaveValidatorsData.
func (db *Db) SaveDelegations(height int64, delegations []types.Delegation) error {
	return db.RunTx(func(tx *sql.Tx) error {
		return db.saveDelegationsWithTx(tx, height, delegations)
	})
}

// saveDelegationsWithTx stores the given delegations using the given transaction.
// NOTE: The provided transaction is never committed nor rolled back. The caller must take care of this
func (db *Db) saveDelegationsWithTx(tx *sql.Tx, height int64, delegations []types.Delegation) error {
	if len(delegations) == 0 {
		return nil
	}

	paramsNumber := 4
	slices := dbutils.SplitDelegations(delegations, paramsNumber)

	for _, delegationSlice := range slices {
		if len(delegationSlice) == 0 {
			continue
		}

		err := db.storeUpToDateDelegations(tx, paramsNumber, height, delegationSlice)
		if err != nil {
			return fmt.Errorf("error while storing up-to-date delegations: %s", err)
		}
	}

	return nil
}

// storeUpToDateDelegations stores the given delegations as the most up-to-date ones using the provided transaction.
// NOTE: The provided transaction is never committed nor rolled back. The caller must take care of this
func (db *Db) storeUpToDateDelegations(tx *sql.Tx, paramsNumber int, height int64, delegations []types.Delegation) error {
	if len(delegations) == 0 {
		return nil
	}

	var accounts []types.Account

	delQry := `
INSERT INTO delegation (validator_address, delegator_address, amount, height) VALUES `
	var delParams []interface{}

	for i, delegation := range delegations {
		// Prepare the account query
		accounts = append(accounts, types.NewAccount(delegation.DelegatorAddress))

		// Get the validator consensus address
		consAddr, err := db.GetValidatorConsensusAddress(delegation.ValidatorOperAddr)
		if err != nil {
			return fmt.Errorf("error while gettting validator consensus address: %s", err)
		}

		// Convert the amount
		coin := dbtypes.NewDbCoin(delegation.Amount)
		value, err := coin.Value()
		if err != nil {
			return fmt.Errorf("error while converting coin to dbcoin: %s", err)
		}

		// Current delegation query
		di := i * paramsNumber
		delQry += fmt.Sprintf("($%d,$%d,$%d,$%d),", di+1, di+2, di+3, di+4)
		delParams = append(delParams, consAddr.String(), delegation.DelegatorAddress, value, height)
	}

	// Store the accounts
	err := db.saveAccountsWithTx(tx, accounts)
	if err != nil {
		return fmt.Errorf("error while storing delegators accounts: %s", err)
	}

	// Insert the delegations
	delQry = delQry[:len(delQry)-1] // Remove the trailing ","
	delQry += ` 
ON CONFLICT ON CONSTRAINT delegation_validator_delegator_unique 
DO UPDATE SET amount = excluded.amount, height = excluded.height
WHERE delegation.height <= excluded.height`

	_, err = tx.Exec(delQry, delParams...)
	if err != nil {
		return fmt.Errorf("error while storing delegations: %s", err)
	}
	return nil
}

// GetUserDelegationsAmount returns the amount of the delegations currently stored for
// the user having the given address
func (db *Db) GetUserDelegationsAmount(address string) (sdk.Coins, error) {
	stmt := `SELECT * FROM delegation WHERE delegator_address = $1`

	var rows []*dbtypes.DelegationRow
	err := db.Sqlx.Select(&rows, stmt, address)
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return sdk.Coins{}, nil
	}

	var amount = sdk.Coins{}
	for _, delegation := range rows {
		amount = amount.Add(delegation.Amount.ToCoin())
	}

	return amount, nil
}

// ReplaceValidatorDelegations replaces all the delegations associated with the given validator with the given ones
func (db *Db) ReplaceValidatorDelegations(height int64, valOperAddr string, delegations []types.Delegation) error {
	return db.RunTx(func(tx *sql.Tx) error {
		// Delete the existing delegations
		stmt := `
DELETE FROM delegation USING validator_info 
WHERE delegation.validator_address = validator_info.consensus_address 
  AND validator_info.operator_address = $1 AND delegation.height <= $2`
		_, err := tx.Exec(stmt, valOperAddr, height)
		if err != nil {
			return fmt.Errorf("error while deleting delegations for validator: %s", err)
		}

		// Store the delegations
		return db.saveDelegationsWithTx(tx, height, delegations)
	})
}

// ReplaceDelegatorDelegations replaces all the delegations associated with the given delegator with the given ones
func (db *Db) ReplaceDelegatorDelegations(height int64, delegator string, delegations []types.Delegation) error {
	return db.RunTx(func(tx *sql.Tx) error {
		// Delete existing delegations
		stmt := `DELETE FROM delegation WHERE delegator_address = $1 AND height <= $2`
		_, err := tx.Exec(stmt, delegator, height)
		if err != nil {
			return fmt.Errorf("error while deleting delegations for delegator: %s", err)
		}

		// Store the delegations
		return db.saveDelegationsWithTx(tx, height, delegations)
	})

}

// --------------------------------------------------------------------------------------------------------------------

// GetDelegators returns the current delegators set
func (db *Db) GetDelegators() ([]string, error) {
	var rows []string
	err := db.Sqlx.Select(&rows, `SELECT DISTINCT (delegator_address) FROM delegation `)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

// --------------------------------------------------------------------------------------------------------------------

// SaveRedelegations saves the given redelegations inside the database.
// It assumes that all the validators as well as all the delegators addresses are
// already present inside the proper tables of the database.
// To store the validators data call SaveValidatorData(s).
// To store the account data call SaveAccount.
func (db *Db) SaveRedelegations(redelegations []types.Redelegation) error {
	paramsNumber := 6
	slices := dbutils.SplitRedelegations(redelegations, paramsNumber)

	for _, redelegationSlice := range slices {
		if len(redelegationSlice) == 0 {
			continue
		}

		err := db.storeUpToDateRedelegations(paramsNumber, redelegationSlice)
		if err != nil {
			return fmt.Errorf("error while storing up-to-date redelegations: %s", err)
		}
	}

	return nil
}

// storeUpToDateRedelegations allows to store the given redelegations as the most up-to-date ones
func (db *Db) storeUpToDateRedelegations(paramsNumber int, redelegations []types.Redelegation) error {

	if len(redelegations) == 0 {
		return nil
	}

	var accounts []types.Account

	rdQry := `
INSERT INTO redelegation 
    (delegator_address, src_validator_address, dst_validator_address, amount, completion_time, height) 
VALUES `
	var rdParams []interface{}

	for i, redelegation := range redelegations {
		// Prepare the account query
		accounts = append(accounts, types.NewAccount(redelegation.DelegatorAddress))

		// Get the validators info
		srcVal, err := db.GetValidator(redelegation.SrcValidator)
		if err != nil {
			return fmt.Errorf("error while getting validator: %s", err)
		}

		dstVal, err := db.GetValidator(redelegation.DstValidator)
		if err != nil {
			return fmt.Errorf("error while getting validator: %s", err)
		}

		// Convert the amount value
		coin := dbtypes.NewDbCoin(redelegation.Amount)
		amountValue, err := coin.Value()
		if err != nil {
			return fmt.Errorf("error while converting coin to dbcoin: %s", err)
		}

		rdi := i * paramsNumber
		rdQry += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d),", rdi+1, rdi+2, rdi+3, rdi+4, rdi+5, rdi+6)
		rdParams = append(rdParams,
			redelegation.DelegatorAddress,
			srcVal.GetConsAddr(), dstVal.GetConsAddr(), amountValue, redelegation.CompletionTime, redelegation.Height)
	}

	// Store the accounts
	err := db.SaveAccounts(accounts)
	if err != nil {
		return fmt.Errorf("error while storing redelegators accounts: %s", err)
	}

	// Insert the redelegations
	rdQry = rdQry[:len(rdQry)-1] // Remove the trailing ","
	rdQry += `
ON CONFLICT ON CONSTRAINT redelegation_validator_delegator_unique 
DO UPDATE SET height = excluded.height
WHERE redelegation.height <= excluded.height`

	_, err = db.Sql.Exec(rdQry, rdParams...)
	if err != nil {
		return fmt.Errorf("error while storing redelegations: %s", err)
	}

	return nil
}

// GetUserRedelegationsAmount returns the amount of the redelegations currently stored for
// the user having the given address
func (db *Db) GetUserRedelegationsAmount(address string) (sdk.Coins, error) {
	stmt := `SELECT * FROM redelegation WHERE delegator_address = $1`

	var rows []*dbtypes.RedelegationRow
	err := db.Sqlx.Select(&rows, stmt, address)
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return sdk.Coins{}, nil
	}

	var amount = sdk.Coins{}
	for _, delegation := range rows {
		amount = amount.Add(delegation.Amount.ToCoin())
	}

	return amount, nil
}

// DeleteCompletedRedelegations deletes all the redelegations
// that have completed before the given timestamp
func (db *Db) DeleteCompletedRedelegations(timestamp time.Time) error {
	stmt := `DELETE FROM redelegation WHERE completion_time < $1`
	_, err := db.Sql.Exec(stmt, timestamp)
	return err
}

// --------------------------------------------------------------------------------------------------------------------

// SaveUnbondingDelegations saves the given unbonding delegations into the database.
// It assumes that all the validators as well as all the delegators addresses are
// already present inside the proper tables of the database.
// To store the validators data call SaveValidatorData(s).
// To store the account data call SaveAccount.
func (db *Db) SaveUnbondingDelegations(delegations []types.UnbondingDelegation) error {
	paramsNumber := 5
	slices := dbutils.SplitUnbondingDelegations(delegations, paramsNumber)

	for _, delegationSlice := range slices {
		if len(delegationSlice) == 0 {
			continue
		}

		err := db.storeUpToDateUnbondingDelegations(paramsNumber, delegationSlice)
		if err != nil {
			return fmt.Errorf("error while storing up-to-date undonding delegations: %s", err)
		}
	}

	return nil
}

// storeUpToDateUnbondingDelegations allows to store the given unbonding delegations as the most up-to-date ones
func (db *Db) storeUpToDateUnbondingDelegations(paramsNumber int, delegations []types.UnbondingDelegation) error {
	if len(delegations) == 0 {
		return nil
	}

	var accounts []types.Account

	udQry := `
INSERT INTO unbonding_delegation (validator_address, delegator_address, amount, completion_timestamp, height)
VALUES `
	var udParams []interface{}

	for i, delegation := range delegations {
		// Prepare the account query
		accounts = append(accounts, types.NewAccount(delegation.DelegatorAddress))

		validator, err := db.GetValidator(delegation.ValidatorOperAddr)
		if err != nil {
			return fmt.Errorf("error while getting validator: %s", err)
		}

		coin := dbtypes.NewDbCoin(delegation.Amount)
		amount, err := coin.Value()
		if err != nil {
			return fmt.Errorf("error while converting coin to dbcoin: %s", err)
		}

		udi := i * paramsNumber
		udQry += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d),", udi+1, udi+2, udi+3, udi+4, udi+5)
		udParams = append(udParams,
			validator.GetConsAddr(), delegation.DelegatorAddress, amount, delegation.CompletionTimestamp, delegation.Height)
	}

	// Store the accounts
	err := db.SaveAccounts(accounts)
	if err != nil {
		return fmt.Errorf("error while storing unbonding delegators accounts: %s", err)
	}

	// Insert the current unbonding delegations
	udQry = udQry[:len(udQry)-1] // Remove the trailing ","
	udQry += `
ON CONFLICT ON CONSTRAINT unbonding_delegation_validator_delegator_unique 
DO UPDATE SET height = excluded.height
WHERE unbonding_delegation.height <= excluded.height`

	_, err = db.Sql.Exec(udQry, udParams...)
	if err != nil {
		return fmt.Errorf("error while storing unbonding delegations: %s", err)
	}

	return nil
}

// GetUserUnBondingDelegationsAmount returns the amount of the redelegations currently stored for
// the user having the given address
func (db *Db) GetUserUnBondingDelegationsAmount(address string) (sdk.Coins, error) {
	stmt := `SELECT * FROM unbonding_delegation WHERE delegator_address = $1`

	var rows []*dbtypes.UnbondingDelegationRow
	err := db.Sqlx.Select(&rows, stmt, address)
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return sdk.Coins{}, nil
	}

	var amount = sdk.Coins{}
	for _, delegation := range rows {
		amount = amount.Add(delegation.Amount.ToCoin())
	}

	return amount, nil
}

// DeleteCompletedUnbondingDelegations deletes all the unbonding delegations
// that have completed before the given timestamp
func (db *Db) DeleteCompletedUnbondingDelegations(timestamp time.Time) error {
	stmt := `DELETE FROM unbonding_delegation WHERE completion_timestamp < $1 RETURNING *`
	_, err := db.Sql.Exec(stmt, timestamp)
	return err
}

// --------------------------------------------------------------------------------------------------------------------

// SaveDelegatorsToRefresh stores the given delegators as accounts to be refreshed on the block at height + 1
func (db *Db) SaveDelegatorsToRefresh(height int64, delegators []string) error {
	if len(delegators) == 0 {
		return nil
	}

	dQry := `INSERT INTO delegators_to_refresh (address, height) VALUES `

	var params []interface{}
	for i, delegator := range delegators {
		di := i * 2

		dQry += fmt.Sprintf("($%d, $%d),", di+1, di+2)
		params = append(params, delegator, height)
	}

	dQry = dQry[:len(dQry)-1] // Remove the trailing ","
	dQry += `
ON CONFLICT ON CONSTRAINT unique_address
DO UPDATE SET height = excluded.height
WHERE delegators_to_refresh.height <= excluded.height`

	_, err := db.Sql.Exec(dQry, params...)
	return err
}

// DeleteDelegatorsToRefresh removes and returns the list of delegators whose balance
// should be updated at the given height
func (db *Db) DeleteDelegatorsToRefresh(height int64) ([]string, error) {
	var delegators []string
	stmt := `DELETE FROM delegators_to_refresh WHERE height < $1 RETURNING address`
	err := db.Sqlx.Select(&delegators, stmt, height)
	if err != nil {
		return nil, err
	}
	return delegators, err
}
