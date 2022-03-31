package accounts

import (
	"git.ooo.ua/vipcoin/chain/x/accounts/types"
	juno "github.com/forbole/juno/v2/types"
)

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
