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

//Insert into Validator Commission Database
func (db BigDipperDb) SaveValidatorCommissions(validators []types.ValidatorCommission) error {
	query := `INSERT INTO validator_commission(validator_address,timestamp,commission,min_self_delegation,height) VALUES `
	var param []interface{}
	for i, validator := range validators {
		vi := i * 5
		query += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d),", vi+1, vi+2, vi+3, vi+4, vi+5)
		param = append(param, validator.ValAddress.String(), validator.Timestamp, validator.Commission,
			validator.MinSelfDelegation, validator.Height)
	}
	query = query[:len(query)-1] // Remove trailing ","
	_, err := db.Sql.Exec(query, param...)
	if err != nil {
		return err
	}
	return nil
}

//UpdateValidatorInfo update a single transaction of validator info
//not update very often
func (db BigDipperDb) UpdateValidatorInfo(validator types.Validator) error {
	query := `UPDATE validator_info 
			  SET moniker = $1, identity = $2, website = $3, security_contact = $4, details = $5
			  WHERE consensus_address=$6;`
	_, err := db.Sql.Exec(query, validator.GetDescription().Moniker, validator.GetDescription().Identity, validator.GetDescription().Website, validator.GetDescription().SecurityContact,
		validator.GetDescription().Details, validator.GetConsAddr().String())
	if err != nil {
		return err
	}
	return nil
}

//SaveEditCommission save a single save edit commission operation
func (db BigDipperDb) SaveEditCommission(data types.ValidatorCommission) error {
	statement := `INSERT INTO validator_commission (validator_address,commission,min_self_delegation,height,timestamp)
	 VALUES ($1,$2,$3,$4,$5);`
	_, err := db.Sql.Exec(statement, data.ValAddress.String(), data.Commission, data.MinSelfDelegation, data.Height, data.Timestamp)
	if err != nil {
		return err
	}
	return nil
}

// GetValidatorsData returns the data of all the validators that are currently stored inside the database.
func (db BigDipperDb) GetValidatorsData() ([]dbtypes.ValidatorData, error) {
	sqlStmt := `SELECT DISTINCT ON (validator.consensus_address)
					validator.consensus_address, validator.consensus_pubkey, validator_info.operator_address,
                    validator_info.self_delegate_address, validator_info.moniker, validator_info.identity,
                    validator_info.website, validator_info.security_contact, validator_info.details
				FROM validator 
				INNER JOIN validator_info 
				ON validator.consensus_address = validator_info.consensus_address
				ORDER BY validator.consensus_address`

	var rows []dbtypes.ValidatorData
	err := db.Sqlx.Select(&rows, sqlStmt)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// SaveSingleValidatorData saves properly the information about the given validator
func (db BigDipperDb) SaveSingleValidatorData(validator types.Validator) error {
	return db.SaveValidatorsData([]types.Validator{validator})
}

// GetValidatorData returns the validator having the given validator address.
// If no validator for such address can be found, an error is returned instead.
func (db BigDipperDb) GetValidatorData(valAddress sdk.ValAddress) (types.Validator, error) {
	var result []dbtypes.ValidatorData
	stmt := `SELECT validator.consensus_address, validator.consensus_pubkey, validator_info.operator_address, 
	         	    validator_info.self_delegate_address, validator_info.moniker, validator_info.identity,
       				validator_info.website, validator_info.security_contact, validator_info.details
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

	selfDelegationAccQuery := `INSERT INTO account (address) VALUES `
	var selfDelegationParam []interface{}

	validatorQuery := `INSERT INTO validator (consensus_address, consensus_pubkey) VALUES `
	var validatorParams []interface{}

	validatorInfoQuery := `INSERT INTO validator_info (consensus_address,operator_address,self_delegate_address,moniker,identity,website,security_contact, details) VALUES`
	var validatorInfoParams []interface{}

	for i, validator := range validators {
		v1 := i * 2 // Starting position for validator params
		vi := i * 8 // Starting position for validator info params

		publicKey, err := sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeConsPub, validator.GetConsPubKey())
		if err != nil {
			return err
		}

		selfDelegationAccQuery += fmt.Sprintf("($%d),", i+1)
		selfDelegationParam = append(selfDelegationParam, validator.GetSelfDelegateAddress().String())

		validatorQuery += fmt.Sprintf("($%d,$%d),", v1+1, v1+2)
		validatorParams = append(validatorParams,
			validator.GetConsAddr().String(), publicKey)

		validatorInfoQuery += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d),", vi+1, vi+2, vi+3, vi+4, vi+5, vi+6, vi+7, vi+8)
		validatorInfoParams = append(validatorInfoParams, validator.GetConsAddr().String(), validator.GetOperator().String(), validator.GetSelfDelegateAddress().String(), validator.GetDescription().Moniker,
			validator.GetDescription().Identity, validator.GetDescription().Website, validator.GetDescription().SecurityContact, validator.GetDescription().Details)
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

	validatorInfoQuery = validatorInfoQuery[:len(validatorInfoQuery)-1] // Remove the trailing ","
	validatorInfoQuery += " ON CONFLICT DO NOTHING"
	_, err = db.Sql.Exec(validatorInfoQuery, validatorInfoParams...)
	return err
}

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

//SaveDelegationsShare sve an array of delegation share
func (db BigDipperDb) SaveDelegationsShares(shares []types.DelegationShare) error {
	stmt := `INSERT INTO validator_delegation_shares (operator_address ,delegator_address,shares,height,timestamp) VALUES`
	var delegationShareParam []interface{}
	for i, share := range shares {
		i1 := i * 5
		stmt += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d),", i1+1, i1+2, i1+3, i1+4, i1+5)
		delegationShareParam = append(delegationShareParam, share.ValidatorAddress.String(), share.DelegatorAddress.String(), share.Shares, share.Height, share.Timestamp)
	}
	stmt = stmt[:len(stmt)-1] // Remove the trailing ","
	stmt += " ON CONFLICT DO NOTHING"
	_, err := db.Sql.Exec(stmt, delegationShareParam...)
	if err != nil {
		return err
	}

	return nil
}
