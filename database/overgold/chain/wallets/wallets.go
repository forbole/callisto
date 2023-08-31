package wallets

import (
	"context"
	"database/sql"

	walletstypes "git.ooo.ua/vipcoin/chain/x/wallets/types"
	"git.ooo.ua/vipcoin/lib/errs"
	"git.ooo.ua/vipcoin/lib/filter"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/jmoiron/sqlx"

	"github.com/forbole/bdjuno/v3/database/types"
)

type (
	// Repository - defines a repository for wallets repository
	Repository struct {
		db  *sqlx.DB
		cdc codec.Codec
	}
)

// NewRepository constructor
func NewRepository(db *sqlx.DB, cdc codec.Codec) *Repository {
	return &Repository{
		db:  db,
		cdc: cdc,
	}
}

// SaveWallets - method that create wallets to the "overgold_chain_wallets_wallets" table
func (r Repository) SaveWallets(wallets ...*walletstypes.Wallet) error {
	if len(wallets) == 0 {
		return nil
	}

	tx, err := r.db.BeginTxx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return errs.Internal{Cause: err.Error()}
	}

	defer tx.Rollback()

	query := `INSERT INTO overgold_chain_wallets_wallets 
       ("address", "account_address", "kind", "state", "balance", "extras", "default_status") 
     VALUES 
       (:address, :account_address, :kind, :state, :balance, :extras, :default_status)`

	for _, wallet := range wallets {
		if _, err := tx.NamedExec(query, toWalletsDatabase(wallet)); err != nil {
			return errs.Internal{Cause: err.Error()}
		}
	}

	return tx.Commit()
}

// GetWallets - method that get wallets from the "overgold_chain_wallets_wallets" table
func (r Repository) GetWallets(filter filter.Filter) ([]*walletstypes.Wallet, error) {
	query, args := filter.Build("overgold_chain_wallets_wallets",
		types.FieldAddress, types.FieldAccountAddress, types.FieldKind, types.FieldState, types.FieldBalance, types.FieldExtras, types.FieldDefaultStatus)

	var result []types.DBWallets
	if err := r.db.Select(&result, query, args...); err != nil {
		return []*walletstypes.Wallet{}, errs.Internal{Cause: err.Error()}
	}

	wallets := make([]*walletstypes.Wallet, 0, len(result))
	for _, w := range result {
		wallets = append(wallets, toWalletDomain(w))
	}

	return wallets, nil
}

// UpdateWallets - method that updates the wallet in the "overgold_chain_wallets_wallets" table
func (r Repository) UpdateWallets(wallets ...*walletstypes.Wallet) error {
	if len(wallets) == 0 {
		return nil
	}

	tx, err := r.db.BeginTxx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return errs.Internal{Cause: err.Error()}
	}

	defer tx.Rollback()

	query := `UPDATE overgold_chain_wallets_wallets SET
				 account_address = :account_address, kind = :kind,
				 state = :state, balance = :balance, extras = :extras, default_status = :default_status
			 WHERE address = :address`

	for _, wallet := range wallets {
		if _, err := tx.NamedExec(query, toWalletsDatabase(wallet)); err != nil {
			return errs.Internal{Cause: err.Error()}
		}
	}

	return tx.Commit()
}

// DeleteWallets - method that delete wallets in the "overgold_chain_wallets_wallets" table
func (r Repository) DeleteWallets(wallets ...*walletstypes.Wallet) error {
	if len(wallets) == 0 {
		return nil
	}

	tx, err := r.db.BeginTxx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return errs.Internal{Cause: err.Error()}
	}

	defer tx.Rollback()

	query := `DELETE FROM overgold_chain_wallets_wallets WHERE address = :address`

	for _, wallet := range wallets {
		if _, err := tx.NamedExec(query, toWalletsDatabase(wallet)); err != nil {
			return errs.Internal{Cause: err.Error()}
		}
	}

	return tx.Commit()
}
