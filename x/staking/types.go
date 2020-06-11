package staking

import (
	"github.com/cosmos/cosmos-sdk/x/staking"
)

var (
	Fetcher = DataFetcher
	Handler = DataHandler
)

// ValidatorInfo represents the information of a validator
type ValidatorInfo struct {
	*staking.Validator
	Delegations *staking.Delegations
}

// NewValidatorInfo allows to create a new ValidatorInfo object
func NewValidatorInfo(info *staking.Validator, delegations *staking.Delegations) ValidatorInfo {
	return ValidatorInfo{
		Validator:   info,
		Delegations: delegations,
	}
}
