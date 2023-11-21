package allowed

import "github.com/lib/pq"

const (
	tableAddresses         = "overgold_allowed_addresses"
	tableCreateAddresses   = "overgold_allowed_create_addresses"
	tableDeleteByAddresses = "overgold_allowed_delete_by_addresses"
	tableDeleteByID        = "overgold_allowed_delete_by_id"
	tableUpdateAddresses   = "overgold_allowed_update_addresses"
)

type (
	// deleteAddressesDB - helper struct for db to overgold_allowed_addresses.
	deleteAddressesDB struct {
		Address pq.StringArray `db:"address"`
		ID      uint64         `db:"id"`
	}
)
