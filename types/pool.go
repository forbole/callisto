package types

import (
	pooltypes "github.com/KYVENetwork/chain/x/pool/types"
)

// PoolList represents the x/pool list
type PoolList struct {
	ID                       uint64
	Name                     string
	Runtime                  string
	Logo                     string
	Config                   string
	StartKey                 string
	CurrentKey               string
	CurrentSummary           string
	CurrentIndex             uint64
	TotalBundles             uint64
	UploadInterval           uint64
	OperatingCost            uint64
	MinDelegation            uint64
	MaxBundleSize            uint64
	Disabled                 bool
	Funders                  []*pooltypes.Funder
	TotalFunds               uint64
	Protocol                 *pooltypes.Protocol
	UpgradePlan              *pooltypes.UpgradePlan
	CurrentStorageProviderID uint32
	CurrentCompressionID     uint32
	Height                   int64
}

// NewPoolList allows to build a new PoolList instance
func NewPoolList(
	id uint64,
	name string,
	runtime string,
	logo string,
	config string,
	startKey string,
	currentKey string,
	currentSummary string,
	currentIndex uint64,
	totalBundles uint64,
	uploadInterval uint64,
	operatingCost uint64,
	minDelegation uint64,
	maxBundleSize uint64,
	disabled bool,
	funders []*pooltypes.Funder,
	totalFunds uint64,
	protocol *pooltypes.Protocol,
	upgradePlan *pooltypes.UpgradePlan,
	currentStorageProviderID uint32,
	currentCompressionID uint32,
	height int64) PoolList {
	return PoolList{
		ID:                       id,
		Name:                     name,
		Runtime:                  runtime,
		Logo:                     logo,
		Config:                   config,
		StartKey:                 startKey,
		CurrentKey:               currentKey,
		CurrentSummary:           currentSummary,
		CurrentIndex:             currentIndex,
		TotalBundles:             totalBundles,
		UploadInterval:           uploadInterval,
		OperatingCost:            operatingCost,
		MinDelegation:            minDelegation,
		MaxBundleSize:            maxBundleSize,
		Disabled:                 disabled,
		Funders:                  funders,
		TotalFunds:               totalFunds,
		Protocol:                 protocol,
		UpgradePlan:              upgradePlan,
		CurrentStorageProviderID: currentStorageProviderID,
		CurrentCompressionID:     currentCompressionID,
		Height:                   height,
	}
}
