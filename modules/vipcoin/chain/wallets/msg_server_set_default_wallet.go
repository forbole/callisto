package wallets

import (
	typesaccount "git.ooo.ua/vipcoin/chain/x/accounts/types"
	typeswallets "git.ooo.ua/vipcoin/chain/x/wallets/types"
	"git.ooo.ua/vipcoin/lib/filter"
	juno "github.com/forbole/juno/v2/types"

	dbtypes "github.com/forbole/bdjuno/v2/database/types"
)

// handleMsgSetDefaultWallet allows to properly handle a MsgSetDefaultWallet
func (m *Module) handleMsgSetDefaultWallet(tx *juno.Tx, index int, msg *typeswallets.MsgSetDefaultWallet) error {
	if err := m.walletsRepo.SaveDefaultWallets(msg, tx.TxHash); err != nil {
		return err
	}

	targetWallet, err := m.walletsRepo.GetWallets(filter.NewFilter().SetArgument(dbtypes.FieldAddress, msg.Address))
	switch {
	case err != nil:
		return err
	case len(targetWallet) != 1:
		return typeswallets.ErrInvalidAddressField
	}

	account, err := m.accountsRepo.GetAccounts(filter.NewFilter().SetArgument(dbtypes.FieldAddress, targetWallet[0].AccountAddress))
	switch {
	case err != nil:
		return err
	case len(account) != 1:
		return typesaccount.ErrInvalidHashField
	}

	for _, walletAddr := range account[0].Wallets {
		w, err := m.walletsRepo.GetWallets(filter.NewFilter().SetArgument(dbtypes.FieldAddress, walletAddr))
		switch {
		case err != nil:
			return err
		case len(w) != 1:
			return typeswallets.ErrInvalidAddressField
		}

		if !w[0].Default {
			continue
		}
		w[0].Default = false
		if err := m.walletsRepo.UpdateWallets(w...); err != nil {
			return err
		}
		break
	}

	targetWallet[0].Default = true

	return m.walletsRepo.UpdateWallets(targetWallet[0])
}
