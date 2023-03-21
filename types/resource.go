package types

import (
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
)

// Resource represents the x/resource
type Resource struct {
	ID           string
	CollectionId string
	Data         []byte
	Name         string
	Version      string
	ResourceType string
	AlsoKnownAs  []*resourcetypes.AlternativeUri
	FromAddress  string
	Height       int64
}

// NewResource allows to build a new Resource instance
func NewResource(id string,
	collectionId string,
	data []byte,
	name string,
	version string,
	resourceType string,
	alsoKnownAs []*resourcetypes.AlternativeUri,
	fromAddress string,
	height int64) *Resource {
	return &Resource{
		ID:           id,
		CollectionId: collectionId,
		Data:         data,
		Name:         name,
		Version:      version,
		ResourceType: resourceType,
		AlsoKnownAs:  alsoKnownAs,
		FromAddress:  fromAddress,
		Height:       height,
	}
}
