package wallets

import (
	typesaccount "git.ooo.ua/vipcoin/chain/x/accounts/types"
	typeswallets "git.ooo.ua/vipcoin/chain/x/wallets/types"
	"git.ooo.ua/vipcoin/lib/filter"
	juno "github.com/forbole/juno/v2/types"

	dbtypes "github.com/forbole/bdjuno/v2/database/types"
)

// MsgCreateWalletWithBalance allows to properly handle a handleMsgCreateWallet
func (m *Module) MsgCreateWalletWithBalance(tx *juno.Tx, index int, msg *typeswallets.MsgCreateWalletWithBalance) error {
	if err := m.walletsRepo.SaveCreateWalletWithBalance(msg, tx.TxHash); err != nil {
		return err
	}

	accountArr, err := m.accountsRepo.GetAccounts(filter.NewFilter().SetArgument(dbtypes.FieldAddress, msg.AccountAddress))
	switch {
	case err != nil:
		return err
	case len(accountArr) != 1:
		return typesaccount.ErrInvalidHashField
	}

	account := accountArr[0]

	// if another wallet is default, we set it to non-default
	if msg.Default {
		for _, walletAddres := range account.Wallets {
			wallet, err := m.walletsRepo.GetWallets(filter.NewFilter().SetArgument(dbtypes.FieldAddress, walletAddres))
			switch {
			case err != nil:
				return err
			case len(wallet) != 1:
				return typeswallets.ErrInvalidAddressField
			}

			if !wallet[0].Default {
				continue
			}

			wallet[0].Default = false
			if err := m.walletsRepo.UpdateWallets(wallet...); err != nil {
				return err
			}
		}
	}

	wallet := typeswallets.Wallet{
		Address:        msg.Address,
		AccountAddress: msg.AccountAddress,
		Kind:           msg.Kind,
		State:          msg.State,
		Extras:         msg.Extras,
		Default:        msg.Default,
	}

	if !msg.Balance.IsZero() {
		wallet.Balance = msg.Balance
	}

	if err := m.walletsRepo.SaveWallets(&wallet); err != nil {
		return err
	}

	// add wallet to account`s wallets list
	account.Wallets = append(account.Wallets, wallet.Address)

	return m.accountsRepo.UpdateAccounts(account)
}
