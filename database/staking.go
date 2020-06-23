package database

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	dbtypes "github.com/forbole/bdjuno/database/types"
	"github.com/forbole/bdjuno/x/staking/types"
)

// SaveStakingPool allows to save for the given height the given stakingtypes pool
func (db BigDipperDb) SaveStakingPool(pool stakingtypes.Pool, height int64, timestamp time.Time) error {
	statement := `INSERT INTO staking_pool (timestamp, height, bonded_tokens, not_bonded_tokens) 
				  VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	_, err := db.Sql.Exec(statement,
		timestamp, height, pool.BondedTokens.Int64(), pool.NotBondedTokens.Int64())
	return err
}

// _________________________________________________________

// SaveValidator saves properly the information about the given validator
func (db BigDipperDb) SaveValidatorData(validator types.Validator) error {
	stmt := `INSERT INTO validator (consensus_address, consensus_pubkey) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err := db.Sql.Exec(stmt,
		validator.GetConsAddr().String(),
		sdk.MustBech32ifyPubKey(sdk.Bech32PubKeyTypeConsPub, validator.GetConsPubKey()),
	)
	if err != nil {
		return err
	}

	stmt = `INSERT INTO validator_info (consensus_address, operator_address) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err = db.Sql.Exec(stmt,
		validator.GetConsAddr().String(), validator.GetOperator().String())
	return err
}

// GetValidatorData returns the validator having the given validator address.
// If no validator for such address can be found, an error is returned instead.
func (db BigDipperDb) GetValidatorData(valAddress sdk.ValAddress) (types.Validator, error) {
	var result []dbtypes.ValidatorData
	stmt := `SELECT validator.consensus_address, validator.consensus_pubkey, validator_info.operator_address 
			 FROM validator INNER JOIN validator_info 
    		 ON validator.consensus_address=validator_info.consensus_address 
			 WHERE validator_info.operator_address = $1`

	if err := db.Sqlx.Select(&result, stmt, valAddress.String()); err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no validator with validator address %s could be found", valAddress.String())
	}

	return result[0], nil
}

// SaveValidatorsData allows the bulk saving of a list of validators
func (db BigDipperDb) SaveValidatorsData(validators []types.Validator) error {
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

// GetValidatorsData returns the data of all the validators that are currently stored inside the database.
func (db BigDipperDb) GetValidatorsData() ([]types.Validator, error) {
	sqlStmt := `SELECT DISTINCT ON (validator.consensus_address)
				validator.consensus_address, validator.consensus_pubkey, validator_info.operator_address 
				FROM validator 
				INNER JOIN validator_info 
				ON validator.consensus_address = validator_info.consensus_address
				ORDER BY validator.consensus_address`

	var rows []dbtypes.ValidatorData
	err := db.Sqlx.Select(&rows, sqlStmt)
	if err != nil {
		return nil, err
	}

	validators := make([]types.Validator, len(rows))
	for index, row := range rows {
		validators[index] = row
	}

	return validators, nil
}

// _________________________________________________________

// SaveValidatorUptime stores into the database the given validator uptime information.
// It assumes that for each uptime information provided, the associated validator data
// have already been saved inside the database properly.
func (db BigDipperDb) SaveValidatorUptime(uptime types.ValidatorUptime) error {
	statement := `INSERT INTO validator_uptime (height, validator_address, signed_blocks_window, missed_blocks_counter)
				  VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	_, err := db.Sql.Exec(statement,
		uptime.Height, uptime.ValidatorAddress.String(), uptime.SignedBlocksWindow, uptime.MissedBlocksCounter)
	return err
}

// _________________________________________________________

