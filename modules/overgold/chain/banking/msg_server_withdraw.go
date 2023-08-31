package banking

import (
	"errors"
	"strings"

	"git.ooo.ua/vipcoin/chain/x/banking/types"
	"git.ooo.ua/vipcoin/lib/errs"
	"git.ooo.ua/vipcoin/lib/filter"
	sdk "github.com/cosmos/cosmos-sdk/types"
	juno "github.com/forbole/juno/v3/types"

	dbtypes "github.com/forbole/bdjuno/v3/database/types"
)

// handleMsgWithdraw allows to properly handle a handleMsgWithdraw
func (m *Module) handleMsgWithdraw(tx *juno.Tx, index int, msg *types.MsgWithdraw) error {
	msg.Wallet = strings.ToLower(msg.Wallet)
	msg.Asset = strings.ToLower(msg.Asset)

	withdraw, err := getWithdrawFromTx(tx, msg)
	if err != nil {
		return err
	}

	wallet, err := m.walletsRepo.GetWallets(filter.NewFilter().SetArgument(dbtypes.FieldAddress, msg.Wallet))
	switch {
	case err != nil:
		return err
	case len(wallet) != 1:
		return types.ErrInvalidAddressField
	}

	asset, err := m.assetRepo.GetAssets(filter.NewFilter().SetArgument(dbtypes.FieldName, msg.Asset))
	switch {
	case err != nil:
		return err
	case len(asset) != 1:
		return types.ErrInvalidAddressField
	}

	if err := m.bankingRepo.SaveWithdraws(withdraw); err != nil {
		if errors.As(err, &errs.AlreadyExists{}) {
			// Transfer already exists, it's ok
			return nil
		}
		return err
	}

	coins := sdk.NewCoins(sdk.NewCoin(msg.Asset, sdk.NewIntFromUint64(msg.Amount)))

	balanceFrom := wallet[0].Balance.AmountOf(msg.Asset).Uint64()
	if balanceFrom >= sdk.NewCoin(msg.Asset, sdk.NewIntFromUint64(msg.Amount)).Amount.Uint64() {
		// subtract coins from sender wallet balance
		wallet[0].Balance = wallet[0].Balance.Sub(coins)
		if err := m.walletsRepo.UpdateWallets(wallet...); err != nil {
			return err
		}
	}

	// add issuance balance in asset
	asset[0].Withdrawn += msg.Amount
	asset[0].InCirculation -= msg.Amount

	if err := m.assetRepo.UpdateAssets(asset...); err != nil {
		return err
	}

	return m.bankingRepo.SaveMsgWithdraw(msg, tx.TxHash)
}
