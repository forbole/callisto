package accounts

import (
	"git.ooo.ua/vipcoin/chain/x/accounts/types"
	juno "github.com/forbole/juno/v2/types"
)

// handleMsgCreateAccount allows to properly handle a handleMsgCreateAccount
func (m *Module) handleMsgCreateAccount(tx *juno.Tx, index int, msg *types.MsgCreateAccount) error {
	if err := m.accountRepo.SaveCreateAccount(msg, tx.TxHash); err != nil {
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
		Hash:       msg.Hash,
		Address:    msg.Address,
		PublicKey:  publicKeyAny,
		Kinds:      msg.Kinds,
		State:      msg.State,
		Extras:     msg.Extras,
		Affiliates: []*types.Affiliate{},
		Wallets:    []string{},
	}

	return m.accountRepo.SaveAccounts(&newAcc)
}
