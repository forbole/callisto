/*
 * Copyright (c) 2023. Business Process Technologies. All rights reserved.
 */

package repository

import (
	"context"
	"fmt"

	walletstypes "git.ooo.ua/vipcoin/chain-client/pkg/api/v4/chain/wallets"
	"git.ooo.ua/vipcoin/chain-client/pkg/client"
	"git.ooo.ua/vipcoin/lib/database"
	"git.ooo.ua/vipcoin/lib/errs"
	"git.ooo.ua/vipcoin/lib/filter"
)

// Repository represents a repository for the users.
type Repository struct {
	db          database.Executor
	chainClient client.VCG
}

// NewRepository constructor Repository.
func NewRepository(db database.Executor, chainClient client.VCG) *Repository {
	return &Repository{db: db, chainClient: chainClient}
}

// GetWallets - method that get wallets from the "overgold_chain_wallets_wallets" table
func (r Repository) GetWallets(ctx context.Context, filter filter.Filter) ([]DBWallets, error) {
	query, args := filter.Build(tableWallets)

	var result []DBWallets
	if err := r.db.SelectContext(ctx, &result, query, args...); err != nil {
		return []DBWallets{}, errs.Internal{Cause: err.Error()}
	}

	if err := r.db.SelectContext(ctx, &result, query, args...); err != nil {
		return nil, errs.Internal{Cause: err.Error()}
	}

	return result, nil

}

// UpdateWallets - method that updates the wallet in the "overgold_chain_wallets_wallets" table
func (r Repository) UpdateWallets(ctx context.Context, wallets ...DBWallets) error {
	query := `
		UPDATE overgold_chain_wallets_wallets
		SET
		    account_address = :account_address, kind = :kind,
		    state = :state, balance = :balance, extras = :extras, default_status = :default_status
		WHERE address = :address
	`

	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return errs.Internal{Cause: err.Error()}
	}

	wrappedStmt := r.db.NewStatement(stmt)

	for _, wallet := range wallets {
		if _, err = wrappedStmt.ExecContext(ctx, wallet); err != nil {
			return errs.Internal{Cause: err.Error()}
		}
	}

	return nil

}

// GetAllWalletFromChain - method that get all wallets from the chain // map[string]*walletstypes.Wallet,
func (r Repository) GetAllWalletFromChain(ctx context.Context) (map[string]*walletstypes.Wallet, error) {
	wallets, err := r.chainClient.GetAllWallets(ctx, client.SetWalletPageLimit(1000000000000000))
	if err != nil {
		return nil, err
	}

	result := make(map[string]*walletstypes.Wallet, len(wallets.Wallets))
	for i, wallet := range wallets.Wallets {
		fmt.Println(i, " кошельков обработано из ", len(wallets.Wallets))
		result[wallet.Address] = wallet
	}

	return result, nil
}
