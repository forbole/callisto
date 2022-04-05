package assets

import (
	assetstypes "git.ooo.ua/vipcoin/chain/x/assets/types"
	extratypes "git.ooo.ua/vipcoin/chain/x/types"

	"github.com/forbole/bdjuno/v2/database/types"
	"github.com/lib/pq"
)

const (
	tableAssets       = "vipcoin_chain_assets_assets"
	tableCreateAssets = "vipcoin_chain_assets_create"
	tableManageAsset  = "vipcoin_chain_assets_manage"
)

// toExtrasDB - mapping func to database model
func toExtrasDB(extras []*extratypes.Extra) types.ExtraDB {
	result := make([]extratypes.Extra, 0, len(extras))
	for _, extra := range extras {
		result = append(result, *extra)
	}

	return types.ExtraDB{Extras: result}
}

// fromExtrasDB - mapping func from database model
func fromExtrasDB(extras types.ExtraDB) []*extratypes.Extra {
	result := make([]*extratypes.Extra, 0, len(extras.Extras))
	for i := range extras.Extras {
		result = append(result, &extras.Extras[i])
	}

	return result
}

// toPoliciesDB - mapping func to database model
func toPoliciesDB(policies []assetstypes.AssetPolicy) pq.Int32Array {
	result := make(pq.Int32Array, 0, len(policies))
	for _, policy := range policies {
		result = append(result, assetstypes.AssetPolicy_value[policy.String()])
	}

	return result
}

// toPoliciesDomain - mapping func to domain model
func toPoliciesDomain(policies pq.Int32Array) []assetstypes.AssetPolicy {
	result := make([]assetstypes.AssetPolicy, 0, len(policies))
	for _, policy := range policies {
		result = append(result, assetstypes.AssetPolicy(policy))
	}

	return result
}

// toAssetDatabase - mapping func to database model
func toAssetDatabase(assets *assetstypes.Asset) types.DBAssets {
	return types.DBAssets{
		Issuer:        assets.Issuer,
		Name:          assets.Name,
		Policies:      toPoliciesDB(assets.Policies),
		State:         int32(assets.State),
		Issued:        assets.Issued,
		Burned:        assets.Burned,
		Withdrawn:     assets.Withdrawn,
		InCirculation: assets.InCirculation,
		Precision:     assets.Properties.Precision,
		FeePercent:    assets.Properties.FeePercent,
		Extras:        toExtrasDB(assets.Extras),
	}
}

// toAssetsArrDatabase- mapping func to database model
func toAssetsArrDatabase(assets ...*assetstypes.Asset) []types.DBAssets {
	result := make([]types.DBAssets, 0, len(assets))
	for _, asset := range assets {
		result = append(result, toAssetDatabase(asset))
	}

	return result
}

// toCreateAssetsArrDatabase - mapping func to database model
func toCreateAssetsArrDatabase(msgs ...*assetstypes.MsgAssetCreate) []types.DBAssetCreate {
	result := make([]types.DBAssetCreate, 0, len(msgs))
	for _, msg := range msgs {
		result = append(result, toCreateAssetDatabase(msg))
	}

	return result
}

// toCreateAssetDatabase - mapping func to database model
func toCreateAssetDatabase(msg *assetstypes.MsgAssetCreate) types.DBAssetCreate {
	return types.DBAssetCreate{
		Creator:    msg.Creator,
		Name:       msg.Name,
		Issuer:     msg.Issuer,
		Policies:   toPoliciesDB(msg.Policies),
		State:      int32(msg.State),
		Precision:  msg.Properties.Precision,
		FeePercent: msg.Properties.FeePercent,
		Extras:     toExtrasDB(msg.Extras),
	}
}

// toCreateAssetDomain - mapping func from database model
func toCreateAssetDomain(asset types.DBAssetCreate) *assetstypes.MsgAssetCreate {
	return &assetstypes.MsgAssetCreate{
		Creator:  asset.Creator,
		Name:     asset.Name,
		Issuer:   asset.Issuer,
		Policies: toPoliciesDomain(asset.Policies),
		State:    assetstypes.AssetState(asset.State),
		Properties: assetstypes.Properties{
			Precision:  asset.Precision,
			FeePercent: asset.FeePercent,
		},
		Extras: fromExtrasDB(asset.Extras),
	}
}

// toManageAssetsArrDatabase - mapping func to database model
func toManageAssetsArrDatabase(msgs ...*assetstypes.MsgAssetManage) []types.DBAssetManage {
	result := make([]types.DBAssetManage, 0, len(msgs))
	for _, msg := range msgs {
		result = append(result, toManageAssetDatabase(msg))
	}

	return result
}

// toManageAssetDatabase - mapping func to database model
func toManageAssetDatabase(msg *assetstypes.MsgAssetManage) types.DBAssetManage {
	return types.DBAssetManage{
		Creator:       msg.Creator,
		Name:          msg.Name,
		Policies:      toPoliciesDB(msg.Policies),
		State:         int32(msg.State),
		Issued:        msg.Issued,
		Burned:        msg.Burned,
		Withdrawn:     msg.Withdrawn,
		InCirculation: msg.InCirculation,
		Precision:     msg.Properties.Precision,
		FeePercent:    msg.Properties.FeePercent,
	}
}

// toManageAssetDomain - mapping func from database model
func toManageAssetDomain(asset types.DBAssetManage) *assetstypes.MsgAssetManage {
	return &assetstypes.MsgAssetManage{
		Creator:  asset.Creator,
		Name:     asset.Name,
		Policies: toPoliciesDomain(asset.Policies),
		State:    assetstypes.AssetState(asset.State),
		Properties: assetstypes.Properties{
			Precision:  asset.Precision,
			FeePercent: asset.FeePercent,
		},
		Issued:        asset.Issued,
		Burned:        asset.Burned,
		Withdrawn:     asset.Withdrawn,
		InCirculation: asset.InCirculation,
	}
}

// toAssetDomain - mapping func to domain model
func toAssetDomain(asset types.DBAssets) *assetstypes.Asset {
	return &assetstypes.Asset{
		Issuer:        asset.Issuer,
		Name:          asset.Name,
		Policies:      toPoliciesDomain(asset.Policies),
		State:         assetstypes.AssetState(asset.State),
		Issued:        asset.Issued,
		Burned:        asset.Burned,
		Withdrawn:     asset.Withdrawn,
		InCirculation: asset.InCirculation,
		Properties: assetstypes.Properties{
			Precision:  asset.Precision,
			FeePercent: asset.FeePercent,
		},
		Extras: fromExtrasDB(asset.Extras),
	}
}
