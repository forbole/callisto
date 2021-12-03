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
		return m.HandleMsgGrantAllowance(tx, cosmosMsg)
	case *feegranttypes.MsgRevokeAllowance:
		return m.HandleMsgRevokeAllowance(tx, cosmosMsg)
	}

	return nil
}

// HandleMsgGrantAllowance allows to properly handle a MsgGrantAllowance
func (m *Module) HandleMsgGrantAllowance(tx *juno.Tx, msg *feegranttypes.MsgGrantAllowance) error {
	allowance := types.NewFeeGrant(msg, tx.Height)
	return m.db.SaveFeeGrantAllowance(allowance)
}

// HandleMsgRevokeAllowance allows to properly handle a MsgRevokeAllowance
func (m *Module) HandleMsgRevokeAllowance(tx *juno.Tx, msg *feegranttypes.MsgRevokeAllowance) error {
	allowanceToDelete := types.NewGrantRemoval(msg.Grantee, msg.Granter, tx.Height)
	return m.db.DeleteFeeGrantAllowance(allowanceToDelete)
}