// SaveRedelegation saves the given delegation inside the database.
// It assumes that the delegator address is already present inside the
// proper database table.
// TIP: To store the validator data call SaveValidator.
func (db BigDipperDb) SaveDelegation(delegation types.Delegation) error {
	accStmt := `INSERT INTO account (address) VALUES ($1) ON CONFLICT DO NOTHING`
	_, err := db.Sql.Exec(accStmt, delegation.ValidatorAddress.String())
	if err != nil {
		return err
	}

	validator, err := db.GetValidatorData(delegation.ValidatorAddress)
	if err != nil {
		return err
	}

	coin := dbtypes.NewDbCoin(delegation.Amount)
	value, err := coin.Value()
	if err != nil {
		return err
	}

	stmt := `INSERT INTO validator_delegation (consensus_address, delegator_address, amount, height, timestamp) 
			 VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
	_, err = db.Sql.Exec(stmt,
		validator.GetConsAddr().String(), delegation.DelegatorAddress.String(), value,
		delegation.Height, delegation.Timestamp,
	)
	return err
}

// SaveDelegations stores inside the database the given delegations data.
// It assumes that the validators addresses are already present inside
// the proper database table.
// TIP: To store the validators data call SaveValidators.
func (db BigDipperDb) SaveDelegations(delegations []types.Delegation) error {
	accountsQuery := `INSERT INTO account (address) VALUES `
	var accountsParams []interface{}

	delegationsQuery := `INSERT INTO validator_delegation 
						 (consensus_address, delegator_address, amount, height, timestamp) VALUES `
	var delegationsParams []interface{}

	for i, delegation := range delegations {
		a1 := i * 1
		d1 := i * 5 // Starting position for the delegations query

		accountsQuery += fmt.Sprintf("($%d),", a1+1)
		accountsParams = append(accountsParams, delegation.DelegatorAddress.String())

		validator, err := db.GetValidatorData(delegation.ValidatorAddress)
		if err != nil {
			return err
		}

		// Convert the amount
		coin := dbtypes.NewDbCoin(delegation.Amount)
		value, err := coin.Value()
		if err != nil {
			return err
		}

		delegationsQuery += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d),", d1+1, d1+2, d1+3, d1+4, d1+5)
		delegationsParams = append(delegationsParams,
			validator.GetConsAddr().String(), delegation.DelegatorAddress.String(), value,
			delegation.Height, delegation.Timestamp,
		)
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

// SaveRedelegation saves the given unbonding delegation inside the database.
// It assumes that the validator address is already present inside the proper
// database table.
// TIP: To store the validators data call SaveValidator.
func (db BigDipperDb) SaveUnbondingDelegation(delegation types.UnbondingDelegation) error {
	accStmt := `INSERT INTO account(address) VALUES ($1) ON CONFLICT DO NOTHING`
	_, err := db.Sql.Exec(accStmt, delegation.DelegatorAddress.String())
	if err != nil {
		return err
	}

	validator, err := db.GetValidatorData(delegation.ValidatorAddress)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO validator_unbonding_delegation 
    		 (consensus_address, delegator_address, amount, completion_timestamp, height, timestamp) 
    		 VALUES ($1, $2, $3, $4, $5, $6)`

	coin := dbtypes.NewDbCoin(delegation.Amount)
	value, err := coin.Value()
	if err != nil {
		return err
	}

	_, err = db.Sql.Exec(stmt,
		validator.GetConsAddr().String(), delegation.DelegatorAddress.String(), value,
		delegation.CompletionTimestamp, delegation.Height, delegation.Timestamp,
	)
	return err
}

// SaveUnbondingDelegations saves the given unbonding delegations into the database.
// It assumes that all the validators as well as all the delegators addresses are
// already present inside the proper tables of the database.
// To store the validators data call SaveValidator(s).
// To store the account data call SaveAccount.
func (db BigDipperDb) SaveUnbondingDelegations(delegations []types.UnbondingDelegation) error {
	// If the delegations are empty just return
	if len(delegations) == 0 {
		return nil
	}

	accQuery := `INSERT INTO account (address) VALUES `
	var accParams []interface{}

	udQuery := `INSERT INTO validator_unbonding_delegation 
			    (consensus_address, delegator_address, amount, completion_timestamp, height, timestamp) 
			    VALUES `
	var delegationsParams []interface{}

	for i, delegation := range delegations {
		a1 := i * 1
		d1 := i * 6 // Starting position for the delegations query

		accQuery += fmt.Sprintf("($%d),", a1+1)
		accParams = append(accParams, delegation.DelegatorAddress.String())

		validator, err := db.GetValidatorData(delegation.ValidatorAddress)
		if err != nil {
			return err
		}

		coin := dbtypes.NewDbCoin(delegation.Amount)
		amount, err := coin.Value()
		if err != nil {
			return err
		}

		udQuery += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d),", d1+1, d1+2, d1+3, d1+4, d1+5, d1+6)
		delegationsParams = append(delegationsParams,
			validator.GetConsAddr().String(), delegation.DelegatorAddress.String(),
			amount, delegation.CompletionTimestamp, delegation.Height, delegation.Timestamp)
	}

	// Insert the delegators
	accQuery = accQuery[:len(accQuery)-1] // Remove the trailing ","
	accQuery += " ON CONFLICT DO NOTHING"
	_, err := db.Sql.Exec(accQuery, accParams...)
	if err != nil {
		return err
	}

	// Insert the delegations
	udQuery = udQuery[:len(udQuery)-1] // Remove the trailing ","
	udQuery += " ON CONFLICT DO NOTHING"
	_, err = db.Sql.Exec(udQuery, delegationsParams...)
	return err
}

