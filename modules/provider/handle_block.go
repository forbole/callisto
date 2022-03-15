package provider

import (
	"fmt"

	juno "github.com/forbole/juno/v2/types"

	"github.com/rs/zerolog/log"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// HandleBlock implements BlockModule
func (m *Module) HandleBlock(
	block *tmctypes.ResultBlock, results *tmctypes.ResultBlockResults, _ []*juno.Tx, _ *tmctypes.ResultValidators,
) error {
	// Update the provider statuses
	err := m.updateProviderStatuses(block.Block.Height)
	if err != nil {
		return fmt.Errorf("error while updating signing info: %s", err)
	}

	return nil
}

// updateProviderStatuses queries the provider statuses and stores them inside the database
func (m *Module) updateProviderStatuses(height int64) error {
	log.Debug().Str("module", "provider").Int64("height", height).Msg("updating provider status")

	providers, err := m.getProviderStatuses(height)
	if err != nil {
		return err
	}

	return m.db.SaveProviders(providers, height)
}
