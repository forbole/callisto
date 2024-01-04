package feeexcluder

import (
	"git.ooo.ua/vipcoin/ovg-chain/x/feeexcluder/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	juno "github.com/forbole/juno/v5/types"
)

// HandleMsg implements MessageModule
func (m *Module) HandleMsg(index int, msg sdk.Msg, tx *juno.Tx) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch feeExcluderMsg := msg.(type) {
	case *types.MsgCreateAddress:
		return m.handleMsgCreateAddress(tx, index, feeExcluderMsg)
	case *types.MsgUpdateAddress:
		return m.handleMsgUpdateAddress(tx, index, feeExcluderMsg)
	case *types.MsgDeleteAddress:
		return m.handleMsgDeleteAddress(tx, index, feeExcluderMsg)
	case *types.MsgCreateTariffs:
		return m.handleMsgCreateTariffs(tx, index, feeExcluderMsg)
	case *types.MsgUpdateTariffs:
		return m.handleMsgUpdateTariffs(tx, index, feeExcluderMsg)
	case *types.MsgDeleteTariffs:
		return m.handleMsgDeleteTariffs(tx, index, feeExcluderMsg)
	default:
		return nil
	}
}
