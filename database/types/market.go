package types

import (
	"database/sql/driver"
	"fmt"

	escrowtypes "github.com/ovrclk/akash/x/escrow/types/v1beta2"
)

// DbLeaseAccountInfo represents a single lease account info stored inside the database
type DbLeaseAccountID struct {
	Scope string
	XID   string
}

// NewDbLeaseAccountID builds a DbLeaseAccountInfo from a lease account ID
func NewDbLeaseAccountID(info escrowtypes.AccountID) DbProviderInfo {
	return DbProviderInfo{
		EMail:   info.Scope,
		Website: info.XID,
	}
}

// Value implements driver.Valuer
func (info *DbLeaseAccountID) Value() (driver.Value, error) {
	return fmt.Sprintf("(%s,%s)", info.Scope, info.XID), nil
}
