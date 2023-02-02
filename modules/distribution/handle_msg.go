package distribution

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	juno "github.com/forbole/juno/v4/types"

	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
)

// HandleMsg implements modules.MessageModule
func (m *Module) HandleMsg(_ int, msg sdk.Msg, tx *juno.Tx) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	if _, ok := msg.(*distrtypes.MsgFundCommunityPool); ok {
		return m.updateCommunityPool(tx.Height)
	}
	return nil
}
