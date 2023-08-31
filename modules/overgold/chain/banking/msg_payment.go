package banking

import (
	"errors"
	"math"
	"strings"

	assets "git.ooo.ua/vipcoin/chain/x/assets/types"
	banking "git.ooo.ua/vipcoin/chain/x/banking/types"
	wallets "git.ooo.ua/vipcoin/chain/x/wallets/types"
	"git.ooo.ua/vipcoin/lib/errs"
	"git.ooo.ua/vipcoin/lib/filter"
	sdk "github.com/cosmos/cosmos-sdk/types"
	juno "github.com/forbole/juno/v3/types"

	dbtypes "github.com/forbole/bdjuno/v3/database/types"
)

// handleMsgPayments allows to properly handle a MsgSetState
func (m *Module) handleMsgPayments(tx *juno.Tx, _ int, msg *banking.MsgPayment) error {
	msg.WalletFrom = strings.ToLower(msg.WalletFrom)
	msg.WalletTo = strings.ToLower(msg.WalletTo)
	msg.Asset = strings.ToLower(msg.Asset)

	asset, err := m.assetRepo.GetAssets(filter.NewFilter().SetArgument(dbtypes.FieldName, msg.Asset))
	switch {
	case err != nil:
		return err
	case len(asset) != 1:
		return assets.ErrNotFoundAsset
	}

	walletFrom, err := m.walletsRepo.GetWallets(filter.NewFilter().SetArgument(dbtypes.FieldAddress, msg.WalletFrom))
	switch {
	case err != nil:
		return err
	case len(walletFrom) != 1:
		return wallets.ErrInvalidAddressField
	}

	walletTo, err := m.walletsRepo.GetWallets(filter.NewFilter().SetArgument(dbtypes.FieldAddress, msg.WalletTo))
	switch {
	case err != nil:
		return err
	case len(walletTo) != 1:
		return wallets.ErrInvalidAddressField
	}

	payment, err := getPaymentFromTx(tx, msg)
	if err != nil {
		return err
	}

	switch payment.Fee {
	case 0:
		err = m.payment(payment, *walletFrom[0], *walletTo[0])
	default:
		err = m.paymentWithFee(tx, payment, *asset[0], *walletFrom[0], *walletTo[0])
	}

	if err != nil {
		return err
	}

	return m.bankingRepo.SaveMsgPayments(msg, tx.TxHash)
}

// payment - creates payment without fee
func (m *Module) payment(
	payment *banking.Payment,
	walletFrom, walletTo wallets.Wallet,
) error {
	if err := m.bankingRepo.SavePayments(payment); err != nil {
		return err
	}

	coin := sdk.NewCoin(payment.BaseTransfer.Asset, sdk.NewIntFromUint64(payment.Amount))

	balanceFrom := walletFrom.Balance.AmountOf(payment.Asset).Uint64()
	if balanceFrom >= coin.Amount.Uint64() {
		// subtract coins from sender wallet balance
		walletFrom.Balance = walletFrom.Balance.Sub(sdk.NewCoins(coin))
		if err := m.walletsRepo.UpdateWallets(&walletFrom); err != nil {
			return err
		}
	}

	// add coins to receiver wallet balance
	walletTo.Balance = walletTo.Balance.Add(coin)

	return m.walletsRepo.UpdateWallets(&walletTo)
}

