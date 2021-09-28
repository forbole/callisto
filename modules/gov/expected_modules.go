package gov

import (
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/forbole/bdjuno/types"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

type AuthModule interface {
	RefreshAccounts(height int64, addresses []string) error
}

type BankModule interface {
	RefreshBalances(height int64, addresses []string) error
}

type StakingModule interface {
	GetStakingPool(height int64) (*types.Pool, error)
	GetValidatorsWithStatus(height int64, status string) ([]stakingtypes.Validator, []types.Validator, error)
	GetValidatorsVotingPowers(height int64, vals *tmctypes.ResultValidators) []types.ValidatorVotingPower
	GetValidatorsStatuses(height int64, validators []stakingtypes.Validator) ([]types.ValidatorStatus, error)
}
