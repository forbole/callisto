package pool

import (
	"fmt"

	pooltypes "github.com/KYVENetwork/chain/x/pool/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	juno "github.com/forbole/juno/v4/types"
)

// HandleMsg implements modules.MessageModule
func (m *Module) HandleMsg(index int, msg sdk.Msg, tx *juno.Tx) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch cosmosMsg := msg.(type) {
	case *pooltypes.MsgCreatePool:
		return m.handleUpdatePools(tx, cosmosMsg)
	}

	return nil
}

// handleUpdatePools allows to properly handle a MsgCreatePool
func (m *Module) handleUpdatePools(tx *juno.Tx, msg *pooltypes.MsgCreatePool) error {
	// refresh info for all pools
	err := m.UpdatePools(tx.Height)
	if err != nil {
		return fmt.Errorf("error while updating pools list: %s", err)
	}

	return nil
}
