package banking

import (
	"math"
	"time"

	accounts "git.ooo.ua/vipcoin/chain/x/accounts/types"
	banking "git.ooo.ua/vipcoin/chain/x/banking/types"
	wallets "git.ooo.ua/vipcoin/chain/x/wallets/types"
	"git.ooo.ua/vipcoin/lib/filter"
	juno "github.com/forbole/juno/v2/types"

	dbtypes "github.com/forbole/bdjuno/v2/database/types"
)

// handleMsgPayments allows to properly handle a MsgSetState
func (m *Module) handleMsgPayments(tx *juno.Tx, index int, msg *banking.MsgPayment) error {
	if err := m.bankingRepo.SaveMsgPayments(msg); err != nil {
		return err
	}

	accountFrom, err := m.accountsRepo.GetAccounts(filter.NewFilter().SetArgument(dbtypes.FieldAddress, msg.WalletFrom))
	switch {
	case err != nil:
		return err
	case len(accountFrom) != 1:
		return banking.ErrInvalidAddressField
	}

	accountTo, err := m.accountsRepo.GetAccounts(filter.NewFilter().SetArgument(dbtypes.FieldAddress, msg.WalletTo))
	switch {
	case err != nil:
		return err
	case len(accountTo) != 1:
		return banking.ErrInvalidAddressField
	}

	walletFrom, err := m.walletsRepo.GetWallets(filter.NewFilter().SetArgument(dbtypes.FieldAddress, msg.WalletFrom))
	switch {
	case err != nil:
		return err
	case len(walletFrom) != 1:
		return banking.ErrInvalidAddressField
	}

	walletTo, err := m.walletsRepo.GetWallets(filter.NewFilter().SetArgument(dbtypes.FieldAddress, msg.WalletTo))
	switch {
	case err != nil:
		return err
	case len(walletTo) != 1:
		return banking.ErrInvalidAddressField
	}

	time, err := time.Parse(time.RFC3339, tx.Timestamp)
	if err != nil {
		return err
	}

	var fee uint64

	switch isChargeFee(accountFrom[0], accountTo[0], walletFrom[0].Kind, walletTo[0].Kind) {
	case false:
		fee = 0
	default:
		assetsByName, err := m.assetRepo.GetAssets(filter.NewFilter().SetArgument(dbtypes.FieldName, msg.Asset))
		switch {
		case err != nil:
			return err
		case len(assetsByName) != 1:
			return banking.ErrInvalidAddressField
		}

		feeRaw := float64(msg.Amount) / 100.0 * (float64(assetsByName[0].Properties.FeePercent) / 100.0) // FeePercent 100 = 1%
		fee = uint64(math.Round(feeRaw))
	}

	payment := &banking.Payment{
		BaseTransfer: banking.BaseTransfer{
			Asset:     msg.Asset,
			Amount:    msg.Amount,
			Kind:      banking.TRANSFER_KIND_PAYMENT,
			Timestamp: time.Unix(),
			Extras:    msg.Extras,
			TxHash:    tx.TxHash,
		},
		WalletFrom: msg.WalletFrom,
		WalletTo:   msg.WalletTo,
		Fee:        fee,
	}

	return m.bankingRepo.SavePayments(payment)
}

// isChargeFee - checks whether it is necessary to charge a commission from the accounts
func isChargeFee(from, to *accounts.Account, toWalletKind, fromWalletKind wallets.WalletKind) bool {
	// allow use payment for same account
	if from.Address == to.Address {
		return false
	}

	// don't charge fee if payment sent from general to system account
	if accounts.IsKind(accounts.ACCOUNT_KIND_GENERAL, from.Kinds...) &&
		accounts.IsKind(accounts.ACCOUNT_KIND_SYSTEM, to.Kinds...) {
		return false
	}

	if fromWalletKind == wallets.WALLET_KIND_HOLDER_NOFEE {
		return false
	}

	if toWalletKind == wallets.WALLET_KIND_HOLDER_NOFEE {
		return false
	}

	return true
}
