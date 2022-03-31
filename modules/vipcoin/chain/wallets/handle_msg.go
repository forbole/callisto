/*
 * Copyright 2022 Business Process Technologies. All rights reserved.
 */

package wallets

import (
	"fmt"

	typeswallets "git.ooo.ua/vipcoin/chain/x/wallets/types"
	"git.ooo.ua/vipcoin/lib/filter"

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
	case *typeswallets.MsgSetWalletState:
		return m.handleMsgSetStates(walletMsg)
	case *typeswallets.MsgCreateWallet:
		return m.handleMsgCreateWallet(tx, index, walletMsg)
	case *typeswallets.MsgSetDefaultWallet:
		return m.handleMsgSetDefaultWallet(walletMsg)
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

// handleMsgSetStates allows to properly handle a MsgSetState
func (m *Module) handleMsgSetStates(msg *typeswallets.MsgSetWalletState) error {
	if err := m.walletsRepo.SaveStates(msg); err != nil {
		return err
	}

	wallets, err := m.walletsRepo.GetWallets(filter.NewFilter().SetArgument(FieldAddress, msg.Address))
	if err != nil {
		return err
	}

	if len(wallets) != 1 {
		return typeswallets.ErrInvalidAddressField
	}

	wallets[0].State = msg.State

	return m.walletsRepo.SaveWallets(wallets...)
}

// handleMsgSetKinds allows to properly handle a handleMsgSetKinds
func (m *Module) handleMsgSetDefaultWallet(msg *typeswallets.MsgSetDefaultWallet) error {
	if err := m.walletsRepo.SaveDefaultWallets(msg); err != nil {
		return err
	}

	wallets, err := m.walletsRepo.GetWallets(filter.NewFilter().SetArgument(FieldAddress, msg.Address))
	if err != nil {
		return err
	}

	if len(wallets) != 1 {
		return typeswallets.ErrInvalidAddressField
	}

	return m.walletsRepo.UpdateWallets(wallets...)
}
