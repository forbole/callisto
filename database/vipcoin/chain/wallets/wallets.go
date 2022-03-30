/*
 * Copyright 2022 Business Process Technologies. All rights reserved.
 */

package wallets

import (
	"context"
	"database/sql"

	walletstypes "git.ooo.ua/vipcoin/chain/x/wallets/types"
	"git.ooo.ua/vipcoin/lib/filter"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/jmoiron/sqlx"

	"github.com/forbole/bdjuno/v2/database/types"
)

type (
	// Repository - defines a repository for wallets repository
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

// SaveWallets - method that create wallets to the "vipcoin_chain_wallets_wallets" table
func (r Repository) SaveWallets(wallets ...*walletstypes.Wallet) error {
	if len(wallets) == 0 {
		return nil
	}

	tx, err := r.db.BeginTxx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	defer tx.Rollback()

	query := `INSERT INTO vipcoin_chain_wallets_wallets 
       ("address", "account_address", "kind", "state", "balance", "extras", "default_status") 
     VALUES 
       (:address, :account_address, :kind, :state, :balance, :extras, :default_status)`

	for _, wallet := range wallets {
		walletDB, err := toWalletsDatabase(wallet)
		if err != nil {
			return err
		}

		if _, err := tx.NamedExec(query, walletDB); err != nil {
			return err
		}
	}

	return tx.Commit()
}

// GetWallets - method that get wallets from the "vipcoin_chain_wallets_wallets" table
func (r Repository) GetWallets(filter filter.Filter) ([]*walletstypes.Wallet, error) {
	query, args := filter.Build("vipcoin_chain_wallets_wallets",
		`address, account_address, kind, state, balance, extras, default_status`)

	var result []types.DBWallets
	if err := r.db.Select(&result, query, args...); err != nil {
		return []*walletstypes.Wallet{}, err
	}

	wallets := make([]*walletstypes.Wallet, 0, len(result))
	for _, w := range result {
		wallet, err := toWalletDomain(w)
		if err != nil {
			return []*walletstypes.Wallet{}, err
		}

		wallets = append(wallets, wallet)
	}

	return wallets, nil
}

// UpdateWallets - method that updates the wallet in the "vipcoin_chain_wallets_wallets" table
func (r Repository) UpdateWallets(wallets ...*walletstypes.Wallet) error {
	if len(wallets) == 0 {
		return nil
	}

	tx, err := r.db.BeginTxx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	defer tx.Rollback()

	query := `UPDATE vipcoin_chain_wallets_wallets SET
				 address = :address, account_address = :account_address, kind = :kind,
				 state = :state, balance = :balance, extras = :extras, default_status = :default_status
			 WHERE address = :address`

	for _, wallet := range wallets {
		walletsDB, err := toWalletsDatabase(wallet)
		if err != nil {
			return err
		}

		if _, err := tx.NamedExec(query, walletsDB); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
