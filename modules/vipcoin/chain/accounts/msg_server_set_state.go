package accounts

import (
	"git.ooo.ua/vipcoin/chain/x/accounts/types"
	assets "git.ooo.ua/vipcoin/chain/x/assets/types"
	wallets "git.ooo.ua/vipcoin/chain/x/wallets/types"
	"git.ooo.ua/vipcoin/lib/filter"
	sdk "github.com/cosmos/cosmos-sdk/types"
	dbtypes "github.com/forbole/bdjuno/v2/database/types"
	juno "github.com/forbole/juno/v2/types"
)

// handleMsgSetState allows to properly handle a handleMsgSetState
func (m *Module) handleMsgSetState(tx *juno.Tx, index int, msg *types.MsgSetState) error {
	acc, err := m.accountRepo.GetAccounts(filter.NewFilter().SetArgument(dbtypes.FieldHash, msg.Hash))
	switch {
	case err != nil:
		return err
	case len(acc) != 1:
		return types.ErrInvalidHashField
	}

	account := acc[0]

	switch msg.State {
	case types.ACCOUNT_STATE_BLOCKED, types.ACCOUNT_STATE_ACTIVE:
		assetArr, err := m.assetRepo.GetAssets(filter.NewFilter().SetArgument(dbtypes.FieldName, assets.AssetVCG))
		switch {
		case err != nil:
			return err
		case len(assetArr) != 1:
			return types.ErrNotFoundAsset
		}

		asset := assetArr[0]

		// Get all user wallets balance
		totalAmount := sdk.NewInt(0)
		for _, walletAdds := range account.Wallets {
			walletArr, err := m.walletsRepo.GetWallets(filter.NewFilter().SetArgument(dbtypes.FieldAddress, walletAdds))
			switch {
			case err != nil:
				return err
			case len(walletArr) != 1:
				return types.ErrInvalidAddressField
			}

			wallet := walletArr[0]

			if wallet.Address == "" {
				continue
			}

			if wallet.State == wallets.WALLET_STATE_DELETED {
				continue
			}

			switch msg.State {
			case types.ACCOUNT_STATE_ACTIVE:
				wallet.State = wallets.WALLET_STATE_ACTIVE
				if err := m.walletsRepo.UpdateWallets(wallet); err != nil {
					return err
				}

			case types.ACCOUNT_STATE_BLOCKED:
				wallet.State = wallets.WALLET_STATE_BLOCKED
				if err := m.walletsRepo.UpdateWallets(wallet); err != nil {
					return err
				}
			}

			coins := wallet.Balance.AmountOf(assets.AssetVCG)
			totalAmount = totalAmount.Add(coins)
		}

		commitAsset := false
		// add or subtotal amount from asset circulation
		if msg.State == types.ACCOUNT_STATE_BLOCKED {
			// Subtracting from asset balance values from blocked wallets
			newBalance := asset.InCirculation - totalAmount.Uint64()
			asset.InCirculation = newBalance
			commitAsset = true
		}

		if account.State == types.ACCOUNT_STATE_BLOCKED && msg.State == types.ACCOUNT_STATE_ACTIVE {
			// Adding to asset balance values from blocked wallets
			newBalance := asset.InCirculation + totalAmount.Uint64()
			asset.InCirculation = newBalance
			commitAsset = true
		}

		if commitAsset {
			// Update asset data
			if err := m.assetRepo.UpdateAssets(asset); err != nil {
				return err
			}
		}
	}

	account.State = msg.State

	if err := m.accountRepo.UpdateAccounts(account); err != nil {
		return err
	}

	return m.accountRepo.SaveState(msg, tx.TxHash)
}
