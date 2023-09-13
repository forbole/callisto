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

	// switch feeexcluderMsg := msg.(type) {
	switch msg.(type) {
	case *types.MsgCreateAddress:
		return nil
	case *types.MsgDeleteAddress:
		return nil
	case *types.MsgUpdateAddress:
		return nil
	case *types.MsgCreateFees:
		return nil
	case *types.MsgDeleteFees:
		return nil
	case *types.MsgUpdateFees:
		return nil
	default:
		return nil
	}
}
