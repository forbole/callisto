package provider

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	juno "github.com/forbole/juno/v2/types"
	providertypes "github.com/ovrclk/akash/x/provider/types/v1beta2"
)

// HandleMsg implements modules.MessageModule
func (m *Module) HandleMsg(_ int, msg sdk.Msg, tx *juno.Tx) error {
	cosmosMsg, ok := msg.(*providertypes.MsgDeleteProvider)
	if ok {
		return m.db.DeleteProvider(cosmosMsg.GetOwner(), tx.Height)
	}

	return nil
}
