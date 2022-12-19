package staking

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/forbole/bdjuno/v3/types"
	juno "github.com/forbole/juno/v3/types"
	"github.com/rs/zerolog/log"
)

func (m *Module) RefreshRedelegations(tx *juno.Tx, index int, delegatorAddr string) error {
	log.Debug().
		Str("module", "staking").
		Str("delegator", delegatorAddr).
		Int64("height", tx.Height).Msg("updating redelegations")

	coin := sdk.NewInt(0)
	var nextKey []byte
	var stop = false
	for !stop {
		res, err := m.source.GetRedelegations(
			tx.Height,
			&stakingtypes.QueryRedelegationsRequest{
				DelegatorAddr: delegatorAddr,
				Pagination: &query.PageRequest{
					Key:   nextKey,
					Limit: 100,
				},
			},
		)
		if err != nil {
			return fmt.Errorf("error while getting delegations: %s", err)
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
			types.NewNativeTokenAmount(delegatorAddr, coin, tx.Height),
		})
	if err != nil {
		return fmt.Errorf("error while savting top accounts delegation from MsgDelegate: %s", err)
	}
	return nil
}
