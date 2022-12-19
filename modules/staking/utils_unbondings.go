package staking

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/forbole/bdjuno/v3/types"
	"github.com/rs/zerolog/log"
)

func (m *Module) RefreshUnbondings(height int64, index int, delegatorAddr string) error {
	log.Debug().
		Str("module", "staking").
		Int64("height", height).Msg("updating unbonding delegations")

	coin := sdk.NewInt(0)
	var nextKey []byte
	var stop = false
	for !stop {
		res, err := m.source.GetUnbondingDelegations(
			height,
			delegatorAddr,
			&query.PageRequest{
				Key:   nextKey,
				Limit: 100,
			},
		)
		if err != nil {
			return fmt.Errorf("error while getting delegations: %s", err)
		}

		nextKey = res.Pagination.NextKey
		stop = len(res.Pagination.NextKey) == 0

		for _, r := range res.UnbondingResponses {
			for _, e := range r.Entries {
				coin = coin.Add(e.Balance)
			}
		}
	}

	err := m.db.SaveTopAccountsBalance("unbonding",
		[]types.NativeTokenAmount{
			types.NewNativeTokenAmount(delegatorAddr, coin, height),
		})
	if err != nil {
		return fmt.Errorf("error while savting top accounts delegation from MsgDelegate: %s", err)
	}
	return nil
}
