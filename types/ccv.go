package types

import (
	ccvconsumertypes "github.com/cosmos/interchain-security/x/ccv/consumer/types"
	ccvprovidertypes "github.com/cosmos/interchain-security/x/ccv/provider/types"
)

// CcvProviderParams represents the parameters of the ccv provider module at a given height
type CcvProviderParams struct {
	ccvprovidertypes.Params
	Height int64
}

// NewCcvProviderParams allows to build a new CcvProviderParams instance
func NewCcvProviderParams(params ccvprovidertypes.Params, height int64) *CcvProviderParams {
	return &CcvProviderParams{
		Params: params,
		Height: height,
	}
}

// CcvConsumerParams represents the parameters of the ccv consumer module at a given height
type CcvConsumerParams struct {
	ccvconsumertypes.Params
	Height int64
}

// NewCcvConsumerParams allows to build a new CcvConsumerParams instance
func NewCcvConsumerParams(params ccvconsumertypes.Params, height int64) *CcvConsumerParams {
	return &CcvConsumerParams{
		Params: params,
		Height: height,
	}
}
