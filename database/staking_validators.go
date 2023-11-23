package database

import (
	"fmt"

	"github.com/forbole/bdjuno/v4/types"

	dbtypes "github.com/forbole/bdjuno/v4/database/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// SaveValidatorData saves properly the information about the given validator.
func (db *Db) SaveValidatorData(validator types.Validator) error {
	return db.SaveValidatorsData([]types.Validator{validator})
}

// SaveValidatorsData allows the bulk saving of a list of validators.
func (db *Db) SaveValidatorsData(validators []types.Validator) error {
	if len(validators) == 0 {
		return nil
	}

	selfDelegationAccQuery := `
INSERT INTO account (address) VALUES `
	var selfDelegationParam []interface{}

	validatorQuery := `
INSERT INTO validator (consensus_address, consensus_pubkey) VALUES `
	var validatorParams []interface{}

	validatorInfoQuery := `
INSERT INTO validator_info (consensus_address, operator_address, self_delegate_address, max_change_rate, max_rate, height) 
VALUES `
	var validatorInfoParams []interface{}

	for i, validator := range validators {
		vp := i * 2 // Starting position for validator params
		vi := i * 6 // Starting position for validator info params

		selfDelegationAccQuery += fmt.Sprintf("($%d),", i+1)
		selfDelegationParam = append(selfDelegationParam,
			validator.GetSelfDelegateAddress())

		validatorQuery += fmt.Sprintf("($%d,$%d),", vp+1, vp+2)
		validatorParams = append(validatorParams,
			validator.GetConsAddr(), validator.GetConsPubKey())

		validatorInfoQuery += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d),", vi+1, vi+2, vi+3, vi+4, vi+5, vi+6)
		validatorInfoParams = append(validatorInfoParams,
			validator.GetConsAddr(), validator.GetOperator(), validator.GetSelfDelegateAddress(),
			validator.GetMaxChangeRate().String(), validator.GetMaxRate().String(), validator.GetHeight(),
		)
	}

	selfDelegationAccQuery = selfDelegationAccQuery[:len(selfDelegationAccQuery)-1] // Remove trailing ","
	selfDelegationAccQuery += " ON CONFLICT DO NOTHING"
	_, err := db.SQL.Exec(selfDelegationAccQuery, selfDelegationParam...)
	if err != nil {
		return fmt.Errorf("error while storing accounts: %s", err)
	}

	validatorQuery = validatorQuery[:len(validatorQuery)-1] // Remove trailing ","
	validatorQuery += " ON CONFLICT DO NOTHING"
	_, err = db.SQL.Exec(validatorQuery, validatorParams...)
	if err != nil {
		return fmt.Errorf("error while storing valdiators: %s", err)
	}

	validatorInfoQuery = validatorInfoQuery[:len(validatorInfoQuery)-1] // Remove the trailing ","
	validatorInfoQuery += `
ON CONFLICT (consensus_address) DO UPDATE 
	SET consensus_address = excluded.consensus_address,
		operator_address = excluded.operator_address,
		self_delegate_address = excluded.self_delegate_address,
		max_change_rate = excluded.max_change_rate,
		max_rate = excluded.max_rate,
		height = excluded.height
WHERE validator_info.height <= excluded.height`
	_, err = db.SQL.Exec(validatorInfoQuery, validatorInfoParams...)
	if err != nil {
		return fmt.Errorf("error while storing validator infos: %s", err)
	}

	return nil
}

// GetValidatorConsensusAddress returns the consensus address of the validator having the given operator address
func (db *Db) GetValidatorConsensusAddress(address string) (sdk.ConsAddress, error) {
	var result []string
	stmt := `SELECT consensus_address FROM validator_info WHERE operator_address = $1`
	err := db.Sqlx.Select(&result, stmt, address)
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("cannot find the consensus address of validator having operator address %s", address)
	}

	return sdk.ConsAddressFromBech32(result[0])
}

// GetValidatorOperatorAddress returns the operator address of the validator having the given consensus address
func (db *Db) GetValidatorOperatorAddress(consAddr string) (sdk.ValAddress, error) {
	var result []string
	stmt := `SELECT operator_address FROM validator_info WHERE consensus_address = $1`
	err := db.Sqlx.Select(&result, stmt, consAddr)
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("cannot find the operator address of validator having consensus address %s", consAddr)
	}

	return sdk.ValAddressFromBech32(result[0])

}

