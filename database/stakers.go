package database

import (
	"encoding/json"
	"fmt"

	"github.com/forbole/bdjuno/v4/types"
)

// SaveStakersParams allows to store the given params inside the database
func (db *Db) SaveStakersParams(params *types.StakersParams) error {
	paramsBz, err := json.Marshal(&params.Params)
	if err != nil {
		return fmt.Errorf("error while marshaling stakers params: %s", err)
	}

	stmt := `
INSERT INTO stakers_params (params, height) 
VALUES ($1, $2)
ON CONFLICT (one_row_id) DO UPDATE 
    SET params = excluded.params,
        height = excluded.height
WHERE stakers_params.height <= excluded.height`

	_, err = db.SQL.Exec(stmt, string(paramsBz), params.Height)
	if err != nil {
		return fmt.Errorf("error while storing stakers params: %s", err)
	}

	return nil
}

// SaveProtocolValidators allows to store a list of protocol validators in database
func (db *Db) SaveProtocolValidators(validators []types.ProtocolValidator) error {
	if len(validators) == 0 {
		return nil
	}

	var accounts []types.Account

	protocolValidatorQuery := `
INSERT INTO protocol_validator (address, height) VALUES `
	var protocolValidatorParams []interface{}

	for i, validator := range validators {
		pv := i * 2
		accounts = append(accounts, types.NewAccount(validator.Address))

		protocolValidatorQuery += fmt.Sprintf("($%d,$%d),", pv+1, pv+2)
		protocolValidatorParams = append(protocolValidatorParams,
			string(validator.Address), validator.Height)

	}

	// Store the accounts
	err := db.SaveAccounts(accounts)
	if err != nil {
		return fmt.Errorf("error while storing protocol validators accounts: %s", err)
	}

	protocolValidatorQuery = protocolValidatorQuery[:len(protocolValidatorQuery)-1] // Remove trailing ","
	protocolValidatorQuery += " ON CONFLICT DO NOTHING"
	_, err = db.SQL.Exec(protocolValidatorQuery, protocolValidatorParams...)
	if err != nil {
		return fmt.Errorf("error while storing protocol validators: %s", err)
	}

	return nil
}

// SaveProtocolValidatorsCommission allows to store a list of protocol validators commission in database
func (db *Db) SaveProtocolValidatorsCommission(validators []types.ProtocolValidatorCommission) error {
	if len(validators) == 0 {
		return nil
	}

	var accounts []types.Account

	protocolValidatorCommQuery := `
INSERT INTO protocol_validator_commission (address, commission, pending_commission_change, self_delegation, height) VALUES `
	var protocolValidatorCommParams []interface{}

	for i, validator := range validators {
		pv := i * 5

		accounts = append(accounts, types.NewAccount(validator.Address))

		pendingCommissionChangeBz, err := json.Marshal(&validator.PendingCommissionChange)
		if err != nil {
			return fmt.Errorf("error while marshaling stakers params: %s", err)
		}

		protocolValidatorCommQuery += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d),", pv+1, pv+2, pv+3, pv+4, pv+5)
		protocolValidatorCommParams = append(protocolValidatorCommParams,
			validator.Address,
			validator.Commission.String(),
			string(pendingCommissionChangeBz),
			validator.SelfDelegation,
			validator.Height)

	}

	// Store the accounts
	err := db.SaveAccounts(accounts)
	if err != nil {
		return fmt.Errorf("error while storing protocol validators accounts: %s", err)
	}

	protocolValidatorCommQuery = protocolValidatorCommQuery[:len(protocolValidatorCommQuery)-1] // Remove the trailing ","
	protocolValidatorCommQuery += `
ON CONFLICT (address) DO UPDATE 
	SET commission = excluded.commission,
		pending_commission_change = excluded.pending_commission_change,
		self_delegation = excluded.self_delegation,
		height = excluded.height
WHERE protocol_validator_commission.height <= excluded.height`
	_, err = db.SQL.Exec(protocolValidatorCommQuery, protocolValidatorCommParams...)
	if err != nil {
		return fmt.Errorf("error while storing protocol validators commission: %s", err)
	}

	return nil
}

