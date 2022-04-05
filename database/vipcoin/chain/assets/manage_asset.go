package assets

import (
	assetstypes "git.ooo.ua/vipcoin/chain/x/assets/types"
	"git.ooo.ua/vipcoin/lib/filter"

	"github.com/forbole/bdjuno/v2/database/types"
)

// SaveManageAsset saves the given asset inside the database
func (r *Repository) SaveManageAsset(msg ...*assetstypes.MsgAssetManage) error {
	if len(msg) == 0 {
		return nil
	}

	query := `INSERT INTO vipcoin_chain_assets_manage 
			(creator, name, policies, state, precision, fee_percent, issued, burned, withdrawn, in_circulation) 
		VALUES 
			(:creator, :name, :policies, :state, :precision, :fee_percent, :issued, :burned, :withdrawn, :in_circulation)`

	if _, err := r.db.NamedExec(query, toManageAssetsArrDatabase(msg...)); err != nil {
		return err
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
		return []*assetstypes.MsgAssetManage{}, err
	}

	assets := make([]*assetstypes.MsgAssetManage, 0, len(result))
	for _, asset := range result {
		assets = append(assets, toManageAssetDomain(asset))
	}

	return assets, nil
}