// GetValidator returns the validator having the given address.
// If no validator for such address can be found, an error is returned instead.
func (db *Db) GetValidator(valAddress string) (types.Validator, error) {
	var result []dbtypes.ValidatorData
	stmt := `
SELECT validator.consensus_address, 
       validator.consensus_pubkey, 
       validator_info.operator_address, 
       validator_info.max_change_rate, 
       validator_info.max_rate,
       validator_info.self_delegate_address
FROM validator INNER JOIN validator_info ON validator.consensus_address=validator_info.consensus_address 
WHERE validator_info.operator_address = $1`

	err := db.Sqlx.Select(&result, stmt, valAddress)
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no validator with validator address %s could be found", valAddress)
	}

	return result[0], nil
}

// GetValidators returns all the validators that are currently stored inside the database.
func (db *Db) GetValidators() ([]types.Validator, error) {
	sqlStmt := `
SELECT DISTINCT ON (validator.consensus_address) 
	validator.consensus_address, 
    validator.consensus_pubkey,
    validator_info.operator_address,
    validator_info.self_delegate_address,                                             
    validator_info.max_rate,
    validator_info.max_change_rate,
    validator_info.height
FROM validator 
INNER JOIN validator_info 
    ON validator.consensus_address = validator_info.consensus_address
ORDER BY validator.consensus_address`

	var rows []dbtypes.ValidatorData
	err := db.Sqlx.Select(&rows, sqlStmt)
	if err != nil {
		return nil, err
	}

	var data = make([]types.Validator, len(rows))
	for index, row := range rows {
		data[index] = row
	}

	return data, nil
}

// GetValidatorBySelfDelegateAddress returns the validator having the given address as the self_delegate_address,
// or an error if such validator cannot be found.
func (db *Db) GetValidatorBySelfDelegateAddress(address string) (types.Validator, error) {
	var result []dbtypes.ValidatorData
	stmt := `
SELECT validator.consensus_address, 
       validator.consensus_pubkey, 
       validator_info.operator_address, 
       validator_info.max_change_rate, 
       validator_info.max_rate,
       validator_info.self_delegate_address
FROM validator INNER JOIN validator_info ON validator.consensus_address=validator_info.consensus_address 
WHERE validator_info.self_delegate_address = $1`

	err := db.Sqlx.Select(&result, stmt, address)
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no validator with self delegate address %s could be found", address)
	}

	return result[0], nil
}

// --------------------------------------------------------------------------------------------------------------------

// SaveValidatorDescription save a single validator description.
// It assumes that the delegator address is already present inside the
// proper database table.
// TIP: To store the validator data call SaveValidatorData.
func (db *Db) SaveValidatorDescription(description types.ValidatorDescription) error {
	consAddr, err := db.GetValidatorConsensusAddress(description.OperatorAddress)
	if err != nil {
		return err
	}

	des, err := description.Description.EnsureLength()
	if err != nil {
		return err
	}

	// Update the existing description with this one, if one is already present
	var avatarURL = description.AvatarURL
	if existing, found := db.getValidatorDescription(consAddr); found {
		des, err = existing.Description.UpdateDescription(des)
		if err != nil {
			return err
		}

		if description.AvatarURL == stakingtypes.DoNotModifyDesc {
			avatarURL = existing.AvatarURL
		}
	}

	// Insert the description
	stmt := `
INSERT INTO validator_description (
	validator_address, moniker, identity, avatar_url, website, security_contact, details, height
)
VALUES($1, $2, $3, $4, $5, $6, $7, $8)
ON CONFLICT (validator_address) DO UPDATE
    SET moniker = excluded.moniker, 
        identity = excluded.identity, 
        avatar_url = excluded.avatar_url,
        website = excluded.website, 
        security_contact = excluded.security_contact, 
        details = excluded.details,
        height = excluded.height
WHERE validator_description.height <= excluded.height`

	_, err = db.SQL.Exec(stmt,
		dbtypes.ToNullString(consAddr.String()),
		dbtypes.ToNullString(des.Moniker),
		dbtypes.ToNullString(des.Identity),
		dbtypes.ToNullString(avatarURL),
		dbtypes.ToNullString(des.Website),
		dbtypes.ToNullString(des.SecurityContact),
		dbtypes.ToNullString(des.Details),
		description.Height,
	)
	if err != nil {
		return fmt.Errorf("error while storing validator description: %s", err)
	}

	return nil
}

