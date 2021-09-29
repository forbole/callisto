package slashing

import (
	"github.com/forbole/bdjuno/types"
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
