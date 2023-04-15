package top_accounts

import (
	"github.com/forbole/bdjuno/v4/types"
)

type AuthModule interface {
	GetAllBaseAccounts(height int64) ([]types.Account, error)
}

type AuthSource interface {
	GetTotalNumberOfAccounts(height int64) (uint64, error)
}

type BankModule interface {
	UpdateBalances(addresses []string, height int64) error
}

type DistrModule interface {
	RefreshDelegatorRewards(delegators []string, height int64) error
}

type StakingModule interface {
	RefreshDelegations(delegatorAddr string, height int64) error
	RefreshRedelegations(delegatorAddr string, height int64) error
	RefreshUnbondings(delegatorAddr string, height int64) error
}
