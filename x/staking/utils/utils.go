package utils

import (
	"context"

	"github.com/cosmos/cosmos-sdk/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"

	"github.com/forbole/bdjuno/x/utils"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/forbole/bdjuno/x/staking/types"
)

func GetValidatorConsAddr(cdc codec.Marshaler, validator stakingtypes.Validator) (sdk.ConsAddress, error) {
	pubKey, err := GetValidatorConsPubKey(cdc, validator)
	if err != nil {
		return nil, err
	}

	return sdk.ConsAddress(pubKey.Address()), err
}

func GetValidatorConsPubKey(cdc codec.Marshaler, validator stakingtypes.Validator) (cryptotypes.PubKey, error) {
	var pubKey cryptotypes.PubKey
	err := cdc.UnpackAny(validator.ConsensusPubkey, &pubKey)
	return pubKey, err
}

// GetDelegations returns the list of all the delegations that the validator having the given address has
// at the given block height (having the given timestamp)
func GetDelegations(
	validatorAddress string, height int64, stakingClient stakingtypes.QueryClient,
) ([]types.Delegation, error) {
	res, err := stakingClient.ValidatorDelegations(
		context.Background(),
		&stakingtypes.QueryValidatorDelegationsRequest{ValidatorAddr: validatorAddress},
		utils.GetHeightRequestHeader(height),
	)
	if err != nil {
		return nil, err
	}

	delegations := make([]types.Delegation, len(res.DelegationResponses))
	for index, delegation := range res.DelegationResponses {
		delegations[index] = types.NewDelegation(
			delegation.Delegation.DelegatorAddress,
			delegation.Delegation.ValidatorAddress,
			delegation.Balance,
			delegation.Delegation.Shares.String(),
			height,
		)
	}

	return delegations, nil
}

// GetUnbondingDelegations returns the list of all the unbonding delegations that the validator having the
// given address has at the given block height (having the given timestamp).
func GetUnbondingDelegations(
	validatorAddress string, bondDenom string, height int64, stakingClient stakingtypes.QueryClient,
) ([]types.UnbondingDelegation, error) {
	res, err := stakingClient.ValidatorUnbondingDelegations(
		context.Background(),
		&stakingtypes.QueryValidatorUnbondingDelegationsRequest{ValidatorAddr: validatorAddress},
		utils.GetHeightRequestHeader(height),
	)
	if err != nil {
		return nil, err
	}

	var unbondingDelegations []types.UnbondingDelegation
	for _, delegation := range res.UnbondingResponses {
		for _, entry := range delegation.Entries {
			unbondingDelegations = append(unbondingDelegations, types.NewUnbondingDelegation(
				delegation.DelegatorAddress,
				delegation.ValidatorAddress,
				sdk.NewCoin(bondDenom, entry.Balance),
				entry.CompletionTime,
				height,
			))
		}
	}

	return unbondingDelegations, nil
}
