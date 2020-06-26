package database

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	dbtypes "github.com/forbole/bdjuno/database/types"
	"github.com/forbole/bdjuno/x/staking/types"
	bstaking "github.com/forbole/bdjuno/x/staking/types"
)

// SaveStakingPool allows to save for the given height the given stakingtypes pool
func (db BigDipperDb) SaveStakingPool(pool stakingtypes.Pool, height int64, timestamp time.Time) error {
	statement := `INSERT INTO staking_pool (timestamp, height, bonded_tokens, not_bonded_tokens) 
				  VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	_, err := db.Sql.Exec(statement,
		timestamp, height, pool.BondedTokens.Int64(), pool.NotBondedTokens.Int64())
	return err
}

func (db BigDipperDb) SaveValidatorCommissions(validators []dbtypes.ValidatorCommission) error {
	query := `INSERT INTO validator_commission(validator_address,timestamp,commission,min_self_delegation,height) VALUES`
	var param []interface{}
	for i, validator := range validators {
		vi := i * 5
		query += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d),", vi+1, vi+2, vi+3, vi+4, vi+5)
		param = append(param, validator.ValidatorAddress, validator.Timestamp.UTC, validator.Commission,
			validator.MinSelfDelegation, validator.Height)
	}
	query = query[:len(query)-1] // Remove trailing ","
	query += " ON CONFLICT DO NOTHING"
	_, err := db.Sql.Exec(query, param...)
	if err != nil {
		return err
	}
	return nil
}

//check the new commission is the same as the one before
func (db BigDipperDb) GetCommission(validator sdk.ValAddress) (dbtypes.ValidatorCommission, error) {

	var result []dbtypes.ValidatorCommission
	if found, _ := db.HasValidator(validator.String()); !found {
		return dbtypes.ValidatorCommission{}, nil
	}
	//query the latest entry and see if the validator detail changed
	query := `SELECT commission,min_self_delegation
				FROM validator_commission
				WHERE timestamp = (
					SELECT MAX(timestamp) 
					FROM validator_commission
					WHERE validator_address = $1
				) and validator_address = $2 ;`

	if err := db.Sqlx.Select(&result, query, validator.String(), validator.String()); err != nil {
		return dbtypes.ValidatorCommission{}, err
	}
	if len(result) == 0 {
		return dbtypes.ValidatorCommission{}, fmt.Errorf("no validator with validator address %s could be found", validator.String())
	}
	return result[0], nil
}

func (db BigDipperDb) UpdateValidatorInfo(validator dbtypes.ValidatorInfoRow) error {
	query := `UPDATE validator_info 
				SET moniker=:Moniker,identity=:Identity,website=:Website,securityContact=:SecurityContact, details=:Details)
				 WHERE consAddress=:ConsAddress`
	_, err := db.Sqlx.NamedExec(query, validator)
	if err != nil {
		return err
	}
	return nil
}

func (db BigDipperDb) SaveValidatorInfo(validators []dbtypes.ValidatorInfoRow) error {
	query := `INSERT INTO validator_info(consensus_address,operator_address,moniker,identity,website,securityContact, details) VALUES`
	var param []interface{}
	for i, validator := range validators {
		vi := i * 7
		query += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,$%d),", vi+1, vi+2, vi+3, vi+4, vi+5, vi+6, vi+7)
		param = append(param, validator.ConsAddress, validator.ValAddress, validator.Moniker,
			validator.Identity, validator.Website, validator.SecurityContact, validator.Details)
	}
	query = query[:len(query)-1] // Remove trailing ","
	query += " ON CONFLICT DO NOTHING"
	_, err := db.Sql.Exec(query, param...)
	if err != nil {
		return err
	}
	return nil
}

func (db BigDipperDb) SaveEditCommission(data dbtypes.ValidatorCommission) error {
	statement := `INSERT INTO validator_commission (validator_address,commissions,min_self_delegtion,height,timestamp)
	 VALUES (:ValidatorAddress  ,
		:Timestamp         ,
		:Commission        ,
		:MinSelfDelegation ,
		:Height            );`
	_, err := db.Sqlx.NamedExec(statement, data)
	if err != nil {
		return err
	}
	return nil
}

/*
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
*/
// GetAccounts returns all the accounts that are currently stored inside the database.
func (db BigDipperDb) GetValidators() ([]bstaking.Validator, error) {
	sqlStmt := `SELECT DISTINCT ON (validator.consensus_address)
				validator.consensus_address, validator.consensus_pubkey, validator_info.operator_address 
				FROM validator 
				INNER JOIN validator_info 
				ON validator.consensus_address = validator_info.consensus_address`

	var rows []bstaking.Validator
	err := db.Sqlx.Select(&rows, sqlStmt)
	if err != nil {
		return nil, err
	}
	return rows, nil
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
