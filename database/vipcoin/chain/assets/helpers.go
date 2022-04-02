package assets

import (
	assetstypes "git.ooo.ua/vipcoin/chain/x/assets/types"
	extratypes "git.ooo.ua/vipcoin/chain/x/types"

	"github.com/forbole/bdjuno/v2/database/types"
	"github.com/lib/pq"
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
