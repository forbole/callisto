package feegrant

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	feegranttypes "github.com/cosmos/cosmos-sdk/x/feegrant"
	juno "github.com/forbole/juno/v2/types"

	"github.com/forbole/bdjuno/v2/types"
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
	allowance, err := msg.GetFeeAllowanceI()
	if err != nil {
		return fmt.Errorf("error while getting fee allowance: %s", err)
	}
	feeGrant, err := feegranttypes.NewGrant(sdk.AccAddress(msg.Granter), sdk.AccAddress(msg.Grantee), allowance)
	if err != nil {
		return fmt.Errorf("error while getting new grant allowance: %s", err)
	}
	feeGrantAllowance := types.NewFeeGrant(feeGrant, tx.Height)
	return m.db.SaveFeeGrantAllowance(feeGrantAllowance)
}

// HandleMsgRevokeAllowance allows to properly handle a MsgRevokeAllowance
func (m *Module) HandleMsgRevokeAllowance(tx *juno.Tx, msg *feegranttypes.MsgRevokeAllowance) error {
	allowanceToDelete := types.NewGrantRemoval(msg.Grantee, msg.Granter, tx.Height)
	return m.db.DeleteFeeGrantAllowance(allowanceToDelete)
}