// getValidatorDescription returns the description of the validator having the given address.
// If no description could be found, returns false instead
func (db *Db) getValidatorDescription(address sdk.ConsAddress) (*types.ValidatorDescription, bool) {
	var result []dbtypes.ValidatorDescriptionRow
	stmt := `SELECT * FROM validator_description WHERE validator_description.validator_address = $1`

	err := db.Sqlx.Select(&result, stmt, address.String())
	if err != nil {
		return nil, false
	}

	if len(result) == 0 {
		return nil, false
	}

	row := result[0]
	description := types.NewValidatorDescription(
		row.ValAddress,
		stakingtypes.NewDescription(
			dbtypes.ToString(row.Moniker),
			dbtypes.ToString(row.Identity),
			dbtypes.ToString(row.Website),
			dbtypes.ToString(row.SecurityContact),
			dbtypes.ToString(row.Details),
		),
		dbtypes.ToString(row.AvatarURL),
		row.Height,
	)
	return &description, true
}

// --------------------------------------------------------------------------------------------------------------------

// SaveValidatorCommission saves a single validator commission.
// It assumes that the delegator address is already present inside the
// proper database table.
// TIP: To store the validator data call SaveValidatorData.
func (db *Db) SaveValidatorCommission(data types.ValidatorCommission) error {
	if data.MinSelfDelegation == nil && data.Commission == nil {
		// Nothing to update
		return nil
	}

	consAddr, err := db.GetValidatorConsensusAddress(data.ValAddress)
	if err != nil {
		return err
	}

	// Get the existing data, if any
	var commission, minSelfDelegation string
	if existing, found := db.getValidatorCommission(consAddr); found {
		if existing.Commission.Valid {
			commission = existing.Commission.String
		}
		if existing.MinSelfDelegation.Valid {
			minSelfDelegation = existing.MinSelfDelegation.String
		}
	}

	// Replace the existing with the current one
	if data.Commission != nil {
		commission = data.Commission.String()
	}
	if data.MinSelfDelegation != nil {
		minSelfDelegation = data.MinSelfDelegation.String()
	}

	// Update the current value
	stmt := `
INSERT INTO validator_commission (validator_address, commission, min_self_delegation, height) 
VALUES ($1, $2, $3, $4)
ON CONFLICT (validator_address) DO UPDATE 
    SET commission = excluded.commission, 
        min_self_delegation = excluded.min_self_delegation,
        height = excluded.height
WHERE validator_commission.height <= excluded.height`
	_, err = db.SQL.Exec(stmt, consAddr.String(), commission, minSelfDelegation, data.Height)
	if err != nil {
		return fmt.Errorf("error while storing validator commission: %s", err)
	}

	return nil
}

// getValidatorCommission returns the commissions of the validator having the given address.
// If no commissions could be found, returns false instead
func (db *Db) getValidatorCommission(address sdk.ConsAddress) (*dbtypes.ValidatorCommissionRow, bool) {
	var rows []dbtypes.ValidatorCommissionRow
	stmt := `SELECT * FROM validator_commission WHERE validator_address = $1`
	err := db.Sqlx.Select(&rows, stmt, address.String())
	if err != nil || len(rows) == 0 {
		return nil, false
	}

	return &rows[0], true
}

// --------------------------------------------------------------------------------------------------------------------

// SaveValidatorsVotingPowers saves the given validator voting powers.
// It assumes that the delegator address is already present inside the
// proper database table.
// TIP: To store the validator data call SaveValidatorData.
func (db *Db) SaveValidatorsVotingPowers(entries []types.ValidatorVotingPower) error {
	if len(entries) == 0 {
		return nil
	}

	stmt := `INSERT INTO validator_voting_power (validator_address, voting_power, height) VALUES `
	var params []interface{}

	for i, entry := range entries {
		pi := i * 3
		stmt += fmt.Sprintf("($%d,$%d,$%d),", pi+1, pi+2, pi+3)
		params = append(params, entry.ConsensusAddress, entry.VotingPower, entry.Height)
	}

	stmt = stmt[:len(stmt)-1]
	stmt += `
ON CONFLICT (validator_address) DO UPDATE 
	SET voting_power = excluded.voting_power, 
		height = excluded.height
WHERE validator_voting_power.height <= excluded.height`

	_, err := db.SQL.Exec(stmt, params...)
	if err != nil {
		return fmt.Errorf("error while storing validators voting power: %s", err)
	}

	return nil
}

