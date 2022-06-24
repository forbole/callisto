package assets

import (
	assetstypes "git.ooo.ua/vipcoin/chain/x/assets/types"
	"git.ooo.ua/vipcoin/lib/filter"
	juno "github.com/forbole/juno/v2/types"

	dbtypes "github.com/forbole/bdjuno/v2/database/types"
)

// handleMsgSetExtraAsset allows to properly handle a handleMsgSetExtraAsset
func (m *Module) handleMsgSetExtraAsset(tx *juno.Tx, index int, msg *assetstypes.MsgAssetSetExtra) error {
	if err := m.assetRepo.SaveExtraAsset(msg, tx.TxHash); err != nil {
		return err
	}

	assets, err := m.assetRepo.GetAssets(filter.NewFilter().SetArgument(dbtypes.FieldName, msg.Name))
	if err != nil {
		return err
	}

	if len(assets) != 1 {
		return assetstypes.ErrInvalidNameField
	}

	assets[0].Extras = msg.Extras

	return m.assetRepo.UpdateAssets(assets...)
}
