package staking

import "time"

type BankModule interface {
	RefreshBalances(height int64, addresses []string) error
}

type DistrModule interface {
	RefreshDelegatorRewards(height int64, delegator string) error
}

type HistoryModule interface {
	UpdateAccountBalanceHistoryWithTime(address string, time time.Time) error
}
