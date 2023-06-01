package types

import (
	"github.com/KYVENetwork/chain/x/query/types"
	stakerstypes "github.com/KYVENetwork/chain/x/stakers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// StakersParams represents the x/stakers parameters
type StakersParams struct {
	stakerstypes.Params
	Height int64
}

// NewStakersParams allows to build a new StakersParams instance
func NewStakersParams(params stakerstypes.Params, height int64) *StakersParams {
	return &StakersParams{
		Params: params,
		Height: height,
	}
}

// ----------------------------------------------------------------------------------------------------------

// ProtocolValidator represents the x/stakers protocol validators
type ProtocolValidator struct {
	Address string
	Height  int64
}

// NewProtocolValidator allows to build a new ProtocolValidator instance
func NewProtocolValidator(
	address string, height int64,
) ProtocolValidator {
	return ProtocolValidator{
		Address: address,
		Height:  height,
	}
}

// ----------------------------------------------------------------------------------------------------------

// ProtocolValidatorCommission contains the data of a protocol
// validator commission at a given height
type ProtocolValidatorCommission struct {
	Address                 string
	Commission              sdk.Dec
	PendingCommissionChange *types.CommissionChangeEntry
	SelfDelegation          uint64
	Height                  int64
}

// NewProtocolValidatorCommission return a new ProtocolValidatorCommission instance
func NewProtocolValidatorCommission(
	address string, commission sdk.Dec, pendingCommissionChange *types.CommissionChangeEntry, selfDelegation uint64, height int64,
) ProtocolValidatorCommission {
	return ProtocolValidatorCommission{
		Address:                 address,
		Commission:              commission,
		PendingCommissionChange: pendingCommissionChange,
		SelfDelegation:          selfDelegation,
		Height:                  height,
	}
}

// --------------------------------------------------------------------------------------------------------------------

// ProtocolValidatorDelegation contains info of a protocol
// validator delegations
type ProtocolValidatorDelegation struct {
	Address         string
	SelfDelegation  uint64
	TotalDelegation uint64
	DelegatorCount  uint64
	Height          int64
}

// NewProtocolValidatorDelegation return a new ProtocolValidatorDelegation object
func NewProtocolValidatorDelegation(
	address string, selfDelegation, totalDelegation, delegatorCount uint64, height int64,
) ProtocolValidatorDelegation {
	return ProtocolValidatorDelegation{
		Address:         address,
		SelfDelegation:  selfDelegation,
		TotalDelegation: totalDelegation,
		DelegatorCount:  delegatorCount,
		Height:          height,
	}
}

// --------------------------------------------------------------------------------------------------------------------

// ProtocolValidatorDescription contains the description of a protocol validator
type ProtocolValidatorDescription struct {
	Address        string
	StakerMetadata *types.StakerMetadata
	AvatarURL      string // URL of the avatar to be used. Will be [do-no-modify] if it shouldn't be edited
	Height         int64
}

// NewProtocolValidatorDescription return a new ProtocolValidatorDescription object
func NewProtocolValidatorDescription(
	address string, stakerMetadata *types.StakerMetadata, avatarURL string, height int64,
) ProtocolValidatorDescription {
	return ProtocolValidatorDescription{
		Address:        address,
		StakerMetadata: stakerMetadata,
		AvatarURL:      avatarURL,
		Height:         height,
	}
}

// --------------------------------------------------------------------------------------------------------------------

// ProtocolValidatorPool contains the info of a protocol validator joined pools
type ProtocolValidatorPool struct {
	Address          string
	ValidatorAddress string
	Balance          uint64
	Pool             string
	Height           int64
}

// NewProtocolValidatorPool return a new ProtocolValidatorPool object
func NewProtocolValidatorPool(
	address, validatorAddress string, balance uint64, pool string, height int64,
) ProtocolValidatorPool {
	return ProtocolValidatorPool{
		Address:          address,
		ValidatorAddress: validatorAddress,
		Balance:          balance,
		Pool:             pool,
		Height:           height,
	}
}
