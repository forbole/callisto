package feegrant

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	feegranttypes "github.com/cosmos/cosmos-sdk/x/feegrant"
	juno "github.com/forbole/juno/v3/types"

	"github.com/forbole/bdjuno/v3/types"
)

// HandleMsg implements modules.MessageModule
func (m *Module) HandleMsg(_ int, msg sdk.Msg, tx *juno.Tx) error {
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
	granter, err := sdk.AccAddressFromBech32(msg.Granter)
	if err != nil {
		return fmt.Errorf("error while parsing granter address: %s", err)
	}
	grantee, err := sdk.AccAddressFromBech32(msg.Grantee)
	if err != nil {
		return fmt.Errorf("error while parsing grantee address: %s", err)
	}
	feeGrant, err := feegranttypes.NewGrant(granter, grantee, allowance)
	if err != nil {
		return fmt.Errorf("error while getting new grant allowance: %s", err)
	}
	return m.db.SaveFeeGrantAllowance(types.NewFeeGrant(feeGrant, tx.Height))
}

// HandleMsgRevokeAllowance allows to properly handle a MsgRevokeAllowance
func (m *Module) HandleMsgRevokeAllowance(tx *juno.Tx, msg *feegranttypes.MsgRevokeAllowance) error {
	return m.db.DeleteFeeGrantAllowance(types.NewGrantRemoval(msg.Grantee, msg.Granter, tx.Height))
}
