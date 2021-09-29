package gov

import (
	"fmt"

	juno "github.com/desmos-labs/juno/types"

	tmctypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/types"
)

// HandleBlock implements modules.BlockModule
func (m *Module) HandleBlock(b *tmctypes.ResultBlock, _ []*juno.Tx, vals *tmctypes.ResultValidators) error {
	err := m.updateProposals(b.Block.Height, vals)
	if err != nil {
		log.Error().Str("module", "gov").Int64("height", b.Block.Height).
			Err(err).Msg("error while updating proposals")
	}

	err = m.updateParams(b.Block.Height)
	if err != nil {
		log.Error().Str("module", "gov").Int64("height", b.Block.Height).
			Err(err).Msg("error while updating params")
	}

	return nil
}

// updateParams updates the governance parameters for the given height
func (m *Module) updateParams(height int64) error {
	depositParams, err := m.source.DepositParams(height)
	if err != nil {
		return fmt.Errorf("error while getting gov deposit params: %s", err)
	}

	votingParams, err := m.source.VotingParams(height)
	if err != nil {
		return fmt.Errorf("error while getting gov voting params: %s", err)
	}

	tallyParams, err := m.source.TallyParams(height)
	if err != nil {
		return fmt.Errorf("error while getting gov tally params: %s", err)
	}

	return m.db.SaveGovParams(types.NewGovParams(
		types.NewVotingParams(votingParams),
		types.NewDepositParam(depositParams),
		types.NewTallyParams(tallyParams),
		height,
	))
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
