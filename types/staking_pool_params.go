package types

import (
	sdkmath "cosmossdk.io/math"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// Pool contains the data of the staking pool at the given height
type Pool struct {
	BondedTokens          sdkmath.Int
	NotBondedTokens       sdkmath.Int
	UnbondingTokens       sdkmath.Int
	StakedNotBondedTokens sdkmath.Int
	Height                int64
}

// NewPool allows to build a new Pool instance
func NewPool(bondedTokens, notBondedTokens, unbondingTokens, stakedNotBondedTokens sdkmath.Int, height int64) *Pool {
	return &Pool{
		BondedTokens:          bondedTokens,
		NotBondedTokens:       notBondedTokens,
		UnbondingTokens:       unbondingTokens,
		StakedNotBondedTokens: stakedNotBondedTokens,
		Height:                height,
	}
}

// PoolSnapshot contains the data of the staking pool snapshot at the given height
type PoolSnapshot struct {
	BondedTokens    sdkmath.Int
	NotBondedTokens sdkmath.Int
	Height          int64
}

// NewPoolSnapshot allows to build a new PoolSnapshot instance
func NewPoolSnapshot(bondedTokens, notBondedTokens sdkmath.Int, height int64) *PoolSnapshot {
	return &PoolSnapshot{
		BondedTokens:    bondedTokens,
		NotBondedTokens: notBondedTokens,
		Height:          height,
	}
}

// --------------------------------------------------------------------------------------------------------------------

// StakingParams represents the parameters of the x/staking module
type StakingParams struct {
	stakingtypes.Params
	Height int64
}

// NewStakingParams returns a new StakingParams instance
func NewStakingParams(params stakingtypes.Params, height int64) StakingParams {
	return StakingParams{
		Params: params,
		Height: height,
	}
}
