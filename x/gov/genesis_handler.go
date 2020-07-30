package gov

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/desmos-labs/juno/parse/worker"
	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/x/gov/types"
	"github.com/rs/zerolog/log"
	tmtypes "github.com/tendermint/tendermint/types"
)

func GenesisHandler(codec *codec.Codec, genesisDoc *tmtypes.GenesisDoc, appState map[string]json.RawMessage, w worker.Worker) error {
	log.Debug().Str("module", "gov").Msg("parsing genesis")

	bigDipperDb, ok := w.Db.(database.BigDipperDb)
	if !ok {
		return fmt.Errorf("given database instance is not a BigDipperDb")
	}
	// Read the genesis state
	var genState gov.GenesisState
	if err := codec.UnmarshalJSON(appState[gov.ModuleName], &genState); err != nil {
		return err
	}

	if err := saveProposals(genState.Proposals, genesisDoc, bigDipperDb); err != nil {
		return err
	}
	return nil
}

func saveProposals(proposals gov.Proposals, genesisDoc *tmtypes.GenesisDoc, db database.BigDipperDb)error {
	return nil
}
