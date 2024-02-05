package types

import (
	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/ibc-go/v4/modules/light-clients/07-tendermint/types"
	ccvconsumertypes "github.com/cosmos/interchain-security/v2/x/ccv/consumer/types"
	ccvprovidertypes "github.com/cosmos/interchain-security/v2/x/ccv/provider/types"
	ccvtypes "github.com/cosmos/interchain-security/v2/x/ccv/types"
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

// CcvProviderChain represents ccv provider chain state at a given height
type CcvProviderChain struct {
	ValsetUpdateID            uint64
	ConsumerStates            []ccvprovidertypes.ConsumerState
	UnbondingOps              []ccvprovidertypes.UnbondingOp
	MatureUnbondingOps        *ccvtypes.MaturedUnbondingOps
	ValsetUpdateIDToHeight    []ccvprovidertypes.ValsetUpdateIdToHeight
	ConsumerAdditionProposals []ccvprovidertypes.ConsumerAdditionProposal
	ConsumerRemovalProposals  []ccvprovidertypes.ConsumerRemovalProposal
	ValidatorConsumerPubkeys  []ccvprovidertypes.ValidatorConsumerPubKey
	ValidatorsByConsumerAddr  []ccvprovidertypes.ValidatorByConsumerAddr
	ConsumerAddrsToPrune      []ccvprovidertypes.ConsumerAddrsToPrune
	Height                    int64
}

// NewNewCcvProviderChain allows to build a new CcvProviderChain instance
func NewCcvProviderChain(valsetUpdateID uint64, consumerStates []ccvprovidertypes.ConsumerState,
	unbondingOps []ccvprovidertypes.UnbondingOp, matureUnbondingOps *ccvtypes.MaturedUnbondingOps,
	valsetUpdateIDToHeight []ccvprovidertypes.ValsetUpdateIdToHeight,
	consumerAdditionProposals []ccvprovidertypes.ConsumerAdditionProposal,
	consumerRemovalProposals []ccvprovidertypes.ConsumerRemovalProposal,
	validatorConsumerPubkeys []ccvprovidertypes.ValidatorConsumerPubKey,
	validatorsByConsumerAddr []ccvprovidertypes.ValidatorByConsumerAddr,
	consumerAddrsToPrune []ccvprovidertypes.ConsumerAddrsToPrune,
	height int64) *CcvProviderChain {
	return &CcvProviderChain{
		ValsetUpdateID:            valsetUpdateID,
		ConsumerStates:            consumerStates,
		UnbondingOps:              unbondingOps,
		MatureUnbondingOps:        matureUnbondingOps,
		ValsetUpdateIDToHeight:    valsetUpdateIDToHeight,
		ConsumerAdditionProposals: consumerAdditionProposals,
		ConsumerRemovalProposals:  consumerRemovalProposals,
		ValidatorConsumerPubkeys:  validatorConsumerPubkeys,
		ValidatorsByConsumerAddr:  validatorsByConsumerAddr,
		ConsumerAddrsToPrune:      consumerAddrsToPrune,
		Height:                    height,
	}
}
