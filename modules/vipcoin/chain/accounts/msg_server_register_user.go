package accounts

import (
	"git.ooo.ua/vipcoin/chain/x/accounts/types"
	wallets "git.ooo.ua/vipcoin/chain/x/wallets/types"
	juno "github.com/forbole/juno/v2/types"
)

// handleMsgRegisterUser allows to properly handle a handleMsgRegisterUser
func (m *Module) handleMsgRegisterUser(tx *juno.Tx, index int, msg *types.MsgRegisterUser) error {
	if err := m.accountRepo.SaveRegisterUser(msg, tx.TxHash); err != nil {
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

	if err := m.accountRepo.SaveAccounts(&newAcc); err != nil {
		return err
	}

	// create wallets
	holder := wallets.Wallet{
		Address:        msg.HolderWallet,
		AccountAddress: msg.Address,
		Kind:           wallets.WALLET_KIND_HOLDER,
		State:          wallets.WALLET_STATE_ACTIVE,
		Extras:         msg.HolderWalletExtras,
		Default:        true,
	}

	if err := m.walletsRepo.SaveWallets(&holder); err != nil {
		return err
	}

	refReward := wallets.Wallet{
		Address:        msg.RefRewardWallet,
		AccountAddress: msg.Address,
		Kind:           wallets.WALLET_KIND_REFERRER_REWARD,
		State:          wallets.WALLET_STATE_ACTIVE,
		Extras:         msg.RefRewardWalletExtras,
		Default:        false,
	}

	return m.walletsRepo.SaveWallets(&refReward)
}
