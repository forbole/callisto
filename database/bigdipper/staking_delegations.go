package bigdipper

import (
	"fmt"

	dbtypes "github.com/forbole/bdjuno/database/types"

	"github.com/forbole/bdjuno/types"
)

// SaveDelegations stores inside the database the given delegations data.
// It assumes that the validators addresses are already present inside
// the proper database table.
// TIP: To store the validators data call SaveValidators.
func (db *Db) SaveDelegations(delegations []types.Delegation) error {
	if len(delegations) == 0 {
		return nil
	}

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

		validator, err := db.GetValidator(delegation.ValidatorOperAddr)
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
			validator.GetConsAddr(), delegation.DelegatorAddress, value, delegation.Height)
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

// DeleteDelegatorDelegations removes all the delegations associated with the given delegator
func (db *Db) DeleteDelegatorDelegations(delegator string) error {
	stmt := `DELETE FROM delegation WHERE delegator_address = $1`
	_, err := db.Sql.Exec(stmt, delegator)
	return err
}

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
DO UPDATE SET amount = excluded.amount, height = excluded.height
WHERE redelegation.height <= excluded.height`
	_, err = db.Sql.Exec(rdQry, rdParams...)
	return err
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
		redelegation.DelegatorAddress, srcVal.GetConsAddr(), dstVal.GetOperator(), redelegation.CompletionTime,
	)
	return err
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
DO UPDATE SET amount = excluded.amount, height = excluded.height
WHERE unbonding_delegation.height <= excluded.height`
	_, err = db.Sql.Exec(udQry, udParams...)
	return err
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
