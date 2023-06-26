package slashing

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/forbole/bdjuno/v4/types"
)

func (m *Module) getSigningInfos(height int64) ([]types.ValidatorSigningInfo, error) {
	signingInfos, err := m.source.GetSigningInfos(height)
	if err != nil {
		return nil, err
	}

	infos := make([]types.ValidatorSigningInfo, len(signingInfos))
	for index, info := range signingInfos {
		infos[index] = types.NewValidatorSigningInfo(
			info.Address,
			info.StartHeight,
			info.IndexOffset,
			info.JailedUntil,
			info.Tombstoned,
			info.MissedBlocksCounter,
			height,
		)
	}

	return infos, nil
}

// GetSigningInfo returns the signing info for the validator having the given consensus address at the specified height
func (m *Module) GetSigningInfo(height int64, consAddr sdk.ConsAddress) (types.ValidatorSigningInfo, error) {
	info, err := m.source.GetSigningInfo(height, consAddr)
	if err != nil {
		return types.ValidatorSigningInfo{}, err
	}

	signingInfo := types.NewValidatorSigningInfo(
		info.Address,
		info.StartHeight,
		info.IndexOffset,
		info.JailedUntil,
		info.Tombstoned,
		info.MissedBlocksCounter,
		height,
	)

	return signingInfo, nil
}
