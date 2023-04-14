package staking

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/forbole/bdjuno/v4/modules/pricefeed"
	"github.com/forbole/bdjuno/v4/types"
	"github.com/rs/zerolog/log"
)

func (m *Module) RefreshDelegations(delegatorAddr string, height int64) error {
	log.Debug().
		Str("module", "staking").
		Int64("height", height).Msg("updating delegations")

	var coin = sdk.Coin{
		Denom:  pricefeed.GetDenom(),
		Amount: sdk.NewInt(0),
	}
	var nextKey []byte
	var stop = false
	for !stop {
		res, err := m.source.GetDelegationsWithPagination(
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

		for _, r := range res.DelegationResponses {
			coin = coin.Add(r.Balance)
		}
	}

	err := m.db.SaveTopAccountsBalance("delegation",
		[]types.NativeTokenAmount{
			types.NewNativeTokenAmount(delegatorAddr, coin.Amount, height),
		})
	if err != nil {
		return fmt.Errorf("error while saving top accounts delegation from MsgDelegate: %s", err)
	}
	return nil
}
