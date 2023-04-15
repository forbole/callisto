package distribution

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/forbole/bdjuno/v4/modules/pricefeed"
	"github.com/forbole/bdjuno/v4/types"
	"github.com/rs/zerolog/log"
)

func (m *Module) RefreshDelegatorRewards(delegators []string, height int64) error {
	log.Debug().
		Str("module", "distribution").
		Int64("height", height).Msg("updating rewards")

	nativeTokenAmounts := make([]types.NativeTokenAmount, len(delegators))
	for index, del := range delegators {
		rews, err := m.source.DelegatorTotalRewards(del, height)
		if err != nil {
			return fmt.Errorf("error while getting delegator rewards: %s", err)
		}

		amount := sdk.NewDec(0)
		for _, r := range rews {
			decCoinAmount := r.Reward.AmountOf(pricefeed.GetDenom())
			amount = amount.Add(decCoinAmount)
		}

		nativeTokenAmounts[index] = types.NewNativeTokenAmount(del, amount.RoundInt(), height)
	}

	err := m.db.SaveTopAccountsBalance("reward", nativeTokenAmounts)
	if err != nil {
		return fmt.Errorf("error while saving delegators rewards amounts: %s", err)
	}

	return nil
}
