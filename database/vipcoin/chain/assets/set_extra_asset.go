package assets

import (
	assetstypes "git.ooo.ua/vipcoin/chain/x/assets/types"
	"git.ooo.ua/vipcoin/lib/filter"

	"github.com/forbole/bdjuno/v2/database/types"
)

// SaveExtraAsset - saves the given extra inside the database
func (r Repository) SaveExtraAsset(msg ...*assetstypes.MsgAssetSetExtra) error {
	if len(msg) == 0 {
		return nil
	}

	query := `INSERT INTO vipcoin_chain_assets_set_extra 
			(creator, name, extras) 
		VALUES 
			(:creator, :name, :extras)`

	if _, err := r.db.NamedExec(query, toSetExtrasDatabase(msg...)); err != nil {
		return err
	}

	return nil
}

// GetExtraAsset - get the given extra from database
func (r Repository) GetExtraAsset(assetFilter filter.Filter) ([]*assetstypes.MsgAssetSetExtra, error) {
	query, args := assetFilter.Build(
		tableSetExtrasAsset,
		types.FieldCreator, types.FieldName, types.FieldExtras,
	)

	var result []types.DBAssetSetExtra
	if err := r.db.Select(&result, query, args...); err != nil {
		return []*assetstypes.MsgAssetSetExtra{}, err
	}

	extras := make([]*assetstypes.MsgAssetSetExtra, 0, len(result))
	for _, extra := range result {
		extras = append(extras, toSetExtraDomain(extra))
	}

	return extras, nil
}
