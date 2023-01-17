package assets

import (
	assetstypes "git.ooo.ua/vipcoin/chain/x/assets/types"
	"git.ooo.ua/vipcoin/lib/errs"
	"git.ooo.ua/vipcoin/lib/filter"

	"github.com/forbole/bdjuno/v2/database/types"
)

// SaveCreateAsset saves the given asset inside the database
func (r *Repository) SaveCreateAsset(msg *assetstypes.MsgAssetCreate, transactionHash string) error {
	query := `INSERT INTO overgold_chain_assets_create 
			(transaction_hash, creator, name, issuer, policies, state, precision, fee_percent, extras) 
		VALUES 
			(:transaction_hash, :creator, :name, :issuer, :policies, :state, :precision, :fee_percent, :extras)`

	if _, err := r.db.NamedExec(query, toCreateAssetDatabase(msg, transactionHash)); err != nil {
		return errs.Internal{Cause: err.Error()}
	}

	return nil
}

// GetCreateAsset get the given wallet from database
func (r *Repository) GetCreateAsset(assetFilter filter.Filter) ([]*assetstypes.MsgAssetCreate, error) {
	query, args := assetFilter.Build(
		tableCreateAssets,
		types.FieldCreator, types.FieldName,
		types.FieldIssuer, types.FieldPolicies,
		types.FieldState, types.FieldPrecision,
		types.FieldFeePercent, types.FieldExtras,
	)

	var result []types.DBAssetCreate

	if err := r.db.Select(&result, query, args...); err != nil {
		return []*assetstypes.MsgAssetCreate{}, errs.Internal{Cause: err.Error()}
	}

	assets := make([]*assetstypes.MsgAssetCreate, 0, len(result))
	for _, asset := range result {
		assets = append(assets, toCreateAssetDomain(asset))
	}

	return assets, nil
}
