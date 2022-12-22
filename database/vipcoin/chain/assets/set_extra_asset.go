package assets

import (
	assetstypes "git.ooo.ua/vipcoin/chain/x/assets/types"
	"git.ooo.ua/vipcoin/lib/errs"
	"git.ooo.ua/vipcoin/lib/filter"

	"github.com/forbole/bdjuno/v2/database/types"
)

// SaveExtraAsset - saves the given extra inside the database
func (r Repository) SaveExtraAsset(msg *assetstypes.MsgAssetSetExtra, transactionHash string) error {
	query := `INSERT INTO overgold_chain_assets_set_extra 
			(transaction_hash, creator, name, extras) 
		VALUES 
			(:transaction_hash, :creator, :name, :extras)`

	if _, err := r.db.NamedExec(query, toSetExtraDatabase(msg, transactionHash)); err != nil {
		return errs.Internal{Cause: err.Error()}
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
		return []*assetstypes.MsgAssetSetExtra{}, errs.Internal{Cause: err.Error()}
	}

	extras := make([]*assetstypes.MsgAssetSetExtra, 0, len(result))
	for _, extra := range result {
		extras = append(extras, toSetExtraDomain(extra))
	}

	return extras, nil
}
