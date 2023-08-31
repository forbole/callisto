package accounts

import (
	"context"
	"database/sql"

	accountstypes "git.ooo.ua/vipcoin/chain/x/accounts/types"
	"git.ooo.ua/vipcoin/lib/errs"
	"git.ooo.ua/vipcoin/lib/filter"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/jmoiron/sqlx"

	"github.com/forbole/bdjuno/v3/database/types"
)

type (
	// Repository - defines a repository for accounts repository
	Repository struct {
		db  *sqlx.DB
		cdc codec.Codec
	}
)

// NewRepository constructor.
func NewRepository(db *sqlx.DB, cdc codec.Codec) *Repository {
	return &Repository{
		db:  db,
		cdc: cdc,
	}
}

func (r Repository) SaveAccounts(accounts ...*accountstypes.Account) error {
	if len(accounts) == 0 {
		return nil
	}

	tx, err := r.db.BeginTxx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return errs.Internal{Cause: err.Error()}
	}

	defer func() {
		_ = tx.Rollback()
	}()

	query := `INSERT INTO overgold_chain_accounts_accounts 
			 ("address", "hash", "public_key", "kinds", "state", "extras", "wallets") 
		 VALUES 
			 (:address, :hash, :public_key, :kinds, :state, :extras, :wallets)`

	for _, acc := range accounts {
		accountDB, err := toAccountDatabase(acc, r.cdc)
		if err != nil {
			return errs.Internal{Cause: err.Error()}
		}

		if _, err := tx.NamedExec(query, accountDB); err != nil {
			return errs.Internal{Cause: err.Error()}
		}

		if err = saveAffiliates(tx, acc.Affiliates, accountDB.Hash); err != nil {
			return errs.Internal{Cause: err.Error()}
		}
	}

	return tx.Commit()
}

func saveAffiliates(tx *sqlx.Tx, affiliates []*accountstypes.Affiliate, accHash string) error {
	if len(affiliates) == 0 {
		return nil
	}

	query := `INSERT INTO overgold_chain_accounts_affiliates
		 (account_hash, address, affiliation_kind, extras)
	 VALUES
		 (:account_hash, :address, :affiliation_kind, :extras)`

	if _, err := tx.NamedExec(query, toAffiliatesDatabase(affiliates, accHash)); err != nil {
		return err
	}

	return nil
}

func (r Repository) UpdateAccounts(accounts ...*accountstypes.Account) error {
	if len(accounts) == 0 {
		return nil
	}

	tx, err := r.db.BeginTxx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return errs.Internal{Cause: err.Error()}
	}

	defer func() {
		_ = tx.Rollback()
	}()

	query := `UPDATE overgold_chain_accounts_accounts SET
				 address = :address, public_key = :public_key, kinds = :kinds,
				 state = :state, extras = :extras, wallets = :wallets
			 WHERE hash = :hash`

	deleteAffiliates := "DELETE FROM overgold_chain_accounts_affiliates WHERE account_hash=($1)"

	for _, acc := range accounts {
		accountDB, err := toAccountDatabase(acc, r.cdc)
		if err != nil {
			return errs.Internal{Cause: err.Error()}
		}

		if _, err := tx.Exec(deleteAffiliates, accountDB.Hash); err != nil {
			return errs.Internal{Cause: err.Error()}
		}

		if _, err := tx.NamedExec(query, accountDB); err != nil {
			return errs.Internal{Cause: err.Error()}
		}

		if err = saveAffiliates(tx, acc.Affiliates, acc.Hash); err != nil {
			return errs.Internal{Cause: err.Error()}
		}
	}

	return tx.Commit()
}

func (r Repository) GetAccounts(accountFilter filter.Filter) ([]*accountstypes.Account, error) {
	query, args := accountFilter.Build(
		tableAccounts,
		types.FieldAddress, types.FieldHash, types.FieldPublicKey, types.FieldKinds,
		types.FieldState, types.FieldExtras, types.FieldWallets,
	)

	var result []types.DBAccount
	if err := r.db.Select(&result, query, args...); err != nil {
		return []*accountstypes.Account{}, errs.Internal{Cause: err.Error()}
	}

	accounts := make([]*accountstypes.Account, 0, len(result))
	for _, acc := range result {
		account, err := toAccountDomain(acc)
		if err != nil {
			return []*accountstypes.Account{}, errs.Internal{Cause: err.Error()}
		}

		if account.Affiliates, err = getAffiliates(r.db, acc.Hash); err != nil {
			return []*accountstypes.Account{}, errs.Internal{Cause: err.Error()}
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}

func getAffiliates(db *sqlx.DB, accHash string) ([]*accountstypes.Affiliate, error) {
	query, args := filter.NewFilter().SetArgument(types.FieldAccountHash, accHash).Build(
		tableAffiliates,
		types.FieldID, types.FieldAccountHash, types.FieldAddress,
		types.FieldAffiliationKind, types.FieldExtras,
	)

	var result []types.DBAffiliates
	if err := db.Select(&result, query, args...); err != nil {
		return []*accountstypes.Affiliate{}, err
	}

	return toAffiliatesDomain(result), nil
}
