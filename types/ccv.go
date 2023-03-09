package types

import (
	"time"

	ibcclienttypes "github.com/cosmos/ibc-go/v4/modules/core/02-client/types"
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

// --------------------------------------------------------------------------------------------------------------------

// CcvConsumerChain represents ccv consumer chain state at a given height
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

// CcvProviderChain represents ccv provider chain state at a given height
type CcvProviderChain struct {
	ValsetUpdateID            uint64
	ConsumerStates            []ccvprovidertypes.ConsumerState
	UnbondingOps              []ccvprovidertypes.UnbondingOp
	MatureUnbondingOps        *ccvtypes.MaturedUnbondingOps
	ValsetUpdateIdToHeight    []ccvprovidertypes.ValsetUpdateIdToHeight
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
	valsetUpdateIdToHeight []ccvprovidertypes.ValsetUpdateIdToHeight,
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
		ValsetUpdateIdToHeight:    valsetUpdateIdToHeight,
		ConsumerAdditionProposals: consumerAdditionProposals,
		ConsumerRemovalProposals:  consumerRemovalProposals,
		ValidatorConsumerPubkeys:  validatorConsumerPubkeys,
		ValidatorsByConsumerAddr:  validatorsByConsumerAddr,
		ConsumerAddrsToPrune:      consumerAddrsToPrune,
		Height:                    height,
	}
}

// --------------------------------------------------------------------------------------------------------------------

// CcvProposalContent represents a single ccv proposal content
type CcvProposalContent struct {
	Title                             string
	Description                       string
	ChainID                           string
	GenesisHash                       string
	BinaryHash                        string
	SpawnTime                         time.Time
	StopTime                          time.Time
	InitialHeight                     ibcclienttypes.Height
	UnbondingPeriod                   time.Duration
	CcvTimeoutPeriod                  time.Duration
	TransferTimeoutPeriod             time.Duration
	ConsumerRedistributionFraction    string
	BlocksPerDistributionTransmission int64
	HistoricalEntries                 int64
}

// NewCcvProposalContent return a new CcvProposalContent instance
func NewCcvProposalContent(
	title string,
	description string,
	chainID string,
	genesisHash string,
	binaryHash string,
	spawnTime time.Time,
	stopTime time.Time,
	initialHeight ibcclienttypes.Height,
	unbondingPeriod time.Duration,
	ccvTimeoutPeriod time.Duration,
	transferTimeoutPeriod time.Duration,
	consumerRedistributionFraction string,
	blocksPerDistributionTransmission int64,
	historicalEntries int64,
) CcvProposalContent {
	return CcvProposalContent{
		Title:                             title,
		Description:                       description,
		ChainID:                           chainID,
		GenesisHash:                       genesisHash,
		BinaryHash:                        binaryHash,
		SpawnTime:                         spawnTime,
		StopTime:                          stopTime,
		InitialHeight:                     initialHeight,
		UnbondingPeriod:                   unbondingPeriod,
		CcvTimeoutPeriod:                  ccvTimeoutPeriod,
		TransferTimeoutPeriod:             transferTimeoutPeriod,
		ConsumerRedistributionFraction:    consumerRedistributionFraction,
		BlocksPerDistributionTransmission: blocksPerDistributionTransmission,
		HistoricalEntries:                 historicalEntries,
	}
}

// CcvProposal represents a single ccv proposal
type CcvProposal struct {
	ProposalRoute   string
	ProposalType    string
	ProposalID      uint64
	Content         CcvProposalContent
	Status          string
	SubmitTime      time.Time
	DepositEndTime  time.Time
	VotingStartTime time.Time
	VotingEndTime   time.Time
	Proposer        string
}

// NewCcvProposal return a new CcvProposal instance
func NewCcvProposal(
	proposalID uint64,
	proposalRoute string,
	proposalType string,
	content CcvProposalContent,
	status string,
	submitTime time.Time,
	depositEndTime time.Time,
	votingStartTime time.Time,
	votingEndTime time.Time,
	proposer string,
) CcvProposal {
	return CcvProposal{
		Content:         content,
		ProposalRoute:   proposalRoute,
		ProposalType:    proposalType,
		ProposalID:      proposalID,
		Status:          status,
		SubmitTime:      submitTime,
		DepositEndTime:  depositEndTime,
		VotingStartTime: votingStartTime,
		VotingEndTime:   votingEndTime,
		Proposer:        proposer,
	}
}
