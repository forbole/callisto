package banking

import (
	"strings"
	"time"

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

	if err := m.bankingRepo.SaveMsgWithdraw(msg); err != nil {
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

	// sub coins from wallet balance
	wallet[0].Balance, _ = wallet[0].Balance.SafeSub(coins)
	if err := m.walletsRepo.UpdateWallets(wallet...); err != nil {
		return err
	}

	// add issuance balance in asset
	asset[0].Withdrawn += msg.Amount
	asset[0].InCirculation -= msg.Amount
	if err := m.assetRepo.UpdateAssets(asset...); err != nil {
		return err
	}

	msgTime, err := time.Parse(time.RFC3339, tx.Timestamp)
	if err != nil {
		return err
	}

	withdraw := &types.Withdraw{
		Wallet: msg.Wallet,
		BaseTransfer: types.BaseTransfer{
			Asset:     msg.Asset,
			Amount:    msg.Amount,
			Kind:      types.TRANSFER_KIND_WITHDRAW,
			Timestamp: msgTime.Unix(),
			Extras:    msg.Extras,
			TxHash:    tx.TxHash,
		},
	}

	return m.bankingRepo.SaveWithdraws(withdraw)
}
