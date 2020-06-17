package staking

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/desmos-labs/juno/parse/worker"
	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/x/staking/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

func GenesisHandler(
	codec *codec.Codec, genesisDoc *tmtypes.GenesisDoc, appState map[string]json.RawMessage, w worker.Worker,
) error {
	bigDipperDb, ok := w.Db.(database.BigDipperDb)
	if !ok {
		return fmt.Errorf("provided database is not a BigDipper database")
	}

	// Read the genesis state
	var genState staking.GenesisState
	if err := codec.UnmarshalJSON(appState[staking.ModuleName], &genState); err != nil {
		return err
	}

	validators := make([]types.Validator, len(genesisDoc.Validators))
	for _, validator := range genesisDoc.Validators {
		val := types.NewValidator(validator.Address)
		validators = append(validators, validator)
	}

	bigDipperDb.SaveValidators()
}