// SaveProtocolValidatorsDelegation allows to store a list of protocol validators delegation  in database
func (db *Db) SaveProtocolValidatorsDelegation(validators []types.ProtocolValidatorDelegation) error {
	if len(validators) == 0 {
		return nil
	}

	var accounts []types.Account

	protocolValidatorsDelegationQuery := `
INSERT INTO protocol_validator_delegation (address, self_delegation, total_delegation, 
	delegator_count, height) VALUES `
	var protocolValidatorDelegationParams []interface{}

	for i, validator := range validators {
		pv := i * 5

		accounts = append(accounts, types.NewAccount(validator.Address))

		protocolValidatorsDelegationQuery += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d),", pv+1, pv+2, pv+3, pv+4, pv+5)
		protocolValidatorDelegationParams = append(protocolValidatorDelegationParams,
			validator.Address,
			validator.SelfDelegation,
			validator.TotalDelegation,
			validator.DelegatorCount,
			validator.Height)
	}

	// Store the accounts
	err := db.SaveAccounts(accounts)
	if err != nil {
		return fmt.Errorf("error while storing protocol validators accounts: %s", err)
	}

	protocolValidatorsDelegationQuery = protocolValidatorsDelegationQuery[:len(protocolValidatorsDelegationQuery)-1] // Remove the trailing ","
	protocolValidatorsDelegationQuery += `
ON CONFLICT (address) DO UPDATE 
	SET self_delegation = excluded.self_delegation,
		total_delegation = excluded.total_delegation,
		delegator_count = excluded.delegator_count,
		height = excluded.height
WHERE protocol_validator_delegation.height <= excluded.height`
	_, err = db.SQL.Exec(protocolValidatorsDelegationQuery, protocolValidatorDelegationParams...)
	if err != nil {
		return fmt.Errorf("error while storing protocol validators delegation: %s", err)
	}

	return nil
}

// SaveProtocolValidatorsDesc allows to store a list of protocol validators description in database
func (db *Db) SaveProtocolValidatorsDescription(validators []types.ProtocolValidatorDescription) error {
	if len(validators) == 0 {
		return nil
	}

	var accounts []types.Account

	protocolValidatorDescQuery := `
INSERT INTO protocol_validator_description (address, moniker, identity, avatar_url, 
	website, security_contact, details, height) VALUES `
	var protocolValidatorDescParams []interface{}

	for i, validator := range validators {
		pv := i * 8

		accounts = append(accounts, types.NewAccount(validator.Address))

		protocolValidatorDescQuery += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d),", pv+1, pv+2, pv+3, pv+4, pv+5, pv+6, pv+7, pv+8)
		protocolValidatorDescParams = append(protocolValidatorDescParams,
			validator.Address,
			validator.StakerMetadata.Moniker,
			validator.StakerMetadata.Identity,
			validator.AvatarURL,
			validator.StakerMetadata.Website,
			validator.StakerMetadata.SecurityContact,
			validator.StakerMetadata.Details,
			validator.Height)
	}

	// Store the accounts
	err := db.SaveAccounts(accounts)
	if err != nil {
		return fmt.Errorf("error while storing protocol validators accounts: %s", err)
	}

	protocolValidatorDescQuery = protocolValidatorDescQuery[:len(protocolValidatorDescQuery)-1] // Remove the trailing ","
	protocolValidatorDescQuery += `
ON CONFLICT (address) DO UPDATE 
	SET moniker = excluded.moniker,
		identity = excluded.identity,
		avatar_url = excluded.avatar_url,
		website = excluded.website,
		security_contact = excluded.security_contact,
		details = excluded.details,
		height = excluded.height
WHERE protocol_validator_description.height <= excluded.height`
	_, err = db.SQL.Exec(protocolValidatorDescQuery, protocolValidatorDescParams...)
	if err != nil {
		return fmt.Errorf("error while storing protocol validators description: %s", err)
	}

	return nil
}

// SaveProtocolValidatorsPool allows to store a pool info that protocol validator
// is participating in inside database
func (db *Db) SaveProtocolValidatorsPool(validators []types.ProtocolValidatorPool) error {
	if len(validators) == 0 {
		return nil
	}

	var accounts []types.Account

	protocolValidatorPoolQuery := `
INSERT INTO protocol_validator_pool (address, validator_address, balance, pool, height) VALUES `
	var protocolValidatorPoolParams []interface{}

	for i, validator := range validators {
		pv := i * 5

		accounts = append(accounts, types.NewAccount(validator.Address))

		protocolValidatorPoolQuery += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d),", pv+1, pv+2, pv+3, pv+4, pv+5)
		protocolValidatorPoolParams = append(protocolValidatorPoolParams,
			validator.Address,
			validator.ValidatorAddress,
			validator.Balance,
			validator.Pool,
			validator.Height)
	}

	// Store the accounts
	err := db.SaveAccounts(accounts)
	if err != nil {
		return fmt.Errorf("error while storing protocol validators accounts: %s", err)
	}

	protocolValidatorPoolQuery = protocolValidatorPoolQuery[:len(protocolValidatorPoolQuery)-1] // Remove the trailing ","
	protocolValidatorPoolQuery += `
ON CONFLICT ON CONSTRAINT unique_protocol_validator_pool DO UPDATE  
	SET validator_address = excluded.validator_address,
		balance = excluded.balance,
		height = excluded.height
WHERE protocol_validator_pool.height <= excluded.height`
	_, err = db.SQL.Exec(protocolValidatorPoolQuery, protocolValidatorPoolParams...)
	if err != nil {
		return fmt.Errorf("error while storing protocol validators pool info: %s", err)
	}

	return nil
}
