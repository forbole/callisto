package types

import (
	"github.com/cosmos/ibc-go/v4/modules/light-clients/07-tendermint/types"
	ccvconsumertypes "github.com/cosmos/interchain-security/x/ccv/consumer/types"
	ccvprovidertypes "github.com/cosmos/interchain-security/x/ccv/provider/types"
	ccvtypes "github.com/cosmos/interchain-security/x/ccv/types"
	abcitypes "github.com/tendermint/tendermint/abci/types"
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

// CcvConsumerChain represents ccv consumer chain at a given height
type CcvConsumerChain struct {
	ProviderClientID            string
	ProviderChannelID           string
	NewChain                    bool
	ProviderClientState         *types.ClientState
	ProviderConsensusState      *types.ConsensusState
	MaturingPackets             []ccvconsumertypes.MaturingVSCPacket
	InitialValSet               []abcitypes.ValidatorUpdate
	HeightToValsetUpdateID      []ccvconsumertypes.HeightToValsetUpdateID
	OutstandingDowntimeSlashing []ccvconsumertypes.OutstandingDowntime
	PendingConsumerPackets      ccvtypes.ConsumerPacketDataList
	LastTransmissionBlockHeight ccvconsumertypes.LastTransmissionBlockHeight
	Height                      int64
}

// NewCcvConsumerChain allows to build a new CcvConsumerChain instance
func NewCcvConsumerChain(providerClientID, providerChannelID string, newChain bool,
	providerClientState *types.ClientState, providerConsensusState *types.ConsensusState,
	maturingPackets []ccvconsumertypes.MaturingVSCPacket, initialValSet []abcitypes.ValidatorUpdate,
	heightToValsetUpdateID []ccvconsumertypes.HeightToValsetUpdateID,
	outstandingDowntimeSlashing []ccvconsumertypes.OutstandingDowntime,
	pendingConsumerPackets ccvtypes.ConsumerPacketDataList,
	lastTransmissionBlockHeight ccvconsumertypes.LastTransmissionBlockHeight,
	height int64) *CcvConsumerChain {
	return &CcvConsumerChain{
		ProviderClientID:            providerClientID,
		ProviderChannelID:           providerChannelID,
		NewChain:                    newChain,
		ProviderClientState:         providerClientState,
		ProviderConsensusState:      providerConsensusState,
		MaturingPackets:             maturingPackets,
		InitialValSet:               initialValSet,
		HeightToValsetUpdateID:      heightToValsetUpdateID,
		OutstandingDowntimeSlashing: outstandingDowntimeSlashing,
		PendingConsumerPackets:      pendingConsumerPackets,
		LastTransmissionBlockHeight: lastTransmissionBlockHeight,
		Height:                      height,
	}
}
