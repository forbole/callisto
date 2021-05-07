package distribution

import "github.com/forbole/bdjuno/types"

// DB represents a generic database that allows to perform x/distribution related operations
type DB interface {
	GetDelegators() ([]string, error)
	GetValidatorsInfo() ([]types.ValidatorInfo, error)
	SaveDelegatorsRewardsAmounts([]types.DelegatorRewardAmount) error
	SaveValidatorCommissionAmount(amount types.ValidatorCommissionAmount) error
}
