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

// handleMsgIssue allows to properly handle a handleMsgIssue
func (m *Module) handleMsgIssue(tx *juno.Tx, index int, msg *types.MsgIssue) error {
	msg.Wallet = strings.ToLower(msg.Wallet)
	msg.Asset = strings.ToLower(msg.Asset)

	if err := m.bankingRepo.SaveMsgIssue(msg); err != nil {
		return err
	}

	asset, err := m.assetRepo.GetAssets(filter.NewFilter().SetArgument(dbtypes.FieldName, msg.Asset))
	switch {
	case err != nil:
		return err
	case len(asset) != 1:
		return types.ErrNotFoundAsset
	}

	wallet, err := m.walletsRepo.GetWallets(filter.NewFilter().SetArgument(dbtypes.FieldAddress, msg.Wallet))
	switch {
	case err != nil:
		return err
	case len(wallet) != 1:
		return types.ErrInvalidAddressField
	}

	coin := sdk.NewCoin(asset[0].Name, sdk.NewIntFromUint64(msg.Amount))

	// add coins to wallet balance
	wallet[0].Balance = wallet[0].Balance.Add(coin)
	if err := m.walletsRepo.UpdateWallets(wallet...); err != nil {
		return err
	}

	// add issuance balance in asset
	asset[0].Issued += msg.Amount
	asset[0].InCirculation += msg.Amount
	if err := m.assetRepo.UpdateAssets(asset...); err != nil {
		return err
	}

	time, err := time.Parse(time.RFC3339, tx.Timestamp)
	if err != nil {
		return err
	}

	issue := &types.Issue{
		Wallet: msg.Wallet,
		BaseTransfer: types.BaseTransfer{
			Asset:     msg.Asset,
			Amount:    msg.Amount,
			Kind:      types.TRANSFER_KIND_ISSUE,
			Timestamp: time.Unix(),
			Extras:    msg.Extras,
			TxHash:    tx.TxHash,
		},
	}
	return m.bankingRepo.SaveIssues(issue)
}
