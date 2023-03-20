package top_accounts

type AuthModule interface {
	GetTotalNumberOfAccounts(height int64) (uint64, error) 
}
type BankModule interface {
	UpdateBalances(addresses []string, height int64) error
}

type DistrModule interface {
	RefreshDelegatorRewards(height int64, delegators []string) error
}

type StakingModule interface {
	RefreshDelegations(height int64, delegatorAddr string) error
	RefreshRedelegations(height int64, delegatorAddr string) error
	RefreshUnbondings(height int64, delegatorAddr string) error
}
