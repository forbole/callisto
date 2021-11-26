package feegrant

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	feegranttypes "github.com/cosmos/cosmos-sdk/x/feegrant"
	"github.com/forbole/bdjuno/v2/types"
	juno "github.com/forbole/juno/v2/types"
)

// HandleMsg implements modules.MessageModule
func (m *Module) HandleMsg(index int, msg sdk.Msg, tx *juno.Tx) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch cosmosMsg := msg.(type) {
	case *feegranttypes.MsgGrantAllowance:
		return m.handleMsgGrantAllowance(tx, index, cosmosMsg)

	case *feegranttypes.MsgRevokeAllowance:
		return m.handleMsgRevokeAllowance(tx, index, cosmosMsg)
	}

	return nil
}

// handleMsgGrantAllowance allows to properly handle a MsgGrantAllowance
func (m *Module) handleMsgGrantAllowance(tx *juno.Tx, index int, msg *feegranttypes.MsgGrantAllowance) error {
	grantAllowance := types.NewFeeGrantAllowance(msg.Grantee, msg.Granter, msg.Allowance, tx.Height)

	return m.db.SaveGrantAllowance(types.FeeGrantAllowance(grantAllowance))
}

// handleMsgRevokeAllowance allows to properly handle a MsgRevokeAllowance
func (m *Module) handleMsgRevokeAllowance(tx *juno.Tx, index int, msg *feegranttypes.MsgRevokeAllowance) error {
	return m.db.RevokeGrantAllowance(msg.Grantee, msg.Granter)
}
