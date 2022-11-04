package upgrade

type StakingModule interface {
	RefreshAllValidatorInfos(height int64) error
}
