package staking

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	liquidstakingtypes "github.com/crescent-network/crescent/v2/x/liquidstaking/types"
	"github.com/forbole/bdjuno/v3/types"
)

type SlashingModule interface {
	GetSigningInfo(height int64, consAddr sdk.ConsAddress) (types.ValidatorSigningInfo, error)
}

type LiquidStakingModule interface {
	GetLiquidValidators(height int64) ([]liquidstakingtypes.LiquidValidatorState, error)
}
