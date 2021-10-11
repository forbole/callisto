package slashing

type StakingModule interface {
	RefreshValidatorDelegations(height int64, valOperAddr string) error
}
