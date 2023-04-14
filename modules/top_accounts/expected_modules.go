package top_accounts

type AuthModule interface {
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
