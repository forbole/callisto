package database

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/forbole/bdjuno/types"

	dbtypes "github.com/forbole/bdjuno/database/types"
)

// SaveDelegations stores inside the database the given delegations data.
// It assumes that the validators addresses are already present inside
// the proper database table.
// TIP: To store the validators data call SaveValidatorsData.
func (db *Db) SaveDelegations(delegations []types.Delegation) error {
	if len(delegations) == 0 {
		return nil
	}

	err := db.storeUpToDateDelegations(delegations)
	if err != nil {
		return fmt.Errorf("error while storing up-to-date delegations: %s", err)
	}

	return nil
}

// storeUpToDateDelegations stores the given delegations as the most up-to-date ones
func (db *Db) storeUpToDateDelegations(delegations []types.Delegation) error {
	accQry := `
INSERT INTO account (address) VALUES `
	var accParams []interface{}

	delQry := `
INSERT INTO delegation (validator_address, delegator_address, amount, height) VALUES `
	var delParams []interface{}

	for i, delegation := range delegations {
		ai := i * 1
		accQry += fmt.Sprintf("($%d),", ai+1)
		accParams = append(accParams, delegation.DelegatorAddress)

		// Get the validator consensus address
		consAddr, err := db.GetValidatorConsensusAddress(delegation.ValidatorOperAddr)
		if err != nil {
			return err
		}

		// Convert the amount
		coin := dbtypes.NewDbCoin(delegation.Amount)
		value, err := coin.Value()
		if err != nil {
			return err
		}

		// Current delegation query
		di := i * 4
		delQry += fmt.Sprintf("($%d,$%d,$%d,$%d),", di+1, di+2, di+3, di+4)
		delParams = append(delParams,
			consAddr.String(), delegation.DelegatorAddress, value, delegation.Height)
	}

	// Insert the accounts
	accQry = accQry[:len(accQry)-1] // Remove the trailing ","
	accQry += " ON CONFLICT DO NOTHING"
	_, err := db.Sql.Exec(accQry, accParams...)
	if err != nil {
		return err
	}

	// Insert the delegations
	delQry = delQry[:len(delQry)-1] // Remove the trailing ","
	delQry += ` 
ON CONFLICT ON CONSTRAINT delegation_validator_delegator_unique 
DO UPDATE SET amount = excluded.amount, height = excluded.height
WHERE delegation.height <= excluded.height`
	_, err = db.Sql.Exec(delQry, delParams...)
	return err
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

// DeleteDelegatorDelegations removes all the delegations associated with the given delegator
func (db *Db) DeleteDelegatorDelegations(delegator string) error {
	stmt := `DELETE FROM delegation WHERE delegator_address = $1`
	_, err := db.Sql.Exec(stmt, delegator)
	return err
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
	if len(redelegations) == 0 {
		return nil
	}

	err := db.storeUpToDateRedelegations(redelegations)
	if err != nil {
		return fmt.Errorf("error while storing up-to-date redelegations: %s", err)
	}

	return nil
}

// storeUpToDateRedelegations allows to store the given redelegations as the most up-to-date ones
func (db *Db) storeUpToDateRedelegations(redelegations []types.Redelegation) error {
	accQry := `
INSERT INTO account (address) VALUES `
	var accParams []interface{}

	rdQry := `
INSERT INTO redelegation 
    (delegator_address, src_validator_address, dst_validator_address, amount, completion_time, height) 
VALUES `
	var rdParams []interface{}

	for i, redelegation := range redelegations {
		a1 := i * 1
		accQry += fmt.Sprintf("($%d),", a1+1)
		accParams = append(accParams, redelegation.DelegatorAddress)

		// Get the validators info
		srcVal, err := db.GetValidator(redelegation.SrcValidator)
		if err != nil {
			return err
		}

		dstVal, err := db.GetValidator(redelegation.DstValidator)
		if err != nil {
			return err
		}

		// Convert the amount value
		coin := dbtypes.NewDbCoin(redelegation.Amount)
		amountValue, err := coin.Value()
		if err != nil {
			return err
		}

		rdi := i * 6
		rdQry += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d),", rdi+1, rdi+2, rdi+3, rdi+4, rdi+5, rdi+6)
		rdParams = append(rdParams,
			redelegation.DelegatorAddress,
			srcVal.GetConsAddr(), dstVal.GetConsAddr(), amountValue, redelegation.CompletionTime, redelegation.Height)
	}

	// Insert the delegators
	accQry = accQry[:len(accQry)-1] // Remove the trailing ","
	accQry += " ON CONFLICT DO NOTHING"
	_, err := db.Sql.Exec(accQry, accParams...)
	if err != nil {
		return err
	}

	// Insert the redelegations
	rdQry = rdQry[:len(rdQry)-1] // Remove the trailing ","
	rdQry += `
ON CONFLICT ON CONSTRAINT redelegation_validator_delegator_unique 
DO UPDATE SET height = excluded.height
WHERE redelegation.height <= excluded.height`
	_, err = db.Sql.Exec(rdQry, rdParams...)
	return err
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

// DeleteRedelegation removes the given redelegation from the database
func (db *Db) DeleteRedelegation(redelegation types.Redelegation) error {
	srcVal, err := db.GetValidator(redelegation.SrcValidator)
	if err != nil {
		return err
	}

	dstVal, err := db.GetValidator(redelegation.DstValidator)
	if err != nil {
		return err
	}

	stmt := `
DELETE FROM redelegation 
WHERE delegator_address = $1 
  AND src_validator_address = $2 
  AND dst_validator_address = $3 
  AND completion_time = $4`
	_, err = db.Sql.Exec(stmt,
		redelegation.DelegatorAddress, srcVal.GetConsAddr(), dstVal.GetConsAddr(), redelegation.CompletionTime,
	)
	return err
}

// DeleteCompletedRedelegations deletes all the redelegations that have completed
// on or before the given timestamp
func (db *Db) DeleteCompletedRedelegations(timestamp time.Time) ([]types.Redelegation, error) {
	stmt := `DELETE FROM redelegation WHERE completion_time <= $1 RETURNING *`

	var rows []dbtypes.RedelegationRow
	err := db.Sqlx.Select(&rows, stmt, timestamp)
	if err != nil {
		return nil, err
	}

	var redelegations = make([]types.Redelegation, len(rows))
	for index, row := range rows {
		srcAddr, err := db.GetValidatorOperatorAddress(row.SrcValidatorAddress)
		if err != nil {
			return nil, err
		}

		dstAddr, err := db.GetValidatorOperatorAddress(row.DstValidatorAddress)
		if err != nil {
			return nil, err
		}

		redelegations[index] = types.NewRedelegation(
			row.DelegatorAddress,
			srcAddr.String(),
			dstAddr.String(),
			row.Amount.ToCoin(),
			row.CompletionTime,
			row.Height,
		)
	}

	return redelegations, err
}

// --------------------------------------------------------------------------------------------------------------------

// SaveUnbondingDelegations saves the given unbonding delegations into the database.
// It assumes that all the validators as well as all the delegators addresses are
// already present inside the proper tables of the database.
// To store the validators data call SaveValidatorData(s).
// To store the account data call SaveAccount.
func (db *Db) SaveUnbondingDelegations(delegations []types.UnbondingDelegation) error {
	// If the delegations are empty just return
	if len(delegations) == 0 {
		return nil
	}

	err := db.storeUpToDateUnbondingDelegations(delegations)
	if err != nil {
		return fmt.Errorf("error while storing up-to-date undonding delegations: %s", err)
	}

	return nil
}

// storeUpToDateUnbondingDelegations allows to store the given unbonding delegations as the most up-to-date ones
func (db *Db) storeUpToDateUnbondingDelegations(delegations []types.UnbondingDelegation) error {
	accQry := `
INSERT INTO account (address) VALUES `
	var accParams []interface{}

	udQry := `
INSERT INTO unbonding_delegation (validator_address, delegator_address, amount, completion_timestamp, height)
VALUES `
	var udParams []interface{}

	for i, delegation := range delegations {
		ai := i * 1
		accQry += fmt.Sprintf("($%d),", ai+1)
		accParams = append(accParams, delegation.DelegatorAddress)

		validator, err := db.GetValidator(delegation.ValidatorOperAddr)
		if err != nil {
			return err
		}

		coin := dbtypes.NewDbCoin(delegation.Amount)
		amount, err := coin.Value()
		if err != nil {
			return err
		}

		udi := i * 5
		udQry += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d),", udi+1, udi+2, udi+3, udi+4, udi+5)
		udParams = append(udParams,
			validator.GetConsAddr(), delegation.DelegatorAddress, amount, delegation.CompletionTimestamp, delegation.Height)
	}

	// Insert the delegators
	accQry = accQry[:len(accQry)-1] // Remove the trailing ","
	accQry += " ON CONFLICT DO NOTHING"
	_, err := db.Sql.Exec(accQry, accParams...)
	if err != nil {
		return err
	}

	// Insert the current unbonding delegations
	udQry = udQry[:len(udQry)-1] // Remove the trailing ","
	udQry += `
ON CONFLICT ON CONSTRAINT unbonding_delegation_validator_delegator_unique 
DO UPDATE SET height = excluded.height
WHERE unbonding_delegation.height <= excluded.height`
	_, err = db.Sql.Exec(udQry, udParams...)
	return err
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

// DeleteUnbondingDelegation removes the given unbonding delegation from the database
func (db *Db) DeleteUnbondingDelegation(delegation types.UnbondingDelegation) error {
	val, err := db.GetValidator(delegation.ValidatorOperAddr)
	if err != nil {
		return err
	}

	stmt := `
DELETE FROM unbonding_delegation 
WHERE delegator_address = $1 
  AND validator_address = $2 
  AND completion_timestamp = $3`
	_, err = db.Sql.Exec(stmt,
		delegation.DelegatorAddress, val.GetConsAddr(), delegation.CompletionTimestamp,
	)
	return err
}

// DeleteCompletedUnbondingDelegations deletes all the unbonding delegations that have completed
// on or before the given timestamp
func (db *Db) DeleteCompletedUnbondingDelegations(timestamp time.Time) ([]types.UnbondingDelegation, error) {
	stmt := `DELETE FROM unbonding_delegation WHERE completion_timestamp <= $1 RETURNING *`

	var rows []dbtypes.UnbondingDelegationRow
	err := db.Sqlx.Select(&rows, stmt, timestamp)
	if err != nil {
		return nil, err
	}

	var delegations = make([]types.UnbondingDelegation, len(rows))
	for index, row := range rows {
		operAddr, err := db.GetValidatorOperatorAddress(row.ConsensusAddress)
		if err != nil {
			return nil, err
		}

		delegations[index] = types.NewUnbondingDelegation(
			row.DelegatorAddress,
			operAddr.String(),
			row.Amount.ToCoin(),
			row.CompletionTimestamp,
			row.Height,
		)
	}

	return delegations, nil
}
