package staking

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/forbole/bdjuno/v2/types"
)

type BankModule interface {
	RefreshBalances(height int64, addresses []string) error
}

type DistrModule interface {
	RefreshDelegatorRewards(height int64, delegator string) error
}

type HistoryModule interface {
	UpdateAccountBalanceHistoryWithTime(address string, time time.Time) error
}

type SlashingModule interface {
	GetSigningInfo(height int64, consAddr sdk.ConsAddress) (types.ValidatorSigningInfo, error)
}
