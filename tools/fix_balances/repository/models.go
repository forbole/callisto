/*
 * Copyright (c) 2023. Business Process Technologies. All rights reserved.
 */

package repository

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/cosmos/cosmos-sdk/types"
)

const (
	// tableUsers is a name of the users table.
	tableWallets = "overgold_chain_wallets_wallets"
)

type (
	// DBWallets represents a single row inside the "overgold_chain_wallets_wallets" table
	DBWallets struct {
		Address        string    `db:"address"`
		AccountAddress string    `db:"account_address"`
		Kind           int32     `db:"kind"`
		State          int32     `db:"state"`
		Balance        BalanceDB `db:"balance"`
		Extras         string    `db:"extras"`
		DefaultStatus  bool      `db:"default_status"`
	}

	// // ExtraDB helprs type
	// ExtraDB struct {
	// 	Extras []extratypes.Extra
	// }

	// BalanceDB helpers type
	BalanceDB struct {
		Balance types.Coins
	}
)

// Value Make the BalanceDB struct implement the driver.Valuer interface. This method
// simply returns the JSON-encoded representation of the struct.
func (b BalanceDB) Value() (driver.Value, error) {
	return json.Marshal(b.Balance)
}

// Scan Make the BalanceDB struct implement the sql.Scanner interface. This method
// simply decodes a JSON-encoded value into the struct fields.
func (b *BalanceDB) Scan(value interface{}) error {
	v, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(v, &b.Balance)
}

// // Value Make the ExtraDB struct implement the driver.Valuer interface. This method
// // simply returns the JSON-encoded representation of the struct.
// func (e ExtraDB) Value() (driver.Value, error) {
// 	return json.Marshal(e.Extras)
// }
//
// // Scan Make the ExtraDB struct implement the sql.Scanner interface. This method
// // simply decodes a JSON-encoded value into the struct fields.
// func (e *ExtraDB) Scan(value interface{}) error {
// 	b, ok := value.([]byte)
// 	if !ok {
// 		return errors.New("type assertion to []byte failed")
// 	}
//
// 	return json.Unmarshal(b, &e.Extras)
// }
