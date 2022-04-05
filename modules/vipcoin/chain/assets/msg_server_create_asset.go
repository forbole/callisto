package assets

import (
	assetstypes "git.ooo.ua/vipcoin/chain/x/assets/types"

	juno "github.com/forbole/juno/v2/types"
)

// handleMsgCreateAsset allows to properly handle a handleMsgCreateAsset
func (m *Module) handleMsgCreateAsset(tx *juno.Tx, index int, msg *assetstypes.MsgAssetCreate) error {
	if err := m.assetRepo.SaveCreateAsset(msg); err != nil {
		return err
	}

	newAsset := assetstypes.Asset{
		Name:       msg.Name,
		Issuer:     msg.Issuer,
		Policies:   msg.Policies,
		State:      msg.State,
		Properties: msg.Properties,
		Extras:     msg.Extras,
	}

	return m.assetRepo.SaveAssets(&newAsset)
}
