package topaccounts

type BankModule interface {
	UpdateBalances(addresses []string, height int64) error
}

type DistrModule interface {
}

type StakingModule interface {
}
