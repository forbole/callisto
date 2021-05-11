package forbolex

import (
	"fmt"

	dbtypes "github.com/forbole/bdjuno/database/types"

	"github.com/forbole/bdjuno/modules/common/staking"
	"github.com/forbole/bdjuno/types"
)

var (
	_ staking.DB = &Db{}
)

// SaveStakingParams implements staking.DB
func (db *Db) SaveStakingParams(params types.StakingParams) error {
	stmt := `
INSERT INTO staking_params (bond_denom) 
VALUES ($1)
ON CONFLICT (one_row_id) DO UPDATE 
    SET bond_denom = excluded.bond_denom`

	_, err := db.Sql.Exec(stmt, params.BondName)
	return err
}

// GetStakingParams implements staking.DB
func (db *Db) GetStakingParams() (*types.StakingParams, error) {
	var rows []dbtypes.StakingParamsRow
	stmt := `SELECT * FROM staking_params LIMIT 1`
	err := db.Sqlx.Select(&rows, stmt)
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, fmt.Errorf("no staking params found")
	}

	return &types.StakingParams{
		BondName: rows[0].BondName,
	}, nil
}

// SaveValidators implements staking.DB
func (db *Db) SaveValidators(validators []types.Validator) error {
	if len(validators) == 0 {
		return nil
	}

	selfDelegationAccQuery := `
INSERT INTO account (address) VALUES `
	var selfDelegationParam []interface{}

	validatorQuery := `
INSERT INTO validator_info (val_oper_addr, val_self_delegate_addr) VALUES `
	var validatorParams []interface{}

	for i, validator := range validators {
		vp := i * 2 // Starting position for validator params

		selfDelegationAccQuery += fmt.Sprintf("($%d),", i+1)
		selfDelegationParam = append(selfDelegationParam,
			validator.GetSelfDelegateAddress())

		validatorQuery += fmt.Sprintf("($%d,$%d),", vp+1, vp+2)
		validatorParams = append(validatorParams,
			validator.GetOperator(), validator.GetSelfDelegateAddress())
	}

	selfDelegationAccQuery = selfDelegationAccQuery[:len(selfDelegationAccQuery)-1] // Remove trailing ","
	selfDelegationAccQuery += " ON CONFLICT DO NOTHING"
	_, err := db.Sql.Exec(selfDelegationAccQuery, selfDelegationParam...)
	if err != nil {
		return err
	}

	validatorQuery = validatorQuery[:len(validatorQuery)-1] // Remove trailing ","
	validatorQuery += " ON CONFLICT DO NOTHING"
	_, err = db.Sql.Exec(validatorQuery, validatorParams...)
	if err != nil {
		return err
	}

	return err
}

// SaveDelegations implements staking.DB
func (db *Db) SaveDelegations(delegations []types.Delegation) error {
	if len(delegations) == 0 {
		return nil
	}

	accQry := `
INSERT INTO account (address) VALUES `
	var accParams []interface{}

	delQry := `
INSERT INTO delegation_history (validator_address, delegator_address, amount, height) VALUES `
	var delParams []interface{}

	for i, delegation := range delegations {
		ai := i * 1
		accQry += fmt.Sprintf("($%d),", ai+1)
		accParams = append(accParams, delegation.DelegatorAddress)

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
			delegation.DelegatorAddress, delegation.DelegatorAddress, value, delegation.Height)
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
DO UPDATE SET amount = excluded.amount`
	_, err = db.Sql.Exec(delQry, delParams...)
	return err
}

// SaveRedelegations implements staking.DB
func (db *Db) SaveRedelegations(redelegations []types.Redelegation) error {
	if len(redelegations) == 0 {
		return nil
	}

	accQry := `
INSERT INTO account (address) VALUES `
	var accParams []interface{}

	rdQry := `
INSERT INTO redelegation_history 
    (delegator_address, src_validator_address, dst_validator_address, amount, completion_time, height) 
VALUES `
	var rdParams []interface{}

	for i, redelegation := range redelegations {
		a1 := i * 1
		accQry += fmt.Sprintf("($%d),", a1+1)
		accParams = append(accParams, redelegation.DelegatorAddress)

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
			redelegation.SrcValidator, redelegation.DstValidator, amountValue,
			redelegation.CompletionTime, redelegation.Height)
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
DO UPDATE SET amount = excluded.amount`
	_, err = db.Sql.Exec(rdQry, rdParams...)
	return err
}

// SaveUnbondingDelegations implements staking.DB
func (db *Db) SaveUnbondingDelegations(delegations []types.UnbondingDelegation) error {
	// If the delegations are empty just return
	if len(delegations) == 0 {
		return nil
	}

	accQry := `
INSERT INTO account (address) VALUES `
	var accParams []interface{}

	udQry := `
INSERT INTO unbonding_delegation_history (validator_address, delegator_address, amount, completion_timestamp, height)
VALUES `
	var udParams []interface{}

	for i, delegation := range delegations {
		ai := i * 1
		accQry += fmt.Sprintf("($%d),", ai+1)
		accParams = append(accParams, delegation.DelegatorAddress)

		coin := dbtypes.NewDbCoin(delegation.Amount)
		amount, err := coin.Value()
		if err != nil {
			return err
		}

		udi := i * 5
		udQry += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d),", udi+1, udi+2, udi+3, udi+4, udi+5)
		udParams = append(udParams,
			delegation.ValidatorOperAddr, delegation.DelegatorAddress, amount,
			delegation.CompletionTimestamp, delegation.Height)
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
DO UPDATE SET amount = excluded.amount`
	_, err = db.Sql.Exec(udQry, udParams...)
	return err
}
