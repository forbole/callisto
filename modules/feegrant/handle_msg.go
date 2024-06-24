package feegrant

import (
	"fmt"

	feegranttypes "cosmossdk.io/x/feegrant"
	sdk "github.com/cosmos/cosmos-sdk/types"
	juno "github.com/forbole/juno/v6/types"
	"github.com/rs/zerolog/log"

	"github.com/forbole/callisto/v4/types"
	"github.com/forbole/callisto/v4/utils"
)

var msgFilter = map[string]bool{
	"/cosmos.feegrant.v1beta1.MsgGrantAllowance":  true,
	"/cosmos.feegrant.v1beta1.MsgRevokeAllowance": true,
}

// HandleMsgExec implements modules.AuthzMessageModule
func (m *Module) HandleMsgExec(index int, _ int, executedMsg juno.Message, tx *juno.Transaction) error {
	return m.HandleMsg(index, executedMsg, tx)
}

// HandleMsg implements modules.MessageModule
func (m *Module) HandleMsg(_ int, msg juno.Message, tx *juno.Transaction) error {
	if _, ok := msgFilter[msg.GetType()]; !ok {
		return nil
	}

	log.Debug().Str("module", "feegrant").Str("hash", tx.TxHash).Uint64("height", tx.Height).Msg(fmt.Sprintf("handling feegrant message %s", msg.GetType()))

	switch msg.GetType() {
	case "/cosmos.feegrant.v1beta1.MsgGrantAllowance":
		cosmosMsg := utils.UnpackMessage(m.cdc, msg.GetBytes(), &feegranttypes.MsgGrantAllowance{})
		return m.HandleMsgGrantAllowance(tx, cosmosMsg)
	case "/cosmos.feegrant.v1beta1.MsgRevokeAllowance":
		cosmosMsg := utils.UnpackMessage(m.cdc, msg.GetBytes(), &feegranttypes.MsgRevokeAllowance{})
		return m.HandleMsgRevokeAllowance(tx, cosmosMsg)
	}

	return nil
}

// HandleMsgGrantAllowance allows to properly handle a MsgGrantAllowance
func (m *Module) HandleMsgGrantAllowance(tx *juno.Transaction, msg *feegranttypes.MsgGrantAllowance) error {
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
	return m.db.SaveFeeGrantAllowance(types.NewFeeGrant(feeGrant, int64(tx.Height)))
}

// HandleMsgRevokeAllowance allows to properly handle a MsgRevokeAllowance
func (m *Module) HandleMsgRevokeAllowance(tx *juno.Transaction, msg *feegranttypes.MsgRevokeAllowance) error {
	return m.db.DeleteFeeGrantAllowance(types.NewGrantRemoval(msg.Grantee, msg.Granter, int64(tx.Height)))
}