// paymentWithFee - creates payment with fee
func (m *Module) paymentWithFee(
	tx *juno.Tx,
	payment *banking.Payment,
	asset assets.Asset,
	walletFrom, walletTo wallets.Wallet,
) error {
	systemReward, systemRefReward, err := m.getSystemTransfers(tx, payment.WalletTo)
	if err != nil {
		return err
	}

	// Getting supplementary wallets
	walletsSystemReward, err := m.walletsRepo.GetWallets(filter.NewFilter().SetArgument(dbtypes.FieldAddress, systemReward.WalletTo))
	switch {
	case err != nil:
		return err
	case len(walletsSystemReward) != 1:
		return wallets.ErrInvalidAddressField
	}

	walletsRefReward, err := m.walletsRepo.GetWallets(filter.NewFilter().SetArgument(dbtypes.FieldAddress, systemRefReward.WalletTo))
	switch {
	case err != nil:
		return err
	case len(walletsRefReward) != 1:
		return wallets.ErrInvalidAddressField
	}

	walletsVoid, err := m.walletsRepo.GetWallets(
		filter.NewFilter().
			SetCondition(filter.ConditionAND).
			SetArgument(dbtypes.FieldAccountAddress, walletsSystemReward[0].AccountAddress).
			SetArgument(dbtypes.FieldKind, wallets.WALLET_KIND_VOID))
	switch {
	case err != nil:
		return err
	case len(walletsVoid) != 1:
		return wallets.ErrInvalidKindField
	}

	if err := m.bankingRepo.SaveSystemTransfers(systemReward); err != nil {
		if errors.As(err, &errs.AlreadyExists{}) {
			// Transfer already exists, it's ok
			return nil
		}
		return err
	}

	if err := m.bankingRepo.SaveSystemTransfers(systemRefReward); err != nil {
		if errors.As(err, &errs.AlreadyExists{}) {
			// Transfer already exists, it's ok
			return nil
		}
		return err
	}

	walletSystemReward := walletsSystemReward[0]
	walletRefReward := walletsRefReward[0]
	walletVoid := walletsVoid[0]

	// ----- General payment -----
	var (
		feeRaw          = float64(payment.Amount) / 100.0 * (float64(asset.Properties.FeePercent) / 100.0) // FeePercent 100 = 1%
		feeSysRewardRaw = feeRaw / 100.0 * 50.0
		feeRefRewardRaw = feeRaw / 100.0 * 25.0
	)

	var (
		fee          = uint64(math.Round(feeRaw))
		feeSysReward = uint64(math.Round(feeSysRewardRaw))
		feeRefReward = uint64(math.Round(feeRefRewardRaw))
		feeVoid      = fee - (feeSysReward + feeRefReward)
	)

	var (
		coin             = sdk.NewCoin(asset.Name, sdk.NewIntFromUint64(payment.Amount))
		coinAfterFee     = sdk.NewCoin(asset.Name, sdk.NewIntFromUint64(payment.Amount-fee))
		coinFeeSysReward = sdk.NewCoin(asset.Name, sdk.NewIntFromUint64(feeSysReward))
		coinFeeRefReward = sdk.NewCoin(asset.Name, sdk.NewIntFromUint64(feeRefReward))
		coinFeeVoid      = sdk.NewCoin(asset.Name, sdk.NewIntFromUint64(feeVoid))
	)

	balanceFrom := walletFrom.Balance.AmountOf(asset.Name).Uint64()
	if balanceFrom >= coin.Amount.Uint64() {
		walletFrom.Balance = walletFrom.Balance.Sub(sdk.NewCoins(coin))
		if err := m.walletsRepo.UpdateWallets(&walletFrom); err != nil {
			return err
		}
	}

	// add coins to receiver wallet balance
	walletTo.Balance = walletTo.Balance.Add(coinAfterFee)
	if err := m.walletsRepo.UpdateWallets(&walletTo); err != nil {
		return err
	}

	// add coins to system reward wallet balance
	walletSystemReward.Balance = walletSystemReward.Balance.Add(coinFeeSysReward)
	if err := m.walletsRepo.UpdateWallets(walletSystemReward); err != nil {
		return err
	}

	// add coins to referrer (if referrer is empty then it will be system ref reward wallet) wallet balance
	walletRefReward.Balance = walletRefReward.Balance.Add(coinFeeRefReward)
	if err := m.walletsRepo.UpdateWallets(walletRefReward); err != nil {
		return err
	}

	// add coins to void wallet balance
	walletVoid.Balance = walletVoid.Balance.Add(coinFeeVoid)
	if err := m.walletsRepo.UpdateWallets(walletVoid); err != nil {
		return err
	}

	asset.Burned += feeVoid
	asset.InCirculation -= feeVoid
	if err := m.assetRepo.UpdateAssets(&asset); err != nil {
		return err
	}

	return m.bankingRepo.SavePayments(payment)
}
