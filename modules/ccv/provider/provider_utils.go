package provider

import (
	"fmt"

	"github.com/forbole/bdjuno/v4/types"
	"github.com/rs/zerolog/log"
)

// UpdateAllConsumerChains updates consumer chains for the given height
func (m *Module) UpdateAllConsumerChains(height int64) error {
	log.Debug().Str("module", "ccvprovider").Int64("height", height).
		Msg("updating all consumer chains")

	// Get all consumer chains
	consumerChains, err := m.source.GetAllConsumerChains(height)
	if err != nil {
		return fmt.Errorf("error while getting consumer chains: %s", err)
	}

	var consumerChainsList []*types.CcvConsumerChain

	for _, cc := range consumerChains {
		consumerChainsList = append(consumerChainsList, types.NewCcvConsumerChain(
			cc.ClientId,
			"",
			false,
			cc.ChainId,
			nil,
			nil,
			nil,
			height))
	}

	return m.db.SaveCcvConsumerChains(consumerChainsList)
}

// RemoveConsumerChain removes consumer chain from database
func (m *Module) RemoveConsumerChain(height int64, chainID string) error {
	log.Debug().Str("module", "ccvprovider").Int64("height", height).
		Msg("removing consumer chain from database")

	return m.db.DeleteConsumerChainFromDB(height, chainID)
}
