package types

import storagetypes "github.com/jackalLabs/canine-chain/v3/x/storage/types"

// storageParams represents the x/storage parameters
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
