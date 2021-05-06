package common

import (
	"context"

	"github.com/cosmos/cosmos-sdk/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	bstakingtypes "github.com/forbole/bdjuno/modules/bigdipper/staking/types"
	utils2 "github.com/forbole/bdjuno/modules/common/utils"

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

// ConvertValidator converts the given staking validator into a BDJuno validator
func ConvertValidator(cdc codec.Marshaler, validator stakingtypes.Validator, height int64) (bstakingtypes.Validator, error) {
	consAddr, err := GetValidatorConsAddr(cdc, validator)
	if err != nil {
		return nil, err
	}

	consPubKey, err := GetValidatorConsPubKey(cdc, validator)
	if err != nil {
		return nil, err
	}

	return bstakingtypes.NewValidator(
		consAddr.String(),
		validator.OperatorAddress,
		consPubKey.String(),
		sdk.AccAddress(validator.GetOperator()).String(),
		&validator.Commission.MaxChangeRate,
		&validator.Commission.MaxRate,
		height,
	), nil
}

// GetValidators returns the validators list at the given height
func GetValidators(height int64, client stakingtypes.QueryClient) ([]stakingtypes.Validator, error) {
	header := utils2.GetHeightRequestHeader(height)

	var validators []stakingtypes.Validator
	var nextKey []byte
	var stop = false
	for !stop {
		res, err := client.Validators(
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
			return nil, err
		}

		nextKey = res.Pagination.NextKey
		stop = len(res.Pagination.NextKey) == 0
		validators = append(validators, res.Validators...)
	}

	return validators, nil
}
