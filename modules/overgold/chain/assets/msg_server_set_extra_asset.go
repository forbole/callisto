package assets

import (
	assetstypes "git.ooo.ua/vipcoin/chain/x/assets/types"
	"git.ooo.ua/vipcoin/lib/filter"
	juno "github.com/forbole/juno/v3/types"

	dbtypes "github.com/forbole/bdjuno/v3/database/types"
)

// handleMsgSetExtraAsset allows to properly handle a handleMsgSetExtraAsset
func (m *Module) handleMsgSetExtraAsset(tx *juno.Tx, index int, msg *assetstypes.MsgAssetSetExtra) error {
	assets, err := m.assetRepo.GetAssets(filter.NewFilter().SetArgument(dbtypes.FieldName, msg.Name))
	if err != nil {
		return err
	}

	if len(assets) != 1 {
		return assetstypes.ErrInvalidNameField
	}

	assets[0].Extras = msg.Extras

	if err := m.assetRepo.UpdateAssets(assets...); err != nil {
		return err
	}

	return m.assetRepo.SaveExtraAsset(msg, tx.TxHash)
}
