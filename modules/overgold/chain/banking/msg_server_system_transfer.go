package banking

import (
	"strings"

	"git.ooo.ua/vipcoin/chain/x/banking/types"
	"git.ooo.ua/vipcoin/lib/filter"
	sdk "github.com/cosmos/cosmos-sdk/types"
	juno "github.com/forbole/juno/v2/types"

	dbtypes "github.com/forbole/bdjuno/v2/database/types"
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

	var hasNeg bool

	oldBalanceWalletFrom := walletFrom[0].Balance
	walletFrom[0].Balance, hasNeg = walletFrom[0].Balance.SafeSub(sdk.NewCoins(coin))
	if hasNeg {
		walletFrom[0].Balance = oldBalanceWalletFrom

		//return errs.Internal{Cause: fmt.Sprintf(
		//	"failed to transfer %s coins from wallet %s to wallet %s: insufficient funds. "+
		//		"The balance of wallet %s is %s coins.[TX_HASH: %s]",
		//	coin.String(), msg.WalletFrom, msg.WalletTo, msg.WalletFrom, oldBalanceWalletFrom.String(), tx.TxHash),
		//}
	}

	if err = m.walletsRepo.UpdateWallets(walletFrom...); err != nil {
		return err
	}

	walletTo[0].Balance = walletTo[0].Balance.Add(coin)
	if err = m.walletsRepo.UpdateWallets(walletTo...); err != nil {
		return err
	}

	if err = m.bankingRepo.SaveSystemTransfers(transfer); err != nil {
		return err
	}

	return m.bankingRepo.SaveMsgSystemTransfers(msg, tx.TxHash)
}
