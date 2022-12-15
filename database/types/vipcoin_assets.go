package types

import (
	"github.com/lib/pq"
)

type (
	// DBAssets represents a single row inside the "overgold_chain_assets_assets" table
	DBAssets struct {
		Issuer        string        `db:"issuer"`
		Name          string        `db:"name"`
		Policies      pq.Int32Array `db:"policies"`
		State         int32         `db:"state"`
		Issued        uint64        `db:"issued"`
		Burned        uint64        `db:"burned"`
		Withdrawn     uint64        `db:"withdrawn"`
		InCirculation uint64        `db:"in_circulation"`
		Precision     int64         `db:"precision"`
		FeePercent    int64         `db:"fee_percent"`
		Extras        ExtraDB       `db:"extras"`
	}

	// DBAssetCreate represents a single row inside the "overgold_chain_assets_create" table
	DBAssetCreate struct {
		Hash       string        `db:"transaction_hash"`
		Creator    string        `db:"creator"`
		Name       string        `db:"name"`
		Issuer     string        `db:"issuer"`
		Policies   pq.Int32Array `db:"policies"`
		State      int32         `db:"state"`
		Precision  int64         `db:"precision"`
		FeePercent int64         `db:"fee_percent"`
		Extras     ExtraDB       `db:"extras"`
	}

	// DBAssetManage represents a single row inside the "overgold_chain_assets_manage" table
	DBAssetManage struct {
		Hash          string        `db:"transaction_hash"`
		Creator       string        `db:"creator"`
		Name          string        `db:"name"`
		Policies      pq.Int32Array `db:"policies"`
		State         int32         `db:"state"`
		Precision     int64         `db:"precision"`
		FeePercent    int64         `db:"fee_percent"`
		Issued        uint64        `db:"issued"`
		Burned        uint64        `db:"burned"`
		Withdrawn     uint64        `db:"withdrawn"`
		InCirculation uint64        `db:"in_circulation"`
	}

	// DBAssetSetExtra represents a single row inside the "overgold_chain_assets_set_extra" table
	DBAssetSetExtra struct {
		Hash    string  `db:"transaction_hash"`
		Creator string  `db:"creator"`
		Name    string  `db:"name"`
		Extras  ExtraDB `db:"extras"`
	}
)
