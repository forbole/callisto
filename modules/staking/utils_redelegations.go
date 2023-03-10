package staking

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/forbole/bdjuno/v4/types"
	"github.com/rs/zerolog/log"
)

func (m *Module) RefreshRedelegations(height int64, delegatorAddr string) error {
	log.Debug().
		Str("module", "staking").
		Int64("height", height).Msg("updating redelegations")

	coin := sdk.NewInt(0)
	var nextKey []byte
	var stop = false
	for !stop {
		res, err := m.source.GetRedelegations(
			height,
			&stakingtypes.QueryRedelegationsRequest{
				DelegatorAddr: delegatorAddr,
				Pagination: &query.PageRequest{
					Key:   nextKey,
					Limit: 100,
				},
			},
		)
		if err != nil {
			return fmt.Errorf("error while getting redelegations: %s", err)
		}

		nextKey = res.Pagination.NextKey
		stop = len(res.Pagination.NextKey) == 0

		for _, r := range res.RedelegationResponses {
			for _, e := range r.Entries {
				coin = coin.Add(e.Balance)
			}
		}
	}

	err := m.db.SaveTopAccountsBalance("redelegation",
		[]types.NativeTokenAmount{
			types.NewNativeTokenAmount(delegatorAddr, coin, height),
		})
	if err != nil {
		return fmt.Errorf("error while saving top accounts redelegation from MsgBeginRedelegate: %s", err)
	}
	return nil
}
