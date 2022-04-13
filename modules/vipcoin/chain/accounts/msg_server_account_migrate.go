package accounts

import (
	"git.ooo.ua/vipcoin/chain/x/accounts/types"
	xtypes "git.ooo.ua/vipcoin/chain/x/types"
	"git.ooo.ua/vipcoin/lib/filter"
	juno "github.com/forbole/juno/v2/types"

	dbtypes "github.com/forbole/bdjuno/v2/database/types"
)

// handleMsgAccountMigrate allows to properly handle a handleMsgAccountMigrate
func (m *Module) handleMsgAccountMigrate(tx *juno.Tx, index int, msg *types.MsgAccountMigrate) error {
	if err := m.accountRepo.SaveAccountMigrate(msg); err != nil {
		return err
	}

	accountArr, err := m.accountRepo.GetAccounts(filter.NewFilter().SetArgument(dbtypes.FieldHash, msg.Hash))
	switch {
	case err != nil:
		return err
	case len(accountArr) != 1:
		return types.ErrInvalidHashField
	}

	account := accountArr[0]

	wallets, err := m.walletsRepo.GetWallets(filter.NewFilter().SetArgument(dbtypes.FieldAccountAddress, account.Address))
	if err != nil {
		return err
	}

	publicKey, err := types.PubKeyFromString(msg.PublicKey)
	if err != nil {
		return types.ErrInvalidPublicKeyField
	}

	account.PublicKey, err = types.PubKeyToAny(publicKey)
	if err != nil {
		return types.ErrInvalidPublicKeyField
	}

	oldAddr := account.Address
	account.Address = msg.Address

	// Change address in affiliates
	for _, affiliate := range account.Affiliates {
		// Get affiliate account
		affiliateAcc, err := m.accountRepo.GetAccounts(filter.NewFilter().SetArgument(dbtypes.FieldAddress, affiliate.Address))
		switch {
		case err != nil:
			return err
		case len(affiliateAcc) != 1:
			return types.ErrInvalidAddressField
		}

		if affiliateAcc[0].Address == "" {
			// skip empty affiliate
			continue
		}

		// search obsolete address
		for i, a := range affiliateAcc[0].Affiliates {
			if xtypes.GetSDKAddress(a.Address).Equals(xtypes.GetSDKAddress(oldAddr)) {
				// Update address to new one
				a.Address = msg.Address
				affiliateAcc[0].Affiliates[i] = a
				// Break loop
				break
			}
		}

		if err := m.accountRepo.UpdateAccounts(affiliateAcc...); err != nil {
			return err
		}
	}

	if err := m.walletsRepo.DeleteWallets(wallets...); err != nil {
		return err
	}

	if err := m.accountRepo.UpdateAccounts(account); err != nil {
		return err
	}

	// change address in wallets
	for index := range wallets {
		if wallets[index].Address == "" {
			// skip empty wallet
			continue
		}

		wallets[index].AccountAddress = msg.Address
	}

	return m.walletsRepo.SaveWallets(wallets...)
}
