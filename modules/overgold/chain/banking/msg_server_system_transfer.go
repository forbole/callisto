package banking

import (
	"strings"

	"git.ooo.ua/vipcoin/chain/x/banking/types"
	"git.ooo.ua/vipcoin/lib/filter"
	sdk "github.com/cosmos/cosmos-sdk/types"
	juno "github.com/forbole/juno/v3/types"

	dbtypes "github.com/forbole/bdjuno/v3/database/types"
)

// handleMsgSystemTransfer allows to properly handle a handleMsgSystemTransfer
func (m *Module) handleMsgSystemTransfer(tx *juno.Tx, index int, msg *types.MsgSystemTransfer) error {
	msg.WalletFrom = strings.ToLower(msg.WalletFrom)
	msg.WalletTo = strings.ToLower(msg.WalletTo)
	msg.Asset = strings.ToLower(msg.Asset)

	transfer, err := getSystemTransferFromTx(tx, msg)
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

	coin := sdk.NewCoin(msg.Asset, sdk.NewIntFromUint64(msg.Amount))

	if walletFrom[0].Address == walletTo[0].Address {
		// If transfer from and to the same wallet, then just update transfer
		if err = m.bankingRepo.SaveSystemTransfers(transfer); err != nil {
			return err
		}

		return m.bankingRepo.SaveMsgSystemTransfers(msg, tx.TxHash)
	}

	balanceFrom := walletFrom[0].Balance.AmountOf(msg.Asset).Uint64()
	if balanceFrom >= msg.Amount {
		walletFrom[0].Balance = walletFrom[0].Balance.Sub(sdk.NewCoins(coin))
		if err := m.walletsRepo.UpdateWallets(walletFrom...); err != nil {
			return err
		}
	}

	walletTo[0].Balance = walletTo[0].Balance.Add(coin)
	if err := m.walletsRepo.UpdateWallets(walletTo...); err != nil {
		return err
	}

	if err = m.bankingRepo.SaveSystemTransfers(transfer); err != nil {
		return err
	}

	return m.bankingRepo.SaveMsgSystemTransfers(msg, tx.TxHash)
}
