package accounts

import (
	"context"
	"database/sql"

	accountstypes "git.ooo.ua/vipcoin/chain/x/accounts/types"
	"git.ooo.ua/vipcoin/lib/filter"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/forbole/bdjuno/v2/database/types"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type (
	// repository - defines a repository for accounts repository
	Repository struct {
		db  *sqlx.DB
		cdc codec.Marshaler
	}
)

// NewRepository constructor.
func NewRepository(db *sqlx.DB, cdc codec.Marshaler) *Repository {
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
		return err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	query := `INSERT INTO vipcoin_chain_accounts_accounts 
			 ("address", "hash", "public_key", "kinds", "state", "extras", "affiliates", "wallets") 
		 VALUES 
			 (:address, :hash, :public_key, :kinds, :state, :extras, :affiliates, :wallets)`

	for _, acc := range accounts {
		accountDB, err := toAccountDatabase(acc, r.cdc)
		if err != nil {
			return err
		}

		if accountDB.Affiliates, err = saveAffiliates(tx, acc.Affiliates); err != nil {
			return err
		}

		if _, err := tx.NamedExec(query, accountDB); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func saveAffiliates(tx *sqlx.Tx, affiliates []*accountstypes.Affiliate) (pq.Int64Array, error) {
	if len(affiliates) == 0 {
		return pq.Int64Array{}, nil
	}

	query := `INSERT INTO vipcoin_chain_accounts_affiliates
		 (address, affiliation_kind, extras)
	 VALUES
		 (:address, :affiliation_kind, :extras)
	 RETURNING id`

	resultID := make(pq.Int64Array, len(affiliates))
	for index, affiliate := range affiliates {
		resp, err := tx.NamedQuery(query, toAffiliatesDatabase(affiliate))
		if err != nil {
			return pq.Int64Array{}, err
		}

		for resp.Next() {
			if err := resp.Scan(&resultID[index]); err != nil {
				return pq.Int64Array{}, err
			}
		}

		if err := resp.Err(); err != nil {
			return pq.Int64Array{}, err
		}
	}

	return resultID, nil
}

func (r Repository) UpdateAccounts(accounts ...*accountstypes.Account) error {
	if len(accounts) == 0 {
		return nil
	}

	tx, err := r.db.BeginTxx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	query := `UPDATE vipcoin_chain_accounts_accounts SET
				 address = :address, public_key = :public_key, kinds = :kinds,
				 state = :state, extras = :extras, affiliates = :affiliates, wallets = :wallets
			 WHERE hash = :hash`

	queryAffiliates := "SELECT affiliates FROM vipcoin_chain_accounts_accounts WHERE hash = $1"

	for _, acc := range accounts {
		accountDB, err := toAccountDatabase(acc, r.cdc)
		if err != nil {
			return err
		}

		var affiliatesID pq.Int64Array
		if err := tx.Get(&affiliatesID, queryAffiliates, accountDB.Hash); err != nil {
			return err
		}

		if err := deleteAffiliates(tx, affiliatesID); err != nil {
			return err
		}

		if accountDB.Affiliates, err = saveAffiliates(tx, acc.Affiliates); err != nil {
			return err
		}

		if _, err := tx.NamedExec(query, accountDB); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func deleteAffiliates(tx *sqlx.Tx, affiliatesID pq.Int64Array) error {
	if len(affiliatesID) == 0 {
		return nil
	}

	query := `DELETE FROM vipcoin_chain_accounts_affiliates 
			 WHERE id=($1)`

	for _, id := range affiliatesID {
		if _, err := tx.Exec(query, id); err != nil {
			return err
		}
	}

	return nil
}

func (r Repository) GetAccounts(accfilter filter.Filter) ([]*accountstypes.Account, error) {
	query, args := accfilter.Build("vipcoin_chain_accounts_accounts",
		`address, hash, public_key, kinds, state, extras, affiliates, wallets`)

	var result []types.DBAccount
	if err := r.db.Select(&result, query, args...); err != nil {
		return []*accountstypes.Account{}, err
	}

	accounts := make([]*accountstypes.Account, 0, len(result))
	for _, acc := range result {
		account, err := toAccountDomain(acc)
		if err != nil {
			return []*accountstypes.Account{}, err
		}

		if account.Affiliates, err = getAffiliates(r.db, acc.Affiliates); err != nil {
			return []*accountstypes.Account{}, err
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}

func getAffiliates(db *sqlx.DB, affiliatesID pq.Int64Array) ([]*accountstypes.Affiliate, error) {
	if len(affiliatesID) == 0 {
		return []*accountstypes.Affiliate{}, nil
	}

	query, args := filter.NewFilter().SetArgument("id", parseID(affiliatesID)...).
		Build("vipcoin_chain_accounts_affiliates", `id, address, affiliation_kind, extras`)

	var result []types.DBAffiliates
	if err := db.Select(&result, query, args...); err != nil {
		return []*accountstypes.Affiliate{}, err
	}

	return toAffiliatesDomain(result), nil
}
