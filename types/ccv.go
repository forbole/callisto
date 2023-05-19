package types

import (
	"github.com/cosmos/ibc-go/v4/modules/light-clients/07-tendermint/types"
	ccvconsumertypes "github.com/cosmos/interchain-security/x/ccv/consumer/types"
	abcitypes "github.com/tendermint/tendermint/abci/types"
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
	ConsumerConsensusAddress string
	ProviderConsensusAddress string
	Height                   int64
}

func NewCCVValidator(consumerConsensusAddress, providerConsensusAddress string, height int64) CCVValidator {
	return CCVValidator{
		ConsumerConsensusAddress: consumerConsensusAddress,
		ProviderConsensusAddress: providerConsensusAddress,
		Height:                   height,
	}
}
