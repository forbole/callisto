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
	InflationShareWeight     uint64
	MinDelegation            uint64
	MaxBundleSize            uint64
	Disabled                 bool
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
	inflationShareWeight uint64,
	minDelegation uint64,
	maxBundleSize uint64,
	disabled bool,
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
		InflationShareWeight:     inflationShareWeight,
		MinDelegation:            minDelegation,
		MaxBundleSize:            maxBundleSize,
		Disabled:                 disabled,
		Protocol:                 protocol,
		UpgradePlan:              upgradePlan,
		CurrentStorageProviderID: currentStorageProviderID,
		CurrentCompressionID:     currentCompressionID,
		Height:                   height,
	}
}
