package gov

import (
	"fmt"

	juno "github.com/forbole/juno/v3/types"

	tmctypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/rs/zerolog/log"
)

// HandleBlock implements modules.BlockModule
func (m *Module) HandleBlock(
	b *tmctypes.ResultBlock, _ *tmctypes.ResultBlockResults, _ []*juno.Tx, vals *tmctypes.ResultValidators,
) error {
	err := m.updateProposals(b.Block.Height, vals)
	if err != nil {
		log.Error().Str("module", "gov").Int64("height", b.Block.Height).
			Err(err).Msg("error while updating proposals")
	}
	return nil
}

// updateProposals updates the proposals
func (m *Module) updateProposals(height int64, blockVals *tmctypes.ResultValidators) error {
	ids, err := m.db.GetOpenProposalsIds()
	if err != nil {
		log.Error().Err(err).Str("module", "gov").Msg("error while getting open ids")
	}

	for _, id := range ids {
		err = m.UpdateProposal(height, blockVals, id)
		if err != nil {
			return fmt.Errorf("error while updating proposal: %s", err)
		}
	}
	return nil
}
