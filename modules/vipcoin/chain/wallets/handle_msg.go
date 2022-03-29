/*
 * Copyright 2022 Business Process Technologies. All rights reserved.
 */

package wallets

import (
	"fmt"

	typeswallets "git.ooo.ua/vipcoin/chain/x/wallets/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	juno "github.com/forbole/juno/v2/types"

	"github.com/forbole/juno/v2/types"
)

// HandleMsg implements MessageModule
func (m *Module) HandleMsg(index int, msg sdk.Msg, tx *types.Tx) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch walletMsg := msg.(type) {
	case *typeswallets.MsgCreateWallet:
		return m.handleMsgCreateWallet(tx, index, walletMsg)
	default:
		errMsg := fmt.Sprintf("unrecognized %s message type: %T", typeswallets.ModuleName, walletMsg)
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
	}
}

// handleMsgCreateWallet allows to properly handle a handleMsgCreateWallet
func (m *Module) handleMsgCreateWallet(tx *juno.Tx, index int, msg *typeswallets.MsgCreateWallet) error {
	if err := m.walletsRepo.SaveCreateWallet(msg); err != nil {
		return err
	}

	newWallet := typeswallets.Wallet{
		Address:        msg.Address,
		AccountAddress: msg.AccountAddress,
		Kind:           msg.Kind,
		State:          msg.State,
		Extras:         msg.Extras,
	}

	return m.walletsRepo.SaveWallets(&newWallet)
}
