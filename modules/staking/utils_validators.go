package staking

import (
	"fmt"
	"time"

	"github.com/forbole/bdjuno/v4/modules/staking/keybase"
	"github.com/forbole/bdjuno/v4/types"

	"github.com/rs/zerolog/log"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// getValidatorConsPubKey returns the consensus public key of the given validator
func (m *Module) getValidatorConsPubKey(validator stakingtypes.Validator) (cryptotypes.PubKey, error) {
	var pubKey cryptotypes.PubKey
	err := m.cdc.UnpackAny(validator.ConsensusPubkey, &pubKey)
	return pubKey, err
}

// getValidatorConsAddr returns the consensus address of the given validator
func (m *Module) getValidatorConsAddr(validator stakingtypes.Validator) (sdk.ConsAddress, error) {
	pubKey, err := m.getValidatorConsPubKey(validator)
	if err != nil {
		return nil, fmt.Errorf("error while getting validator consensus pub key: %s", err)
	}

	return sdk.ConsAddress(pubKey.Address()), err
}

// ---------------------------------------------------------------------------------------------------------------------

// ConvertValidator converts the given staking validator into a BDJuno validator
func (m *Module) convertValidator(height int64, validator stakingtypes.Validator) (types.Validator, error) {
	consAddr, err := m.getValidatorConsAddr(validator)
	if err != nil {
		return nil, fmt.Errorf("error while getting validator consensus address: %s", err)
	}

	consPubKey, err := m.getValidatorConsPubKey(validator)
	if err != nil {
		return nil, fmt.Errorf("error while getting validator consensus pub key: %s", err)
	}

	return types.NewValidator(
		consAddr.String(),
		validator.OperatorAddress,
		consPubKey.String(),
		sdk.AccAddress(validator.GetOperator()).String(),
		&validator.Commission.MaxChangeRate,
		&validator.Commission.MaxRate,
		height,
	), nil
}

// convertValidatorDescription returns a new types.ValidatorDescription object by fetching the avatar URL
// using the Keybase APIs
func (m *Module) convertValidatorDescription(height int64, opAddr string, description stakingtypes.Description) types.ValidatorDescription {
	var avatarURL string

	if description.Identity == stakingtypes.DoNotModifyDesc {
		avatarURL = stakingtypes.DoNotModifyDesc
	} else {
		url, err := keybase.GetAvatarURL(description.Identity)
		if err != nil {
			url = ""
		}
		avatarURL = url
	}

	return types.NewValidatorDescription(opAddr, description, avatarURL, height)
}

// --------------------------------------------------------------------------------------------------------------------

// RefreshAllValidatorInfos refreshes the info of all the validators at the given height
func (m *Module) RefreshAllValidatorInfos(height int64, timestamp time.Time) error {
	// Get all validators
	validators, err := m.source.GetValidatorsWithStatus(height, "")
	if err != nil {
		return fmt.Errorf("error while getting validators: %s", err)
	}

	// Refresh each validator
	for _, validator := range validators {
		err = m.RefreshValidatorInfos(height, timestamp, validator.OperatorAddress)
		if err != nil {
			return fmt.Errorf("error while refreshing validator: %s", err)
		}
	}

	return nil
}

// RefreshValidatorInfos refreshes the info for the validator with the given operator address at the provided height.
// The information include the description as well as the current status
func (m *Module) RefreshValidatorInfos(height int64, timestamp time.Time, valOper string) error {
	// Get the validator information
	stakingValidator, err := m.source.GetValidator(height, valOper)
	if err != nil {
		return err
	}

	validator, err := m.convertValidator(height, stakingValidator)
	if err != nil {
		return fmt.Errorf("error while converting validator: %s", err)
	}

	// Save the validator generic data
	err = m.db.SaveValidatorsData([]types.Validator{validator})
	if err != nil {
		return err
	}

	// Save the description
	err = m.db.SaveValidatorDescription(m.convertValidatorDescription(
		height,
		stakingValidator.OperatorAddress,
		stakingValidator.Description,
	))
	if err != nil {
		return err
	}

	// Save the commission
	err = m.db.SaveValidatorCommission(types.NewValidatorCommission(
		stakingValidator.OperatorAddress,
		&stakingValidator.Commission.Rate,
		&stakingValidator.MinSelfDelegation,
		height,
	))
	if err != nil {
		return err
	}

	return m.updateValidatorsStatusesVotingPowersAndSnapshots(height, timestamp, []stakingtypes.Validator{stakingValidator})
}

// GetValidatorsWithStatus returns the list of all the validators having the given status at the given height
func (m *Module) GetValidatorsWithStatus(height int64, status string) ([]stakingtypes.Validator, []types.Validator, error) {
	validators, err := m.source.GetValidatorsWithStatus(height, status)
	if err != nil {
		return nil, nil, err
	}

	var vals = make([]types.Validator, len(validators))
	for index, val := range validators {
		validator, err := m.convertValidator(height, val)
		if err != nil {
			return nil, nil, fmt.Errorf("error while converting validator: %s", err)
		}

		vals[index] = validator
	}

	return validators, vals, nil
}

// getValidators returns the validators list at the given height
func (m *Module) getValidators(height int64) ([]stakingtypes.Validator, []types.Validator, error) {
	return m.GetValidatorsWithStatus(height, "")
}

// updateValidators updates the list of validators that are present at the given height
func (m *Module) updateValidators(height int64) ([]stakingtypes.Validator, error) {
	log.Debug().Str("module", "staking").Int64("height", height).
		Msg("updating validators")

	vals, validators, err := m.getValidators(height)
	if err != nil {
		return nil, fmt.Errorf("error while getting validator: %s", err)
	}

	err = m.db.SaveValidatorsData(validators)
	if err != nil {
		return nil, err
	}

	return vals, err
}

// --------------------------------------------------------------------------------------------------------------------

func (m *Module) GetValidatorsStatuses(height int64, validators []stakingtypes.Validator) ([]types.ValidatorStatus, error) {
	statuses := make([]types.ValidatorStatus, len(validators))
	for index, validator := range validators {
		consAddr, err := m.getValidatorConsAddr(validator)
		if err != nil {
			return nil, fmt.Errorf("error while getting validator consensus address: %s", err)
		}

		consPubKey, err := m.getValidatorConsPubKey(validator)
		if err != nil {
			return nil, fmt.Errorf("error while getting validator consensus public key: %s", err)
		}

		statuses[index] = types.NewValidatorStatus(
			consAddr.String(),
			consPubKey.String(),
			int(validator.GetStatus()),
			validator.IsJailed(),
			height,
		)
	}

	return statuses, nil
}

// UpdateBondedValidatorsStatusesAndVotingPowers allows to update the status and voting power
// for all the bonded validators.
func (m *Module) UpdateBondedValidatorsStatusesAndVotingPowers() error {
	// Get the latest block height from db
	block, err := m.db.GetLastBlockHeightAndTimestamp()
	if err != nil {
		return fmt.Errorf("error while getting latest block height from db: %s", err)
	}

	// Get all the validators that are still bonded
	validators, _, err := m.GetValidatorsWithStatus(block.Height, stakingtypes.Bonded.String())
	if err != nil {
		return fmt.Errorf("error while getting validators with bonded status: %s", err)
	}

	return m.updateValidatorsStatusesVotingPowersAndSnapshots(block.Height, block.BlockTimestamp, validators)
}

func (m *Module) updateValidatorsStatusesVotingPowersAndSnapshots(height int64, timestamp time.Time, validators []stakingtypes.Validator) error {
	// Update the validators statuses and voting power
	err := m.updateValidatorsStatusesAndVotingPowers(height, validators)
	if err != nil {
		return fmt.Errorf("error while updating validators status and voting power: %s", err)
	}

	// Get all the active on-chain proposals
	ids, err := m.db.GetOpenProposalsIds(timestamp)
	if err != nil {
		return fmt.Errorf("error while getting open proposals ids: %s", err)
	}

	// Update the validators voting power associating them to the given proposals.
	// This will make sure that we can track the voting power of each validator for each proposal
	// once it closes, even if it changes in the future.
	for _, id := range ids {
		err = m.updateProposalValidatorStatusSnapshot(height, id, validators)
		if err != nil {
			return fmt.Errorf("error while updating proposal validator status snapshots: %s", err)
		}
	}

	return nil
}

// updateProposalValidatorStatusSnapshot updates validators snapshot for the proposal having the given id
func (m *Module) updateProposalValidatorStatusSnapshot(height int64, proposalID uint64, validators []stakingtypes.Validator) error {
	snapshots := make([]types.ProposalValidatorStatusSnapshot, len(validators))
	for index, validator := range validators {
		consAddr, err := validator.GetConsAddr()
		if err != nil {
			return err
		}

		var votingPower = validator.Tokens.Int64() / 1_000_000
		if validator.Status != stakingtypes.Bonded || validator.IsJailed() {
			votingPower = 0
		}

		snapshots[index] = types.NewProposalValidatorStatusSnapshot(
			proposalID,
			consAddr.String(),
			votingPower,
			validator.Status,
			validator.Jailed,
			height,
		)
	}

	log.Debug().Str("module", "staking").Msg("refreshing proposal validator statuses snapshots")

	return m.db.SaveProposalValidatorsStatusesSnapshots(snapshots)
}

// updateValidatorsStatusesAndVotingPowers updates validators status
// and validators voting power
func (m *Module) updateValidatorsStatusesAndVotingPowers(height int64, validators []stakingtypes.Validator) error {
	votingPowers := make([]types.ValidatorVotingPower, len(validators))
	statuses := make([]types.ValidatorStatus, len(validators))

	for index, validator := range validators {
		consAddr, err := validator.GetConsAddr()
		if err != nil {
			return err
		}

		if found, _ := m.db.HasValidator(consAddr.String()); !found {
			continue
		}

		consPubKey, err := m.getValidatorConsPubKey(validator)
		if err != nil {
			return err
		}

		var votingPower = validator.Tokens.Int64() / 1_000_000
		if validator.Status != stakingtypes.Bonded || validator.IsJailed() {
			votingPower = 0
		}

		votingPowers[index] = types.NewValidatorVotingPower(
			consAddr.String(),
			votingPower,
			height,
		)

		statuses[index] = types.NewValidatorStatus(
			consAddr.String(),
			consPubKey.String(),
			int(validator.GetStatus()),
			validator.IsJailed(),
			height,
		)
	}

	log.Debug().Str("module", "staking").Msg("refreshing validator voting power")
	// Save validators voting powers in db
	err := m.db.SaveValidatorsVotingPowers(votingPowers)
	if err != nil {
		log.Error().Str("module", "staking").Err(err).Int64("height", height).
			Msg("error while saving validators voting powers")
	}

	log.Debug().Str("module", "staking").Msg("refreshing validator statuses")
	// Save validators statuses in db
	err = m.db.SaveValidatorsStatuses(statuses)
	if err != nil {
		log.Error().Str("module", "staking").Err(err).
			Int64("height", height).
			Msg("error while saving validators statuses")
	}

	return nil
}
