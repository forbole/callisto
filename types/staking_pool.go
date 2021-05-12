package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// Pool contains the data of the staking pool at the given height
type Pool struct {
	BondedTokens    sdk.Int
	NotBondedTokens sdk.Int
	Height          int64
}

// NewPool allows to build a new Pool instance
func NewPool(bondedTokens, notBondedTokens sdk.Int, height int64) Pool {
	return Pool{
		BondedTokens:    bondedTokens,
		NotBondedTokens: notBondedTokens,
		Height:          height,
	}
}
