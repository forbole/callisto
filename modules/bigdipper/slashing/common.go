package slashing

import (
	"context"

	"github.com/cosmos/cosmos-sdk/types/query"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"

	bslashingtypes "github.com/forbole/bdjuno/modules/bigdipper/slashing/types"
	utils2 "github.com/forbole/bdjuno/modules/common/utils"
)

func GetSigningInfos(height int64, client slashingtypes.QueryClient) ([]bslashingtypes.ValidatorSigningInfo, error) {
	var signingInfos []slashingtypes.ValidatorSigningInfo

	header := utils2.GetHeightRequestHeader(height)

	var nextKey []byte
	var stop = false
	for !stop {
		res, err := client.SigningInfos(
			context.Background(),
			&slashingtypes.QuerySigningInfosRequest{
				Pagination: &query.PageRequest{
					Key:   nextKey,
					Limit: 1000, // Query 1000 signing infos at a time
				},
			},
			header,
		)
		if err != nil {
			return nil, err
		}

		nextKey = res.Pagination.NextKey
		stop = len(res.Pagination.NextKey) == 0
		signingInfos = append(signingInfos, res.Info...)
	}

	infos := make([]bslashingtypes.ValidatorSigningInfo, len(signingInfos))
	for index, info := range signingInfos {
		infos[index] = bslashingtypes.NewValidatorSigningInfo(
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
