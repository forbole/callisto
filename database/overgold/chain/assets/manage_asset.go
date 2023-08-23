package assets

import (
	assetstypes "git.ooo.ua/vipcoin/chain/x/assets/types"
	"git.ooo.ua/vipcoin/lib/errs"
	"git.ooo.ua/vipcoin/lib/filter"

	"github.com/forbole/bdjuno/v3/database/types"
)

// SaveManageAsset saves the given asset inside the database
func (r *Repository) SaveManageAsset(msg *assetstypes.MsgAssetManage, transactionHash string) error {
	query := `INSERT INTO overgold_chain_assets_manage 
			(transaction_hash, creator, name, policies, state, precision, fee_percent, issued, burned, withdrawn, in_circulation) 
		VALUES 
			(:transaction_hash, :creator, :name, :policies, :state, :precision, :fee_percent, :issued, :burned, :withdrawn, :in_circulation)`

	if _, err := r.db.NamedExec(query, toManageAssetDatabase(msg, transactionHash)); err != nil {
		return errs.Internal{Cause: err.Error()}
	}

	return nil
}

// GetManageAsset get the given wallet from database
func (r *Repository) GetManageAsset(assetFilter filter.Filter) ([]*assetstypes.MsgAssetManage, error) {
	query, args := assetFilter.Build(
		tableManageAsset,
		types.FieldCreator, types.FieldName,
		types.FieldPolicies, types.FieldState,
		types.FieldPrecision, types.FieldFeePercent,
		types.FieldIssued, types.FieldBurned,
		types.FieldWithdrawn, types.FieldInCirculation,
	)

	var result []types.DBAssetManage

	if err := r.db.Select(&result, query, args...); err != nil {
		return []*assetstypes.MsgAssetManage{}, errs.Internal{Cause: err.Error()}
	}

	assets := make([]*assetstypes.MsgAssetManage, 0, len(result))
	for _, asset := range result {
		assets = append(assets, toManageAssetDomain(asset))
	}

	return assets, nil
}