// --------------------------------------------------------------------------------------------------------------------

// SaveValidatorsStatuses save validator jail and status in the given height and timestamp
func (db *Db) SaveValidatorsStatuses(statuses []types.ValidatorStatus) error {
	if len(statuses) == 0 {
		return nil
	}

	validatorStmt := `INSERT INTO validator (consensus_address, consensus_pubkey) VALUES`
	var valParams []interface{}

	statusStmt := `INSERT INTO validator_status (validator_address, status, jailed, height) VALUES `
	var statusParams []interface{}

	for i, status := range statuses {
		vi := i * 2
		validatorStmt += fmt.Sprintf("($%d, $%d),", vi+1, vi+2)
		valParams = append(valParams, status.ConsensusAddress, status.ConsensusPubKey)

		si := i * 4
		statusStmt += fmt.Sprintf("($%d,$%d,$%d,$%d),", si+1, si+2, si+3, si+4)
		statusParams = append(statusParams, status.ConsensusAddress, status.Status, status.Jailed, status.Height)
	}

	validatorStmt = validatorStmt[:len(validatorStmt)-1]
	validatorStmt += "ON CONFLICT DO NOTHING"
	_, err := db.SQL.Exec(validatorStmt, valParams...)
	if err != nil {
		return fmt.Errorf("error while storing validators: %s", err)
	}

	statusStmt = statusStmt[:len(statusStmt)-1]
	statusStmt += `
ON CONFLICT (validator_address) DO UPDATE 
	SET status = excluded.status,
	    jailed = excluded.jailed,
	    height = excluded.height
WHERE validator_status.height <= excluded.height`
	_, err = db.SQL.Exec(statusStmt, statusParams...)
	if err != nil {
		return fmt.Errorf("error while stroring validators statuses: %s", err)
	}

	return nil
}

// saveDoubleSignVote saves the given vote inside the database, returning the row id
func (db *Db) saveDoubleSignVote(vote types.DoubleSignVote) (int64, error) {
	stmt := `
INSERT INTO double_sign_vote 
    (type, height, round, block_id, validator_address, validator_index, signature) 
VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT DO NOTHING RETURNING id`

	var id int64
	err := db.SQL.QueryRow(stmt,
		vote.Type, vote.Height, vote.Round, vote.BlockID, vote.ValidatorAddress, vote.ValidatorIndex, vote.Signature,
	).Scan(&id)
	return id, err
}

// SaveDoubleSignEvidences saves the given double sign evidences inside the database
func (db *Db) SaveDoubleSignEvidences(evidence []types.DoubleSignEvidence) error {
	if len(evidence) == 0 {
		return nil
	}

	stmt := `
INSERT INTO double_sign_evidence (height, vote_a_id, vote_b_id) 
VALUES `

	var doubleSignEvidence []interface{}

	for i, ev := range evidence {
		voteA, err := db.saveDoubleSignVote(ev.VoteA)
		if err != nil {
			return fmt.Errorf("error while storing double sign vote: %s", err)
		}

		voteB, err := db.saveDoubleSignVote(ev.VoteB)
		if err != nil {
			return fmt.Errorf("error while storing double sign vote: %s", err)
		}

		si := i * 3
		stmt += fmt.Sprintf("($%d,$%d,$%d),", si+1, si+2, si+3)
		doubleSignEvidence = append(doubleSignEvidence, ev.Height, voteA, voteB)

	}

	stmt = stmt[:len(stmt)-1] // remove tailing ","
	stmt += " ON CONFLICT DO NOTHING"
	_, err := db.SQL.Exec(stmt, doubleSignEvidence...)
	if err != nil {
		return fmt.Errorf("error while storing double sign evidences: %s", err)
	}

	return nil
}
