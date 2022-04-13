package banking

import (
	"math"
	"time"

	accounts "git.ooo.ua/vipcoin/chain/x/accounts/types"
	assets "git.ooo.ua/vipcoin/chain/x/assets/types"
	banking "git.ooo.ua/vipcoin/chain/x/banking/types"
	wallets "git.ooo.ua/vipcoin/chain/x/wallets/types"
	"git.ooo.ua/vipcoin/lib/filter"
	sdk "github.com/cosmos/cosmos-sdk/types"
	juno "github.com/forbole/juno/v2/types"

	dbtypes "github.com/forbole/bdjuno/v2/database/types"
)

// handleMsgPayments allows to properly handle a MsgSetState
func (m *Module) handleMsgPayments(tx *juno.Tx, index int, msg *banking.MsgPayment) error {
	if err := m.bankingRepo.SaveMsgPayments(msg); err != nil {
		return err
	}

	asset, err := m.assetRepo.GetAssets(filter.NewFilter().SetArgument(dbtypes.FieldName, msg.Asset))
	switch {
	case err != nil:
		return err
	case len(asset) != 1:
		return banking.ErrNotFoundAsset
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

	accountFrom, err := m.accountsRepo.GetAccounts(filter.NewFilter().SetArgument(dbtypes.FieldAddress, walletFrom[0].AccountAddress))
	switch {
	case err != nil:
		return err
	case len(accountFrom) != 1:
		return banking.ErrInvalidAddressField
	}

	accountTo, err := m.accountsRepo.GetAccounts(filter.NewFilter().SetArgument(dbtypes.FieldAddress, walletTo[0].AccountAddress))
	switch {
	case err != nil:
		return err
	case len(accountTo) != 1:
		return banking.ErrInvalidAddressField
	}

	var payment *banking.Payment

	switch isChargeFee(accountFrom[0], accountTo[0], walletFrom[0].Kind, walletTo[0].Kind) {
	case false:
		payment, err = m.payment(msg, *asset[0], *walletFrom[0], *walletTo[0])
		if err != nil {
			return err
		}
	default:
		payment, err = m.paymentWithFee(msg, *asset[0], *walletFrom[0], *walletTo[0])
		if err != nil {
			return err
		}
	}

	time, err := time.Parse(time.RFC3339, tx.Timestamp)
	if err != nil {
		return err
	}

	payment.BaseTransfer.Timestamp = time.Unix()
	payment.BaseTransfer.TxHash = tx.TxHash

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

// payment - creates payment without fee
func (m *Module) payment(
	msg *banking.MsgPayment,
	asset assets.Asset,
	walletFrom, walletTo wallets.Wallet,
) (*banking.Payment, error) {
	coin := sdk.NewCoin(asset.Name, sdk.NewIntFromUint64(msg.Amount))

	// subtract coins from sender wallet balance
	walletFrom.Balance = walletFrom.Balance.Sub(sdk.NewCoins(coin))
	if err := m.walletsRepo.UpdateWallets(&walletFrom); err != nil {
		return &banking.Payment{}, err
	}

	// add coins to receiver wallet balance
	walletTo.Balance = walletTo.Balance.Add(coin)
	if err := m.walletsRepo.UpdateWallets(&walletTo); err != nil {
		return &banking.Payment{}, err
	}

	payment := &banking.Payment{
		WalletFrom: msg.WalletFrom,
		WalletTo:   msg.WalletTo,
		Fee:        0,
		BaseTransfer: banking.BaseTransfer{
			Asset:  msg.Asset,
			Amount: msg.Amount,
			Extras: msg.Extras,
			Kind:   banking.TRANSFER_KIND_PAYMENT,
		},
	}

	if walletFrom.Kind == wallets.WALLET_KIND_SYSTEM_DEFERRED {
		payment.Kind = banking.TRANSFER_KIND_DEFERRED
	}

	return payment, nil
}

// paymentWithFee - creates payment with fee
func (m *Module) paymentWithFee(
	msg *banking.MsgPayment,
	asset assets.Asset,
	walletFrom, walletTo wallets.Wallet,
) (*banking.Payment, error) {
	// ----- General payment -----
	var (
		feeRaw       = float64(msg.Amount) / 100.0 * (float64(asset.Properties.FeePercent) / 100.0) // FeePercent 100 = 1%
		fee          = uint64(math.Round(feeRaw))
		coin         = sdk.NewCoin(asset.Name, sdk.NewIntFromUint64(msg.Amount))
		coinAfterFee = sdk.NewCoin(asset.Name, sdk.NewIntFromUint64(msg.Amount-fee))
	)

	// subtract coins from sender wallet balance
	walletFrom.Balance = walletFrom.Balance.Sub(sdk.NewCoins(coin))
	if err := m.walletsRepo.UpdateWallets(&walletFrom); err != nil {
		return &banking.Payment{}, err
	}

	// add coins to receiver wallet balance
	walletTo.Balance = walletTo.Balance.Add(coinAfterFee)
	if err := m.walletsRepo.UpdateWallets(&walletTo); err != nil {
		return &banking.Payment{}, err
	}

	payment := &banking.Payment{
		WalletFrom: msg.WalletFrom,
		WalletTo:   msg.WalletTo,
		Fee:        fee,
		BaseTransfer: banking.BaseTransfer{
			Asset:  msg.Asset,
			Amount: msg.Amount,
			Extras: msg.Extras,
			Kind:   banking.TRANSFER_KIND_PAYMENT,
		},
	}

	if walletFrom.Kind == wallets.WALLET_KIND_SYSTEM_DEFERRED {
		payment.Kind = banking.TRANSFER_KIND_DEFERRED
	}

	return payment, nil
}
