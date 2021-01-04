package database

import (
	"fmt"

	dbtypes "github.com/forbole/bdjuno/database/types"
	"github.com/forbole/bdjuno/x/staking/types"
)

// SaveHistoricalRedelegation saves the given delegation inside the database.
// It assumes that the validator address is already present inside
// the proper database table.
// TIP: To store the validator data call SaveValidatorData.
func (db *BigDipperDb) SaveHistoricalDelegation(delegation types.Delegation) error {
	accStmt := `INSERT INTO account (address) VALUES ($1) ON CONFLICT DO NOTHING`
	_, err := db.Sql.Exec(accStmt, delegation.ValidatorAddress)
	if err != nil {
		return err
	}

	validator, err := db.GetValidator(delegation.ValidatorAddress)
	if err != nil {
		return err
	}

	coin := dbtypes.NewDbCoin(delegation.Amount)
	value, err := coin.Value()
	if err != nil {
		return err
	}

	stmt := `
INSERT INTO delegation_history (validator_address, delegator_address, amount, shares, height) 
VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
	_, err = db.Sql.Exec(stmt,
		validator.GetConsAddr(), delegation.DelegatorAddress, value, delegation.Shares, delegation.Height,
	)
	return err
}

// SaveCurrentDelegations stores inside the database the given delegations data.
// It assumes that the validators addresses are already present inside
// the proper database table.
// TIP: To store the validators data call SaveValidators.
func (db *BigDipperDb) SaveCurrentDelegations(delegations []types.Delegation) error {
	if len(delegations) == 0 {
		return nil
	}

	accQry := `
INSERT INTO account (address) VALUES `
	var accParams []interface{}

	delQry := `
INSERT INTO delegation (validator_address, delegator_address, amount, shares) VALUES `
	var delParams []interface{}

	delHistQry := `
INSERT INTO delegation_history (validator_address, delegator_address, amount, shares, height)
VALUES `
	var delHistParams []interface{}

	for i, delegation := range delegations {
		ai := i * 1
		accQry += fmt.Sprintf("($%d),", ai+1)
		accParams = append(accParams, delegation.DelegatorAddress)

		validator, err := db.GetValidator(delegation.ValidatorAddress)
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
			validator.GetConsAddr(), delegation.DelegatorAddress, value, delegation.Shares)

		// Historical delegation query
		dhi := i * 5
		delHistQry += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d),", dhi+1, dhi+2, dhi+3, dhi+4, dhi+5)
		delHistParams = append(delHistParams,
			validator.GetConsAddr(), delegation.DelegatorAddress, value, delegation.Shares, delegation.Height)
	}

	// Insert the accounts
	accQry = accQry[:len(accQry)-1] // Remove the trailing ","
	accQry += " ON CONFLICT DO NOTHING"
	_, err := db.Sql.Exec(accQry, accParams...)
	if err != nil {
		return err
	}

	// Remove the current delegations
	_, err = db.Sql.Exec(`DELETE  FROM delegation WHERE TRUE`)
	if err != nil {
		return err
	}

	// Insert the delegations
	delQry = delQry[:len(delQry)-1] // Remove the trailing ","
	delQry += " ON CONFLICT DO NOTHING"
	_, err = db.Sql.Exec(delQry, delParams...)
	if err != nil {
		return err
	}

	// Insert the delegations historical data
	delHistQry = delHistQry[:len(delHistQry)-1] // Remove the trailing ","
	delHistQry += " ON CONFLICT DO NOTHING"
	_, err = db.Sql.Exec(delHistQry, delHistParams...)
	return err
}

// ________________________________________________

// SaveHistoricalUnbondingDelegation saves the given unbonding delegation inside the database.
// It assumes that the validator address is already present inside the proper
// database table.
// TIP: To store the validators data call SaveValidatorData.
func (db *BigDipperDb) SaveHistoricalUnbondingDelegation(delegation types.UnbondingDelegation) error {
	accStmt := `INSERT INTO account(address) VALUES ($1) ON CONFLICT DO NOTHING`
	_, err := db.Sql.Exec(accStmt, delegation.DelegatorAddress)
	if err != nil {
		return err
	}

	validator, err := db.GetValidator(delegation.ValidatorAddress)
	if err != nil {
		return err
	}

	stmt := `
INSERT INTO unbonding_delegation_history (validator_address, delegator_address, amount, completion_timestamp, height) 
VALUES ($1, $2, $3, $4, $5)`

	coin := dbtypes.NewDbCoin(delegation.Amount)
	value, err := coin.Value()
	if err != nil {
		return err
	}

	_, err = db.Sql.Exec(stmt,
		validator.GetConsAddr(), delegation.DelegatorAddress, value,
		delegation.CompletionTimestamp, delegation.Height,
	)
	return err
}

// SaveCurrentUnbondingDelegations saves the given unbonding delegations into the database.
// It assumes that all the validators as well as all the delegators addresses are
// already present inside the proper tables of the database.
// To store the validators data call SaveValidatorData(s).
// To store the account data call SaveAccount.
func (db *BigDipperDb) SaveCurrentUnbondingDelegations(delegations []types.UnbondingDelegation) error {
	// If the delegations are empty just return
	if len(delegations) == 0 {
		return nil
	}

	accQry := `
INSERT INTO account (address) VALUES `
	var accParams []interface{}

	udQry := `
INSERT INTO unbonding_delegation (validator_address, delegator_address, amount, completion_timestamp)
VALUES `
	var udParams []interface{}

	udHistQry := `
INSERT INTO unbonding_delegation_history (validator_address, delegator_address, amount, completion_timestamp, height)
VALUES `
	var udHistParams []interface{}

	for i, delegation := range delegations {
		ai := i * 1
		accQry += fmt.Sprintf("($%d),", ai+1)
		accParams = append(accParams, delegation.DelegatorAddress)

		validator, err := db.GetValidator(delegation.ValidatorAddress)
		if err != nil {
			return err
		}

		coin := dbtypes.NewDbCoin(delegation.Amount)
		amount, err := coin.Value()
		if err != nil {
			return err
		}

		uhi := i * 5
		udHistQry += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d),", uhi+1, uhi+2, uhi+3, uhi+4, uhi+5)
		udHistParams = append(udHistParams,
			validator.GetConsAddr(), delegation.DelegatorAddress, amount,
			delegation.CompletionTimestamp, delegation.Height,
		)

		udi := i * 4
		udQry += fmt.Sprintf("($%d,$%d,$%d,$%d),", udi+1, udi+2, udi+3, udi+4)
		udParams = append(udParams,
			validator.GetConsAddr(), delegation.DelegatorAddress, amount, delegation.CompletionTimestamp)
	}

	// Insert the delegators
	accQry = accQry[:len(accQry)-1] // Remove the trailing ","
	accQry += " ON CONFLICT DO NOTHING"
	_, err := db.Sql.Exec(accQry, accParams...)
	if err != nil {
		return err
	}

	// Delete the current unbonding delegations
	_, err = db.Sql.Exec(`DELETE FROM unbonding_delegation WHERE TRUE`)
	if err != nil {
		return err
	}

	// Insert the current unbonding delegations
	udQry = udQry[:len(udQry)-1] // Remove the trailing ","
	udQry += " ON CONFLICT DO NOTHING"
	_, err = db.Sql.Exec(udQry, udParams...)
	if err != nil {
		return err
	}

	// Insert the historical unbonding delegations
	udHistQry = udHistQry[:len(udHistQry)-1] // Remove the trailing ","
	udHistQry += " ON CONFLICT DO NOTHING"
	_, err = db.Sql.Exec(udHistQry, udHistParams...)
	return err
}

// ________________________________________________

// SaveHistoricalRedelegation saves the given re-delegation inside the database.
// It assumes that the validator info are already present inside the
// proper tables of the database.
// To store the validators data call SaveValidatorData(s).
func (db *BigDipperDb) SaveHistoricalRedelegation(redelegation types.Redelegation) error {
	accStmt := `INSERT INTO account (address) VALUES ($1) ON CONFLICT DO NOTHING`
	_, err := db.Sql.Exec(accStmt, redelegation.DelegatorAddress)
	if err != nil {
		return err
	}

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

	// Insert the data
	stmt := `
INSERT INTO redelegation_history 
    (delegator_address, src_validator_address, dst_validator_address, amount, completion_time, height) 
VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT DO NOTHING`

	_, err = db.Sql.Exec(stmt,
		redelegation.DelegatorAddress, srcVal.GetConsAddr(), dstVal.GetConsAddr(),
		amountValue, redelegation.CompletionTime, redelegation.CreationHeight,
	)
	return err
}

// SaveCurrentRedelegations saves the given redelegations inside the database.
// It assumes that all the validators as well as all the delegators addresses are
// already present inside the proper tables of the database.
// To store the validators data call SaveValidatorData(s).
// To store the account data call SaveAccount.
func (db *BigDipperDb) SaveCurrentRedelegations(redelegations []types.Redelegation) error {
	if len(redelegations) == 0 {
		return nil
	}

	accQry := `
INSERT INTO account (address) VALUES `
	var accParams []interface{}

	rdQry := `
INSERT INTO redelegation (delegator_address, src_validator_address, dst_validator_address, amount, completion_time) 
VALUES `
	var rdParams []interface{}

	rdHisQry := `
INSERT INTO redelegation_history 
	(delegator_address, src_validator_address, dst_validator_address, amount, completion_time, height) 
VALUES `
	var rdHisParams []interface{}

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

		rdi := i * 5
		rdQry += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d),", rdi+1, rdi+2, rdi+3, rdi+4, rdi+5)
		rdParams = append(rdParams,
			redelegation.DelegatorAddress,
			srcVal.GetConsAddr(), dstVal.GetConsAddr(), amountValue, redelegation.CompletionTime)

		rdHi := i * 6
		rdHisQry += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d),",
			rdHi+1, rdHi+2, rdHi+3, rdHi+4, rdHi+5, rdHi+6)
		rdHisParams = append(rdHisParams,
			redelegation.DelegatorAddress, srcVal.GetConsAddr(), dstVal.GetConsAddr(),
			amountValue, redelegation.CompletionTime, redelegation.CreationHeight)
	}

	// Insert the delegators
	accQry = accQry[:len(accQry)-1] // Remove the trailing ","
	accQry += " ON CONFLICT DO NOTHING"
	_, err := db.Sql.Exec(accQry, accParams...)
	if err != nil {
		return err
	}

	// Delete the current redelegations
	_, err = db.Sql.Exec(`DELETE FROM redelegation WHERE TRUE`)
	if err != nil {
		return err
	}

	// Insert the redelegations
	rdQry = rdQry[:len(rdQry)-1] // Remove the trailing ","
	rdQry += " ON CONFLICT DO NOTHING"
	_, err = db.Sql.Exec(rdQry, rdParams...)
	if err != nil {
		return err
	}

	// Insert the historical redelegations
	rdHisQry = rdHisQry[:len(rdHisQry)-1] // Remove the trailing ","
	rdHisQry += " ON CONFLICT DO NOTHING"
	_, err = db.Sql.Exec(rdHisQry, rdHisParams...)
	return err
}
