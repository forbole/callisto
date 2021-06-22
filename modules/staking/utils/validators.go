package utils

import (
	"context"

	"github.com/desmos-labs/juno/client"

	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/modules/staking/keybase"
	"github.com/forbole/bdjuno/types"

	"github.com/rs/zerolog/log"

	"github.com/cosmos/cosmos-sdk/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// GetValidatorConsPubKey returns the consensus public key of the given validator
func GetValidatorConsPubKey(cdc codec.Marshaler, validator stakingtypes.Validator) (cryptotypes.PubKey, error) {
	var pubKey cryptotypes.PubKey
	err := cdc.UnpackAny(validator.ConsensusPubkey, &pubKey)
	return pubKey, err
}

// GetValidatorConsAddr returns the consensus address of the given validator
func GetValidatorConsAddr(cdc codec.Marshaler, validator stakingtypes.Validator) (sdk.ConsAddress, error) {
	pubKey, err := GetValidatorConsPubKey(cdc, validator)
	if err != nil {
		return nil, err
	}

	return sdk.ConsAddress(pubKey.Address()), err
}

// ---------------------------------------------------------------------------------------------------------------------

// ConvertValidator converts the given staking validator into a BDJuno validator
func ConvertValidator(
	cdc codec.Marshaler, validator stakingtypes.Validator, height int64,
) (types.Validator, error) {
	consAddr, err := GetValidatorConsAddr(cdc, validator)
	if err != nil {
		return nil, err
	}

	consPubKey, err := GetValidatorConsPubKey(cdc, validator)
	if err != nil {
		return nil, err
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

// ConvertValidatorDescription returns a new types.ValidatorDescription object by fetching the avatar URL
// using the Keybase APIs
func ConvertValidatorDescription(
	opAddr string, description stakingtypes.Description, height int64,
) (types.ValidatorDescription, error) {
	avatarURL, err := keybase.GetAvatarURL(description.Identity)
	if err != nil {
		return types.ValidatorDescription{}, err
	}

	return types.NewValidatorDescription(opAddr, description, avatarURL, height), nil
}

// --------------------------------------------------------------------------------------------------------------------

// GetValidators returns the validators list at the given height
func GetValidators(
	height int64, stakingClient stakingtypes.QueryClient, cdc codec.Marshaler,
) ([]stakingtypes.Validator, []types.Validator, error) {
	header := client.GetHeightRequestHeader(height)

	var validators []stakingtypes.Validator
	var nextKey []byte
	var stop = false
	for !stop {
		res, err := stakingClient.Validators(
			context.Background(),
			&stakingtypes.QueryValidatorsRequest{
				Status: "", // Query all the statues
				Pagination: &query.PageRequest{
					Key:   nextKey,
					Limit: 100, // Query 100 validators at time
				},
			},
			header,
		)
		if err != nil {
			return nil, nil, err
		}

		nextKey = res.Pagination.NextKey
		stop = len(res.Pagination.NextKey) == 0
		validators = append(validators, res.Validators...)
	}

	var vals = make([]types.Validator, len(validators))
	for index, val := range validators {
		validator, err := ConvertValidator(cdc, val, height)
		if err != nil {
			return nil, nil, err
		}

		vals[index] = validator
	}

	return validators, vals, nil
}

// UpdateValidators updates the list of validators that are present at the given height
func UpdateValidators(
	height int64, client stakingtypes.QueryClient, cdc codec.Marshaler, db *database.Db,
) ([]stakingtypes.Validator, error) {
	log.Debug().Str("module", "staking").Int64("height", height).
		Msg("updating validators")

	vals, validators, err := GetValidators(height, client, cdc)
	if err != nil {
		return nil, err
	}

	err = db.SaveValidatorsData(validators)
	if err != nil {
		return nil, err
	}

	return vals, err
}
