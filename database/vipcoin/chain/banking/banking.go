/*
 * Copyright 2022 Business Process Technologies. All rights reserved.
 */

package banking

import (
	"database/sql"

	bankingtypes "git.ooo.ua/vipcoin/chain/x/banking/types"
	"git.ooo.ua/vipcoin/lib/filter"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/jmoiron/sqlx"
	"golang.org/x/net/context"

	"github.com/forbole/bdjuno/v2/database/types"
)

type (
	// Repository - defines a repository for banking repository
	Repository struct {
		db  *sqlx.DB
		cdc codec.Marshaler
	}
)

// NewRepository constructor
func NewRepository(db *sqlx.DB, cdc codec.Marshaler) *Repository {
	return &Repository{
		db:  db,
		cdc: cdc,
	}
}

// SaveBaseTransfers - method that create transfers to the "vipcoin_chain_banking_transfers" table
func (r Repository) SaveBaseTransfers(transfers ...*bankingtypes.BaseTransfer) error {
	if len(transfers) == 0 {
		return nil
	}

	tx, err := r.db.BeginTxx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	defer tx.Rollback()

	query := `INSERT INTO vipcoin_chain_banking_transfers 
       ("id", "asset", "amount", "kind", "extras", "timestamp", "tx_hash") 
     VALUES 
       (:id, :asset, :amount, :kind, :extras, :timestamp, :tx_hash)`

	for _, transfer := range transfers {
		transferDB := toTransferDatabase(transfer)

		if _, err := tx.NamedExec(query, transferDB); err != nil {
			return err
		}
	}

	return tx.Commit()
}

// GetBaseTransfers - method that get transfers from the "vipcoin_chain_banking_transfers" table
func (r Repository) GetBaseTransfers(filter filter.Filter) ([]*bankingtypes.BaseTransfer, error) {
	query, args := filter.Build("vipcoin_chain_banking_transfers",
		`id, asset, amount, kind, extras, timestamp, tx_hash`)

	var result []types.DBTransfer
	if err := r.db.Select(&result, query, args...); err != nil {
		return []*bankingtypes.BaseTransfer{}, err
	}

	transfers := make([]*bankingtypes.BaseTransfer, 0, len(result))
	for _, transfer := range result {
		transfers = append(transfers, toTransferDomain(transfer))
	}

	return transfers, nil
}
