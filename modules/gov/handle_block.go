package gov

import (
	"fmt"

	juno "github.com/forbole/juno/v2/types"

	tmctypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/rs/zerolog/log"
)

// HandleBlock implements modules.BlockModule
func (m *Module) HandleBlock(
	b *tmctypes.ResultBlock, _ *tmctypes.ResultBlockResults, _ []*juno.Tx, vals *tmctypes.ResultValidators,
) error {
	err := m.updateParamChangeProposals(b.Block.Height)
	if err != nil {
		log.Error().Str("module", "gov").Int64("height", b.Block.Height).
			Err(err).Msg("error while updating params from ParameterChangeProposals")
	}

	err = m.updateProposals(b.Block.Height, vals)
	if err != nil {
		log.Error().Str("module", "gov").Int64("height", b.Block.Height).
			Err(err).Msg("error while updating proposals")
	}
	return nil
}

// updateParamChangeProposals updates the params if a ParamChangeProposal passed
func (m *Module) updateParamChangeProposals(height int64) error {
	// Get the parameter change proposals
	proposals, err := m.db.GetOpenParamChangeProposals()
	if err != nil {
		log.Error().Err(err).Str("module", "gov").Msg("error while getting ParameterChangeProposal in voting period")
	}

	for _, proposal := range proposals {
		err = m.UpdateParamsFromParamChangeProposal(height, proposal.ProposalID, proposal)
		if err != nil {
			return fmt.Errorf("error while updating parameter change proposal: %s", err)
		}
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
