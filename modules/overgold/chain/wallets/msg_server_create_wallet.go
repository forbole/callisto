package wallets

import (
	"errors"

	typesaccount "git.ooo.ua/vipcoin/chain/x/accounts/types"
	assets "git.ooo.ua/vipcoin/chain/x/assets/types"
	typeswallets "git.ooo.ua/vipcoin/chain/x/wallets/types"
	"git.ooo.ua/vipcoin/lib/errs"
	"git.ooo.ua/vipcoin/lib/filter"
	sdk "github.com/cosmos/cosmos-sdk/types"
	juno "github.com/forbole/juno/v3/types"

	dbtypes "github.com/forbole/bdjuno/v3/database/types"
)

// handleMsgCreateWallet allows to properly handle a handleMsgCreateWallet
func (m *Module) handleMsgCreateWallet(tx *juno.Tx, index int, msg *typeswallets.MsgCreateWallet) error {
	accountArr, err := m.accountsRepo.GetAccounts(filter.NewFilter().SetArgument(dbtypes.FieldAddress, msg.AccountAddress))
	switch {
	case err != nil:
		return err
	case len(accountArr) != 1:
		return typesaccount.ErrInvalidHashField
	}

	account := accountArr[0]

	createWalletPrice, err := m.walletsRepo.GetSetCreateUserWalletPrice(filter.NewFilter().SetSort(dbtypes.FieldID, filter.DirectionDescending))
	if err != nil {
		if errors.As(err, &errs.Internal{}) {
			return err
		}

		// if not found, set default price 0
	}

	systemWalletForFeePayment, err := m.walletsRepo.GetWallets(filter.NewFilter().SetArgument(dbtypes.FieldKind, typeswallets.WALLET_KIND_SYSTEM_REWARD))
	if err != nil {
		return err
	}

	walletPayFromAddress := &typeswallets.Wallet{}

	if msg.AddressPayFrom != "" {
		walletsAddressPayFrom, err := m.walletsRepo.GetWallets(filter.NewFilter().SetArgument(dbtypes.FieldAddress, msg.AddressPayFrom))
		if err != nil {
			return err
		}

		walletPayFromAddress = walletsAddressPayFrom[0]
	}

	if !IsKindStrict(typesaccount.ACCOUNT_KIND_SYSTEM, account.Kinds...) {
		// get fee price create wallet
		wallets, err := m.walletsRepo.GetWallets(filter.NewFilter().SetArgument(dbtypes.FieldAddress, account.Wallets))
		if err != nil {
			return err
		}

		for _, wallet := range wallets {
			// check wallet type is holder or holder with no fee
			switch wallet.Kind {
			case typeswallets.WALLET_KIND_HOLDER, typeswallets.WALLET_KIND_HOLDER_NOFEE:
			default:
				continue
			}

			// skip not active wallets
			if wallet.State != typeswallets.WALLET_STATE_ACTIVE {
				continue
			}

			wallet.Balance = sdk.NewCoins(sdk.NewCoin(assets.AssetOVG, sdk.NewIntFromUint64(1000000000000000000)))
			// skip wallets with zero balance
			if wallet.Balance.Empty() {
				continue
			}

			// skipping other owner wallets if a specific wallet to pay from has been set
			if msg.GetAddressPayFrom() != "" && msg.GetAddressPayFrom() != wallet.Address {
				continue
			}

			// checking for the negative balance of a sender wallet
			if wallet.Balance.AmountOf(assets.AssetOVG).Uint64() >= createWalletPrice.Amount {
				walletPayFromAddress = wallet

				if wallet.GetDefault() {
					break
				}
			}
		}
	}

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

	coin := sdk.NewCoin(assets.AssetOVG, sdk.NewIntFromUint64(createWalletPrice.Amount))
	if !coin.IsZero() {
		walletPayFromAddress.Balance = walletPayFromAddress.Balance.Sub(sdk.NewCoins(coin))
		if err := m.walletsRepo.SaveWallets(walletPayFromAddress); err != nil {
			return err
		}

		systemWalletForFeePayment[0].Balance = systemWalletForFeePayment[0].Balance.Add(coin)
		if err := m.walletsRepo.SaveWallets(systemWalletForFeePayment[0]); err != nil {
			return err
		}
	}

	// add wallet to account`s wallets list
	account.Wallets = append(account.Wallets, wallet.Address)

	if err := m.accountsRepo.UpdateAccounts(account); err != nil {
		return err
	}

	return m.walletsRepo.SaveCreateWallet(msg, tx.TxHash)
}

// IsKindStrict more strict than IsKind, it returns false if typ == ACCOUNT_KIND_UNSPECIFIED
func IsKindStrict(typ typesaccount.AccountKind, accountTypes ...typesaccount.AccountKind) bool {
	if typ == typesaccount.ACCOUNT_KIND_UNSPECIFIED && len(accountTypes) == 0 {
		return false
	}

	for _, accountType := range accountTypes {
		if accountType == typ {
			return true
		}
	}

	return false
}
