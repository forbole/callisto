package staking

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/forbole/bdjuno/v3/types"
)

type SlashingModule interface {
	GetSigningInfo(height int64, consAddr sdk.ConsAddress) (types.ValidatorSigningInfo, error)
}
