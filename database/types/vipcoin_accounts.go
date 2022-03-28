/*
 * Copyright 2022 Business Process Technologies. All rights reserved.
 */

package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	extratypes "git.ooo.ua/vipcoin/chain/x/types"
	"github.com/lib/pq"
)

type (
	// DBAccount represents a single row inside the "vipcoin_chain_accounts_accounts" table
	DBAccount struct {
		Address    string         `db:"address"`
		Hash       string         `db:"hash"`
		PublicKey  string         `db:"public_key"`
		Kinds      pq.Int32Array  `db:"kinds"`
		State      int32          `db:"state"`
		Extras     ExtraDB        `db:"extras"`
		Affiliates pq.Int64Array  `db:"affiliates"`
		Wallets    pq.StringArray `db:"wallets"`
	}

	// DBAffiliates represents a single row inside the "vipcoin_chain_accounts_affiliates" table
	DBAffiliates struct {
		Id              uint64  `db:"id"`
		Address         string  `db:"address"`
		AffiliationKind int32   `db:"affiliation_kind"`
		Extras          ExtraDB `db:"extras"`
	}

	// ExtraDB helprs type
	ExtraDB struct {
		Extras []extratypes.Extra
	}
)

// Make the ExtraDB struct implement the driver.Valuer interface. This method
// simply returns the JSON-encoded representation of the struct.
func (e ExtraDB) Value() (driver.Value, error) {
	return json.Marshal(e.Extras)
}

// Make the ExtraDB struct implement the sql.Scanner interface. This method
// simply decodes a JSON-encoded value into the struct fields.
func (e *ExtraDB) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &e.Extras)
}
