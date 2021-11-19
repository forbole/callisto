package consensus

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/forbole/juno/v2/types"
)

type AuthModule interface {
	RefreshAccounts(height int64, addresses []string) error
}

type BankModule interface {
	HandleMsg(_ int, msg sdk.Msg, tx *types.Tx) error
}

type DistrModule interface {
	HandleMsg(_ int, msg sdk.Msg, tx *types.Tx) error
}

type GovModule interface {
	HandleMsg(_ int, msg sdk.Msg, tx *types.Tx) error
}

type StakingModule interface {
	HandleMsg(_ int, msg sdk.Msg, tx *types.Tx) error
}
