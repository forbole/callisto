package slashing

import (
	"context"

	"github.com/cosmos/cosmos-sdk/types/query"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"

	"github.com/forbole/bdjuno/x/slashing/types"
	"github.com/forbole/bdjuno/x/utils"
)

func GetSigningInfos(height int64, client slashingtypes.QueryClient) ([]types.ValidatorSigningInfo, error) {
	var signingInfos []slashingtypes.ValidatorSigningInfo

	header := utils.GetHeightRequestHeader(height)

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
