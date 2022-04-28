package wasm

import (
	wasmdtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	juno "github.com/forbole/juno/v3/types"
)

// HandleMsg implements modules.MessageModule
func (m *Module) HandleMsg(_ int, msg sdk.Msg, tx *juno.Tx) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch cosmosMsg := msg.(type) {
	case *wasmdtypes.MsgStoreCode:
		return m.HandleMsgStoreCode(tx, cosmosMsg)

	}

	return nil
}

// HandleMsgStoreCode allows to properly handle a MsgStoreCode
func (m *Module) HandleMsgStoreCode(tx *juno.Tx, msg *wasmdtypes.MsgStoreCode) error {
	return nil
}
