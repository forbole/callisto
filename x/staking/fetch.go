package staking

import (
	"encoding/hex"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/desmos-labs/juno/parse/worker"
	djuno "github.com/desmos-labs/juno/types"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

// getValidatorDetails queries the details of the validator having the given address from the LCD endpoint
func getValidatorDetails(validatorAddress tmtypes.Address, w worker.Worker) (*staking.Validator, error) {
	// Create the validator address from the hex string
	valOperAddr := sdk.ValAddress(validatorAddress.Bytes())

	// Get the validator details
	var validator staking.Validator
	endpoint := fmt.Sprintf("staking/validators/%s", valOperAddr.String())
	if err := w.ClientProxy.QueryLCD(endpoint, &validator); err != nil {
		return nil, err
	}

	return &validator, nil
}

// getValidatorDetails queries the delegations of the validator having the given address from the LCD endpoint
func getValidatorDelegations(validatorAddress tmtypes.Address, w worker.Worker) (*staking.Delegations, error) {
	// Create the validator address from the hex string
	valOperAddr, err := sdk.ValAddressFromHex(hex.EncodeToString(validatorAddress))
	if err != nil {
		return nil, err
	}

	// Get the validator delegations
	var validator staking.Delegations
	endpoint := fmt.Sprintf("staking/validators/%s/delegations", valOperAddr.String())
	if err := w.ClientProxy.QueryLCD(endpoint, &validator); err != nil {
		return nil, err
	}

	return &validator, nil
}

// DataFetcher queries the proper data that should be later handled for the staking module.
func DataFetcher(
	block *coretypes.ResultBlock, txs []djuno.Tx, vals *tmctypes.ResultValidators, w worker.Worker,
) (interface{}, error) {

	var validators []ValidatorInfo

	for _, val := range vals.Validators {
		// Get the validator details
		validator, err := getValidatorDetails(val.Address, w)
		if err != nil {
			return nil, err
		}

		// Get the delegations
		delegations, err := getValidatorDelegations(val.Address, w)
		if err != nil {
			return nil, err
		}

		validatorInfo := NewValidatorInfo(validator, delegations)
		validators = append(validators, validatorInfo)
	}

	return validators, nil
}
