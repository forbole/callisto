package banking

import (
	"errors"

	"git.ooo.ua/vipcoin/chain/x/banking/types"
	"git.ooo.ua/vipcoin/lib/errs"
	"git.ooo.ua/vipcoin/lib/filter"
	sdk "github.com/cosmos/cosmos-sdk/types"
	juno "github.com/forbole/juno/v3/types"

	dbtypes "github.com/forbole/bdjuno/v3/database/types"
)

// handleMsgSystemRewardTransfer allows to properly handle a MsgSystemRewardTransfer
func (m *Module) handleMsgSystemRewardTransfer(tx *juno.Tx, index int, msg *types.MsgSystemRewardTransfer) error {
	system := &types.MsgSystemTransfer{
		WalletFrom: msg.WalletFrom,
		WalletTo:   msg.WalletTo,
		Amount:     msg.Amount,
		Extras:     msg.Extras,
	}

	transfer, err := getSystemTransferFromTx(tx, system)
	if err != nil {
		return err
	}

	walletFrom, err := m.walletsRepo.GetWallets(filter.NewFilter().SetArgument(dbtypes.FieldAddress, msg.WalletFrom))
	switch {
	case err != nil:
		return err
	case len(walletFrom) != 1:
		return types.ErrInvalidAddressField
	}

	walletTo, err := m.walletsRepo.GetWallets(filter.NewFilter().SetArgument(dbtypes.FieldAddress, msg.WalletTo))
	switch {
	case err != nil:
		return err
	case len(walletTo) != 1:
		return types.ErrInvalidAddressField
	}

	if err := m.bankingRepo.SaveSystemTransfers(transfer); err != nil {
		if errors.As(err, &errs.AlreadyExists{}) {
			// Transfer already exists, it's ok
			return nil
		}
		return err
	}

	coin := sdk.NewCoin(msg.Asset, sdk.NewIntFromUint64(msg.Amount))

	balanceFrom := walletFrom[0].Balance.AmountOf(msg.Asset).Uint64()
	if balanceFrom >= coin.Amount.Uint64() {
		// subtract coins from sender wallet balance
		if err := m.walletsRepo.UpdateWallets(walletFrom...); err != nil {
			return err
		}
	}

	walletTo[0].Balance = walletTo[0].Balance.Add(coin)
	if err := m.walletsRepo.UpdateWallets(walletTo...); err != nil {
		return err
	}

	return m.bankingRepo.SaveMsgSystemRewardTransfers(msg, tx.TxHash)
}
