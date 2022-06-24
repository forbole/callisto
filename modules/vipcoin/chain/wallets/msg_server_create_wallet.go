package wallets

import (
	typesaccount "git.ooo.ua/vipcoin/chain/x/accounts/types"
	typeswallets "git.ooo.ua/vipcoin/chain/x/wallets/types"
	"git.ooo.ua/vipcoin/lib/filter"
	juno "github.com/forbole/juno/v2/types"

	dbtypes "github.com/forbole/bdjuno/v2/database/types"
)

// handleMsgCreateWallet allows to properly handle a handleMsgCreateWallet
func (m *Module) handleMsgCreateWallet(tx *juno.Tx, index int, msg *typeswallets.MsgCreateWallet) error {
	if err := m.walletsRepo.SaveCreateWallet(msg, tx.TxHash); err != nil {
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

	isDefault := func() bool {
		// check if account has wallet with type "holder", if its not then wallet will be default wallet
		if msg.Kind == typeswallets.WALLET_KIND_HOLDER {
			for _, walletAddr := range account.Wallets {
				w, err := m.walletsRepo.GetWallets(filter.NewFilter().SetArgument(dbtypes.FieldAddress, walletAddr))
				if err != nil {
					return false
				}

				if w[0].Address == "" { // skip empty
					continue
				}

				if w[0].Kind == typeswallets.WALLET_KIND_HOLDER {
					return false
				}
			}
			return true
		}
		return false
	}

	wallet := typeswallets.Wallet{
		Address:        msg.Address,
		AccountAddress: msg.AccountAddress,
		Kind:           msg.Kind,
		State:          msg.State,
		Extras:         msg.Extras,
		Default:        isDefault(),
	}

	if err := m.walletsRepo.SaveWallets(&wallet); err != nil {
		return err
	}

	// add wallet to account`s wallets list
	account.Wallets = append(account.Wallets, wallet.Address)

	return m.accountsRepo.UpdateAccounts(account)
}
