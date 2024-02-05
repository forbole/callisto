package types

import storagetypes "github.com/jackalLabs/canine-chain/v3/x/storage/types"

// StorageParams represents the x/storage parameters
type StorageParams struct {
	storagetypes.Params
	Height int64
}

// NewStorageParams allows to build a new StorageParams instance
func NewStorageParams(params storagetypes.Params, height int64) *StorageParams {
	return &StorageParams{
		Params: params,
		Height: height,
	}
}

// StorageProvidersList represents the x/storage providers
type StorageProvidersList struct {
	Address         string
	IP              string
	TotalSpace      string
	BurnedContracts string
	Creator         string
	KeybaseIdentity string
	AuthClaimers    []string
	Height          int64
}

// NewStorageProvidersList allows to build a new StorageProvidersList instance
func NewStorageProvidersList(
	address string,
	ip string,
	totalspace string,
	burnedContracts string,
	creator string,
	keybaseIdentity string,
	authClaimers []string,
	height int64) *StorageProvidersList {
	return &StorageProvidersList{
		Address:         address,
		IP:              ip,
		TotalSpace:      totalspace,
		BurnedContracts: burnedContracts,
		Creator:         creator,
		KeybaseIdentity: keybaseIdentity,
		AuthClaimers:    authClaimers,
		Height:          height,
	}
}
