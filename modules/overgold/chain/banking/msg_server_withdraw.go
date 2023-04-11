package banking

import (
	"strings"

	"git.ooo.ua/vipcoin/chain/x/banking/types"
	"git.ooo.ua/vipcoin/lib/filter"
	sdk "github.com/cosmos/cosmos-sdk/types"
	juno "github.com/forbole/juno/v2/types"

	dbtypes "github.com/forbole/bdjuno/v2/database/types"
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

	coins := sdk.NewCoins(sdk.NewCoin(msg.Asset, sdk.NewIntFromUint64(msg.Amount)))

	// TODO dev only fix for enormous coins amount
	oldBalance := wallet[0].Balance
	var hasNeg bool
	// sub coins from wallet balance
	// wallet[0].Balance = wallet[0].Balance.Sub(coins) //

	wallet[0].Balance, hasNeg = wallet[0].Balance.SafeSub(coins)
	// If the balance is negative [wallet amount - withdraw amount = is negative], wallet amount remains the same.
	if hasNeg {
		wallet[0].Balance = oldBalance
	}

	if err = m.walletsRepo.UpdateWallets(wallet...); err != nil {
		return err
	}

	// add issuance balance in asset
	asset[0].Withdrawn += msg.Amount

	// TODO dev only fix for enormous coins amount
	if asset[0].InCirculation >= msg.Amount { // remove
		asset[0].InCirculation -= msg.Amount
	} // remove

	if err = m.assetRepo.UpdateAssets(asset...); err != nil {
		return err
	}

	if err = m.bankingRepo.SaveWithdraws(withdraw); err != nil {
		return err
	}

	return m.bankingRepo.SaveMsgWithdraw(msg, tx.TxHash)
}
