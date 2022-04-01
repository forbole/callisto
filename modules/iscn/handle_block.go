package iscn

import (
	"fmt"

	juno "github.com/forbole/juno/v3/types"

	"github.com/rs/zerolog/log"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/forbole/bdjuno/v2/types"
)

// HandleBlock implements modules.BlockModule
func (m *Module) HandleBlock(block *tmctypes.ResultBlock, _ *tmctypes.ResultBlockResults, _ []*juno.Tx, vals *tmctypes.ResultValidators) error {
	// Update the params
	err := m.updateParams(block.Block.Height)
	if err != nil {
		return err
	}

	return nil
}

// updateParams gets the updated iscn params and stores them inside the database
func (m *Module) updateParams(height int64) error {
	log.Debug().Str("module", "iscn").Int64("height", height).Msg("updating iscn params")

	params, err := m.source.GetParams(height)
	if err != nil {
		return fmt.Errorf("error while getting iscn params: %s", err)
	}

	err = m.db.SaveIscnParams(types.NewIscnParams(params, height))
	if err != nil {
		return fmt.Errorf("error while saving iscn params: %s", err)
	}

	return nil
}
