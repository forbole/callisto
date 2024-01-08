package types

import (
	abcitypes "github.com/cometbft/cometbft/abci/types"
	types "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"
	ccvconsumertypes "github.com/cosmos/interchain-security/v3/x/ccv/consumer/types"
)

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

// --------------------------------------------------------------------------------------------------------------------

// CcvConsumerChain represents ccv consumer chain state at a given height
type CcvConsumerChain struct {
	ProviderClientID       string
	ProviderChannelID      string
	ChainID                string
	ProviderClientState    *types.ClientState
	ProviderConsensusState *types.ConsensusState
	InitialValSet          []abcitypes.ValidatorUpdate
	Height                 int64
}

// NewCcvConsumerChain allows to build a new CcvConsumerChain instance
func NewCcvConsumerChain(providerClientID, providerChannelID string, chainID string,
	providerClientState *types.ClientState, providerConsensusState *types.ConsensusState,
	initialValSet []abcitypes.ValidatorUpdate, height int64) *CcvConsumerChain {
	return &CcvConsumerChain{
		ProviderClientID:       providerClientID,
		ProviderChannelID:      providerChannelID,
		ChainID:                chainID,
		ProviderClientState:    providerClientState,
		ProviderConsensusState: providerConsensusState,
		InitialValSet:          initialValSet,
		Height:                 height,
	}
}

// --------------------------------------------------------------------------------------------------------------------

type CCVValidator struct {
	ConsumerConsensusAddress    string
	ConsumerSelfDelegateAddress string
	ConsumerOperatorAddress     string
	ProviderConsensusAddress    string
	ProviderSelfDelegateAddress string
	ProviderOperatorAddress     string
	Height                      int64
}

func NewCCVValidator(consumerConsensusAddress, consumerSelfDelegateAddress, consumerOperatorAddress,
	providerConsensusAddress, providerSelfDelegateAddress, providerOperatorAddress string, height int64) CCVValidator {
	return CCVValidator{
		ConsumerConsensusAddress:    consumerConsensusAddress,
		ConsumerSelfDelegateAddress: consumerSelfDelegateAddress,
		ConsumerOperatorAddress:     consumerOperatorAddress,
		ProviderConsensusAddress:    providerConsensusAddress,
		ProviderSelfDelegateAddress: providerSelfDelegateAddress,
		ProviderOperatorAddress:     providerOperatorAddress,
		Height:                      height,
	}
}
