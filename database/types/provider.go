package types

import (
	"database/sql/driver"
	"fmt"

	providertypes "github.com/akash-network/node/x/provider/types/v1beta2"
	"github.com/forbole/bdjuno/v3/types"
)

// ProviderRow represents a single row inside the provider table
type ProviderRow struct {
	OwnerAddress string         `db:"owner_address"`
	HostURI      string         `db:"host_uri"`
	Attributes   string         `db:"attributes"`
	Info         DbProviderInfo `db:"info"`
	Height       int64          `db:"height"`
}

// NewProviderRow allows to build a new ProviderRow instance
func NewProviderRow(ownerAddress string, hostURI string, attributes string, info DbProviderInfo, height int64) ProviderRow {
	return ProviderRow{
		OwnerAddress: ownerAddress,
		HostURI:      hostURI,
		Attributes:   attributes,
		Info:         info,
		Height:       height,
	}
}

// Equal tells whether a and b contain the same data
func (a ProviderRow) Equal(b ProviderRow) bool {
	return a.OwnerAddress == b.OwnerAddress &&
		a.HostURI == b.HostURI &&
		a.Attributes == b.Attributes &&
		a.Info.Equal(b.Info) &&
		a.Height == b.Height
}

// DbProviderInfo represents the information stored inside the database about a single provider info
type DbProviderInfo struct {
	EMail   string
	Website string
}

// NewDbInfo builds a DbInfo starting from an akash provider info
func NewDbInfo(info providertypes.ProviderInfo) DbProviderInfo {
	return DbProviderInfo{
		EMail:   info.EMail,
		Website: info.Website,
	}
}

// Equal tells whether both object are the same
func (info DbProviderInfo) Equal(b DbProviderInfo) bool {
	return info.EMail == b.EMail && info.Website == b.Website
}

// Value implements driver.Valuer
func (info *DbProviderInfo) Value() (driver.Value, error) {
	return fmt.Sprintf("(%s,%s)", info.EMail, info.Website), nil
}

// DbAkashResource represents the information stored inside the database about a single akash resource
type DbAkashResource struct {
	CPU              uint64
	Memory           uint64
	StorageEphemeral uint64
}

// NewDbAkashResource builds a DbAkashResource from an akash resource
func NewDbAkashResource(resource *types.Resource) DbAkashResource {
	return DbAkashResource{
		CPU:              resource.CPU,
		Memory:           resource.Memory,
		StorageEphemeral: resource.StorageEphemeral,
	}
}

// Equal tells whether both object are the same
func (resource DbAkashResource) Equal(b DbAkashResource) bool {
	return resource.CPU == b.CPU && resource.Memory == b.Memory && resource.StorageEphemeral == b.StorageEphemeral
}

// Value implements driver.Valuer
func (resource *DbAkashResource) Value() (driver.Value, error) {
	return fmt.Sprintf("(%v,%v,%v)", resource.CPU, resource.Memory, resource.StorageEphemeral), nil
}
