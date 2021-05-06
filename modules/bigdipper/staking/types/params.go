package types

// StakingParams contains all the parameters related to the staking module
type StakingParams struct {
	BondName string
}

// NewStakingParams allows to build a new StakingParams
func NewStakingParams(bondDenom string) StakingParams {
	return StakingParams{
		BondName: bondDenom,
	}
}