// SaveRedelegation saves the given re-delegation inside the database.
// It assumes that the validator info are already present inside the
// proper tables of the database.
// To store the validators data call SaveValidator(s).
func (db BigDipperDb) SaveRedelegation(redelegation types.Redelegation) error {
	accStmt := `INSERT INTO account (address) VALUES ($1) ON CONFLICT DO NOTHING`
	_, err := db.Sql.Exec(accStmt, redelegation.DelegatorAddress.String())
	if err != nil {
		return err
	}

	// Get the validators info
	srcVal, err := db.GetValidatorData(redelegation.SrcValidator)
	if err != nil {
		return err
	}

	dstVal, err := db.GetValidatorData(redelegation.DstValidator)
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
	stmt := `INSERT INTO validator_redelegation 
    		 (delegator_address, src_validator_address, dst_validator_address, amount, height, completion_time) 
    		 VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT DO NOTHING`

	_, err = db.Sql.Exec(stmt,
		redelegation.DelegatorAddress.String(), srcVal.GetConsAddr().String(), dstVal.GetConsAddr().String(),
		amountValue, redelegation.CreationHeight, redelegation.CompletionTime,
	)
	return err
}

// SaveRedelegations saves the given redelegations inside the database.
// It assumes that all the validators as well as all the delegators addresses are
// already present inside the proper tables of the database.
// To store the validators data call SaveValidator(s).
// To store the account data call SaveAccount.
func (db BigDipperDb) SaveRedelegations(redelegations []types.Redelegation) error {
	accQuery := `INSERT INTO account (address) VALUES `
	var accParams []interface{}

	redelQuery := `INSERT INTO validator_redelegation 
    		 (delegator_address, src_validator_address, dst_validator_address, amount, height, completion_time) 
    		 VALUES `
	var redelParams []interface{}

	for i, redelegation := range redelegations {
		a1 := i * 1 // Starting position for the account query
		r1 := i * 6 // Starting position for the redelegation query

		accQuery += fmt.Sprintf("($%d),", a1+1)
		accParams = append(accParams, redelegation.DelegatorAddress.String())

		// Get the validators info
		srcVal, err := db.GetValidatorData(redelegation.SrcValidator)
		if err != nil {
			return err
		}

		dstVal, err := db.GetValidatorData(redelegation.DstValidator)
		if err != nil {
			return err
		}

		// Convert the amount value
		coin := dbtypes.NewDbCoin(redelegation.Amount)
		amountValue, err := coin.Value()
		if err != nil {
			return err
		}

		redelQuery += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d),", r1+1, r1+2, r1+3, r1+4, r1+5, r1+6)
		redelParams = append(redelParams,
			redelegation.DelegatorAddress.String(), srcVal.GetConsAddr().String(), dstVal.GetConsAddr().String(),
			amountValue, redelegation.CreationHeight, redelegation.CompletionTime)
	}

	// Inser the delegators
	accQuery = accQuery[:len(accQuery)-1] // Remove the trailing ","
	accQuery += " ON CONFLICT DO NOTHING"
	_, err := db.Sql.Exec(accQuery, accParams...)
	if err != nil {
		return err
	}

	// Insert the redelegations
	redelQuery = redelQuery[:len(redelQuery)-1] // Remove the trailing ","
	redelQuery += " ON CONFLICT DO NOTHING"
	_, err = db.Sql.Exec(redelQuery, redelParams...)
	return err
}
