package feegrant

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	feegranttypes "github.com/cosmos/cosmos-sdk/x/feegrant"
	juno "github.com/forbole/juno/v2/types"
)

// HandleMsg implements modules.MessageModule
func (m *Module) HandleMsg(index int, msg sdk.Msg, tx *juno.Tx) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch cosmosMsg := msg.(type) {
	case *feegranttypes.MsgGrantAllowance:
		return m.HandleMsgGrantAllowance(tx, cosmosMsg)
	case *feegranttypes.MsgRevokeAllowance:
		return m.HandleMsgRevokeAllowance(cosmosMsg)
	}

	return nil
}

// handleMsgGrantAllowance allows to properly handle a MsgGrantAllowance
func (m *Module) HandleMsgGrantAllowance(tx *juno.Tx, msg *feegranttypes.MsgGrantAllowance) error {
	return m.db.SaveFeeGrantAllowance(msg, tx.Height)
}

// handleMsgRevokeAllowance allows to properly handle a MsgRevokeAllowance
func (m *Module) HandleMsgRevokeAllowance(msg *feegranttypes.MsgRevokeAllowance) error {
	return m.db.RevokeFeeGrantAllowance(msg.Grantee, msg.Granter)
}
