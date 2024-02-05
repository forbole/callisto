package consumer

import (
	"encoding/json"
	"fmt"

	tmtypes "github.com/cometbft/cometbft/types"

	"github.com/forbole/bdjuno/v4/types"

	ccvconsumertypes "github.com/cosmos/interchain-security/v2/x/ccv/consumer/types"
	"github.com/rs/zerolog/log"
)

// HandleGenesis implements modules.Module
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	log.Debug().Str("module", "ccvconsumer").Msg("parsing genesis")

	// Read the genesis state
	var genState ccvconsumertypes.GenesisState
	err := m.cdc.UnmarshalJSON(appState[ccvconsumertypes.ModuleName], &genState)
	if err != nil {
		return fmt.Errorf("error while reading ccv consumer genesis data: %s", err)
	}

	// Save the params
	err = m.db.SaveCcvConsumerParams(types.NewCcvConsumerParams(genState.Params, doc.InitialHeight))
	if err != nil {
		return fmt.Errorf("error while storing genesis ccv consumer params: %s", err)
	}

	var consumerChains []*types.CcvConsumerChain
	consumerChains = append(consumerChains, types.NewCcvConsumerChain(
		genState.ProviderClientId,
		genState.ProviderChannelId,
		genState.ProviderClientState.ChainId,
		genState.ProviderClientState,
		genState.ProviderConsensusState,
		// genState.InitialValSet, // disable until we can handle cometbft/cometbft instead of /tendermint/tendermint/abci/types".ValidatorUpdate
		nil,
		doc.InitialHeight))

	// Save the consumer chain info
	err = m.db.SaveCcvConsumerChains(consumerChains)

	if err != nil {
		return fmt.Errorf("error while storing genesis ccv consumer chain info: %s", err)
	}

	return nil
}
