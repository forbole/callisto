package assets

import (
	assetstypes "git.ooo.ua/vipcoin/chain/x/assets/types"
	"git.ooo.ua/vipcoin/lib/filter"
	juno "github.com/forbole/juno/v2/types"

	dbtypes "github.com/forbole/bdjuno/v2/database/types"
)

// handleMsgManageAsset allows to properly handle a handleMsgCreateAsset
func (m *Module) handleMsgManageAsset(tx *juno.Tx, index int, msg *assetstypes.MsgAssetManage) error {
	assetsArr, err := m.assetRepo.GetAssets(filter.NewFilter().SetArgument(dbtypes.FieldName, msg.Name))
	switch {
	case err != nil:
		return err
	case len(assetsArr) != 1:
		return assetstypes.ErrInvalidNameField
	}

	asset := assetsArr[0]

	asset.Properties = msg.Properties
	asset.State = msg.State
	asset.Policies = msg.Policies

	if msg.Issued != 0 {
		asset.Issued = msg.Issued
	}

	if msg.Withdrawn != 0 {
		asset.Withdrawn = msg.Withdrawn
	}

	if msg.Burned != 0 {
		asset.Burned = msg.Burned
	}

	if msg.InCirculation != 0 {
		asset.InCirculation = msg.InCirculation
	}

	if err := m.assetRepo.UpdateAssets(asset); err != nil {
		return err
	}

	return m.assetRepo.SaveManageAsset(msg, tx.TxHash)
}
