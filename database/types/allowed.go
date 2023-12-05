package types

import "github.com/lib/pq"

type (
	// AllowedAddresses - db model for 'overgold_allowed_addresses'
	AllowedAddresses struct {
		ID      uint64         `db:"id"`
		Creator string         `db:"creator"`
		Address pq.StringArray `db:"address"`
	}

	// AllowedCreateAddresses - db model for 'overgold_allowed_create_addresses'
	AllowedCreateAddresses struct {
		ID      uint64         `db:"id"`
		TxHash  string         `db:"tx_hash"`
		Creator string         `db:"creator"`
		Address pq.StringArray `db:"address"`
	}

	// AllowedUpdateAddresses - db model for 'overgold_allowed_update_addresses'
	AllowedUpdateAddresses struct {
		ID      uint64         `db:"id"`
		TxHash  string         `db:"tx_hash"`
		Creator string         `db:"creator"`
		Address pq.StringArray `db:"address"`
	}

	// AllowedDeleteByID - db model for 'overgold_allowed_delete_by_id'
	AllowedDeleteByID struct {
		ID      uint64 `db:"id"`
		TxHash  string `db:"tx_hash"`
		Creator string `db:"creator"`
	}

	// AllowedDeleteByAddresses - db model for 'overgold_allowed_delete_by_addresses'
	AllowedDeleteByAddresses struct {
		ID      uint64         `db:"id"`
		TxHash  string         `db:"tx_hash"`
		Creator string         `db:"creator"`
		Address pq.StringArray `db:"address"`
	}
)
