package banking

import (
	"time"

	"git.ooo.ua/vipcoin/chain/x/banking/types"
	"git.ooo.ua/vipcoin/lib/filter"
	sdk "github.com/cosmos/cosmos-sdk/types"
	juno "github.com/forbole/juno/v2/types"

	dbtypes "github.com/forbole/bdjuno/v2/database/types"
)

// handleMsgSystemTransfer allows to properly handle a handleMsgSystemTransfer
func (m *Module) handleMsgSystemTransfer(tx *juno.Tx, index int, msg *types.MsgSystemTransfer) error {
	if err := m.bankingRepo.SaveMsgSystemTransfers(msg); err != nil {
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

	walletFrom[0].Balance = walletFrom[0].Balance.Sub(sdk.NewCoins(coin))
	if err := m.walletsRepo.UpdateWallets(walletFrom...); err != nil {
		return err
	}

	walletTo[0].Balance = walletTo[0].Balance.Add(coin)
	if err := m.walletsRepo.UpdateWallets(walletTo...); err != nil {
		return err
	}

	time, err := time.Parse(time.RFC3339, tx.Timestamp)
	if err != nil {
		return err
	}

	transfer := &types.SystemTransfer{
		WalletFrom: msg.WalletFrom,
		WalletTo:   msg.WalletTo,
		BaseTransfer: types.BaseTransfer{
			Asset:     msg.Asset,
			Amount:    msg.Amount,
			Kind:      types.TRANSFER_KIND_SYSTEM,
			Extras:    msg.Extras,
			Timestamp: time.Unix(),
			TxHash:    tx.TxHash,
		},
	}
	return m.bankingRepo.SaveSystemTransfers(transfer)
}
