package provider

import (
	"encoding/json"
	"fmt"

	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/forbole/bdjuno/v4/types"

	ccvprovidertypes "github.com/cosmos/interchain-security/x/ccv/provider/types"
	"github.com/rs/zerolog/log"
)

// HandleGenesis implements modules.Module
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	log.Debug().Str("module", "ccvprovider").Msg("parsing genesis")

	// Read the genesis state
	var genState ccvprovidertypes.GenesisState
	err := m.cdc.UnmarshalJSON(appState[ccvprovidertypes.ModuleName], &genState)
	if err != nil {
		return fmt.Errorf("error while reading ccv provider genesis data: %s", err)
	}

	// Save the params
	err = m.db.SaveCcvProviderParams(types.NewCcvProviderParams(genState.Params, doc.InitialHeight))
	if err != nil {
		return fmt.Errorf("error while storing genesis ccv provider params: %s", err)
	}

	// Save the provider chain info
	err = m.db.SaveCcvProviderChain(types.NewCcvProviderChain(
		genState.ValsetUpdateId,
		genState.ConsumerStates,
		genState.UnbondingOps,
		genState.MatureUnbondingOps,
		genState.ValsetUpdateIdToHeight,
		genState.ConsumerAdditionProposals,
		genState.ConsumerRemovalProposals,
		genState.ValidatorConsumerPubkeys,
		genState.ValidatorsByConsumerAddr,
		genState.ConsumerAddrsToPrune,
		doc.InitialHeight))

	if err != nil {
		return fmt.Errorf("error while storing genesis ccv provider chain info: %s", err)
	}

	return nil
}
