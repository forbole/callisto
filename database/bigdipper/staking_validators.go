package bigdipper

import (
	"fmt"

	dbtypes "github.com/forbole/bdjuno/database/types"

	"github.com/forbole/bdjuno/types"

	bstakingtypes "github.com/forbole/bdjuno/modules/bigdipper/staking/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	bddbtypes "github.com/forbole/bdjuno/database/bigdipper/types"
)

// SaveValidatorData saves properly the information about the given validator.
func (db *Db) SaveValidatorData(validator types.Validator) error {
	return db.SaveValidators([]types.Validator{validator})
}

// SaveValidators allows the bulk saving of a list of validators.
func (db *Db) SaveValidators(validators []types.Validator) error {
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
	validatorInfoQuery += `
ON CONFLICT (consensus_address) DO UPDATE 
	SET consensus_address = excluded.consensus_address,
		operator_address = excluded.operator_address,
		self_delegate_address = excluded.self_delegate_address,
		max_change_rate = excluded.max_change_rate,
		max_rate = excluded.max_rate,
		height = excluded.height
WHERE validator_info.height <= excluded.height`
	_, err = db.Sql.Exec(validatorInfoQuery, validatorInfoParams...)
	return err
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

// GetValidator returns the validator having the given address.
// If no validator for such address can be found, an error is returned instead.
func (db *Db) GetValidator(valAddress string) (types.Validator, error) {
	var result []bddbtypes.ValidatorData
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
func (db *Db) GetValidators() ([]bddbtypes.ValidatorData, error) {
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

	var rows []bddbtypes.ValidatorData
	err := db.Sqlx.Select(&rows, sqlStmt)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// --------------------------------------------------------------------------------------------------------------------

// SaveValidatorDescription save a single validator description.
// It assumes that the delegator address is already present inside the
// proper database table.
// TIP: To store the validator data call SaveValidatorData.
func (db *Db) SaveValidatorDescription(description bstakingtypes.ValidatorDescription) error {
	consAddr, err := db.GetValidatorConsensusAddress(description.OperatorAddress)
	if err != nil {
		return err
	}

	des, err := description.Description.EnsureLength()
	if err != nil {
		return err
	}

	// Update the existing description with this one, if one is already present
	if existing, found := db.getValidatorDescription(consAddr); found {
		des, err = existing.Description.UpdateDescription(des)
		if err != nil {
			return err
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

	_, err = db.Sql.Exec(stmt,
		dbtypes.ToNullString(consAddr.String()),
		dbtypes.ToNullString(des.Moniker),
		dbtypes.ToNullString(des.Identity),
		dbtypes.ToNullString(description.AvatarURL),
		dbtypes.ToNullString(des.Website),
		dbtypes.ToNullString(des.SecurityContact),
		dbtypes.ToNullString(des.Details),
		description.Height,
	)
	return err
}

// getValidatorDescription returns the description of the validator having the given address.
// If no description could be found, returns false instead
func (db *Db) getValidatorDescription(address sdk.ConsAddress) (*bstakingtypes.ValidatorDescription, bool) {
	var result []bddbtypes.ValidatorDescriptionRow
	stmt := `SELECT * FROM validator_description WHERE validator_description.validator_address = $1`

	err := db.Sqlx.Select(&result, stmt, address.String())
	if err != nil {
		return nil, false
	}

	if len(result) == 0 {
		return nil, false
	}

	row := result[0]
	description := bstakingtypes.NewValidatorDescription(
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

// ________________________________________________

// SaveValidatorCommission saves a single validator commission.
// It assumes that the delegator address is already present inside the
// proper database table.
// TIP: To store the validator data call SaveValidatorData.
func (db *Db) SaveValidatorCommission(data bstakingtypes.ValidatorCommission) error {
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
	_, err = db.Sql.Exec(stmt, consAddr.String(), commission, minSelfDelegation, data.Height)
	return err
}

// getValidatorCommission returns the commissions of the validator having the given address.
// If no commissions could be found, returns false instead
func (db *Db) getValidatorCommission(address sdk.ConsAddress) (*bddbtypes.ValidatorCommissionRow, bool) {
	var rows []bddbtypes.ValidatorCommissionRow
	stmt := `SELECT * FROM validator_commission WHERE validator_address = $1`
	err := db.Sqlx.Select(&rows, stmt, address.String())
	if err != nil || len(rows) == 0 {
		return nil, false
	}

	return &rows[0], true
}

// ________________________________________________

// SaveValidatorsVotingPowers saves the given validator voting powers.
// It assumes that the delegator address is already present inside the
// proper database table.
// TIP: To store the validator data call SaveValidatorData.
func (db *Db) SaveValidatorsVotingPowers(entries []bstakingtypes.ValidatorVotingPower) error {
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

	_, err := db.Sql.Exec(stmt, params...)
	return err
}

//---------------------------------------------------

// SaveValidatorsStatuses save validator jail and status in the given height and timestamp
func (db *Db) SaveValidatorsStatuses(statuses []bstakingtypes.ValidatorStatus) error {
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
	_, err := db.Sql.Exec(validatorStmt, valParams...)
	if err != nil {
		return err
	}

	statusStmt = statusStmt[:len(statusStmt)-1]
	statusStmt += `
ON CONFLICT (validator_address) DO UPDATE 
	SET status = excluded.status, 
	    jailed = excluded.jailed,
	    height = excluded.height
WHERE validator_status.height <= excluded.height`
	_, err = db.Sql.Exec(statusStmt, statusParams...)
	return err
}

// saveDoubleSignVote saves the given vote inside the database, returning the row id
func (db *Db) saveDoubleSignVote(vote bstakingtypes.DoubleSignVote) (int64, error) {
	stmt := `
INSERT INTO double_sign_vote 
    (type, height, round, block_id, validator_address, validator_index, signature) 
VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT DO NOTHING RETURNING id`

	var id int64
	err := db.Sql.QueryRow(stmt,
		vote.Type, vote.Height, vote.Round, vote.BlockID, vote.ValidatorAddress, vote.ValidatorIndex, vote.Signature,
	).Scan(&id)
	return id, err
}

// SaveDoubleSignEvidence saves the given double sign evidence inside the proper tables
func (db *Db) SaveDoubleSignEvidence(evidence bstakingtypes.DoubleSignEvidence) error {
	voteA, err := db.saveDoubleSignVote(evidence.VoteA)
	if err != nil {
		return err
	}

	voteB, err := db.saveDoubleSignVote(evidence.VoteB)
	if err != nil {
		return err
	}

	stmt := `
INSERT INTO double_sign_evidence (height, vote_a_id, vote_b_id) 
VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	_, err = db.Sql.Exec(stmt, evidence.Height, voteA, voteB)
	return err
}
