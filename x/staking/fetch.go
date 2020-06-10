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

func getValidatorDetails(validatorAddress tmtypes.Address, w worker.Worker) (*staking.Validator, error) {
	// Create the validator address from the hex string
	valOperAddr, err := sdk.ValAddressFromHex(hex.EncodeToString(validatorAddress))
	if err != nil {
		return nil, err
	}

	// Get the validator address
	var validator staking.Validator
	endpoint := fmt.Sprintf("staking/validators/%s", valOperAddr.String())
	if err := w.ClientProxy.QueryLCD(endpoint, &validator); err != nil {
		return nil, err
	}

	return &validator, nil
}

// DataFetcher queries the proper data that should be later handled for the staking module.
func DataFetcher(
	block *coretypes.ResultBlock, txs []djuno.Tx, vals *tmctypes.ResultValidators, w worker.Worker,
) (interface{}, error) {

	var validators []*staking.Validator

	for _, val := range vals.Validators {
		// Get the validator details
		validator, err := getValidatorDetails(val.Address, w)
		if err != nil {
			return nil, err
		}
		validators = append(validators, validator)

		// TODO: Get the delegations
	}

	// TODO: Return the proper data
	return nil, nil
}
