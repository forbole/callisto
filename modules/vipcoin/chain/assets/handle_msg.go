package assets

import (
	assetstypes "git.ooo.ua/vipcoin/chain/x/assets/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	juno "github.com/forbole/juno/v2/types"
)

// HandleMsg implements MessageModule
func (m *Module) HandleMsg(index int, msg sdk.Msg, tx *juno.Tx) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch assetMsg := msg.(type) {
	case *assetstypes.MsgAssetCreate:
		return m.handleMsgCreateAsset(tx, index, assetMsg)
	case *assetstypes.MsgAssetManage:
		return m.handleMsgManageAsset(tx, index, assetMsg)
	case *assetstypes.MsgAssetSetExtra:
		return m.handleMsgSetExtraAsset(tx, index, assetMsg)
	default:
		return nil
	}
}
