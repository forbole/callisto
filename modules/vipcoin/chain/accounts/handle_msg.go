/*
 * Copyright 2022 Business Process Technologies. All rights reserved.
 */

package accounts

import (
	"fmt"

	"git.ooo.ua/vipcoin/chain/x/accounts/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	juno "github.com/forbole/juno/v2/types"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// HandleMsg implements MessageModule
func (m *Module) HandleMsg(index int, msg sdk.Msg, tx *juno.Tx) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch accountMsg := msg.(type) {
	case *types.MsgRegisterUser:
		return m.handleMsgRegisterUser(tx, index, accountMsg)
	default:
		errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, accountMsg)
		fmt.Println(sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg))
		// return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
	}

	return nil
}

// handleMsgRegisterUser allows to properly handle a handleMsgRegisterUser
func (m *Module) handleMsgRegisterUser(tx *juno.Tx, index int, msg *types.MsgRegisterUser) error {
	if err := m.accountRepo.SaveRegisterUser(msg); err != nil {
		return err
	}

	publicKey, err := types.PubKeyFromString(msg.PublicKey)
	if err != nil {
		return types.ErrInvalidPublicKeyField
	}

	publicKeyAny, err := types.PubKeyToAny(publicKey)
	if err != nil {
		return types.ErrInvalidPublicKeyField
	}

	newAcc := types.Account{
		Hash:      msg.Hash,
		Address:   msg.Address,
		Kinds:     []types.AccountKind{types.ACCOUNT_KIND_GENERAL},
		State:     types.ACCOUNT_STATE_ACTIVE,
		PublicKey: publicKeyAny,
		Wallets:   []string{msg.HolderWallet, msg.RefRewardWallet},
	}

	// TODO: Add write wallets.

	return m.accountRepo.SaveAccounts(&newAcc)
}
