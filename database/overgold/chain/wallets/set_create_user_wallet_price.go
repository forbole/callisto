package wallets

import (
	"database/sql"
	"errors"

	walletstypes "git.ooo.ua/vipcoin/chain/x/wallets/types"
	"git.ooo.ua/vipcoin/lib/errs"
	"git.ooo.ua/vipcoin/lib/filter"

	"github.com/forbole/bdjuno/v3/database/types"
)

// SaveSetCreateUserWalletPrice - inserts messages into the "overgold_chain_wallets_set_create_user_wallet_price" table
func (r Repository) SaveSetCreateUserWalletPrice(messages *walletstypes.MsgSetCreateUserWalletPrice, transactionHash string) error {
	query := `INSERT INTO overgold_chain_wallets_set_create_user_wallet_price 
			(transaction_hash, creator, amount) 
			VALUES 
			(:transaction_hash, :creator, :amount)`

	if _, err := r.db.NamedExec(query, toSetCreateUserWalletPriceDatabase(messages, transactionHash)); err != nil {
		return errs.Internal{Cause: err.Error()}
	}

	return nil
}

// GetSetCreateUserWalletPrice - get the given messages from the "overgold_chain_wallets_set_create_user_wallet_price" table
func (r Repository) GetSetCreateUserWalletPrice(filter filter.Filter) (*walletstypes.MsgSetCreateUserWalletPrice, error) {
	query, args := filter.Build(types.TableOvergoldChainWalletsSetCreateUserWalletPrice, types.FieldCreator, types.FieldAmount, types.FieldTransactionHash)

	var wallet types.DBSetCreateUserWalletPrice
	if err := r.db.QueryRow(query, args...).Scan(&wallet.Creator, &wallet.Amount, &wallet.Hash); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &walletstypes.MsgSetCreateUserWalletPrice{}, errs.NotFound{What: "msg set create user wallet price"}
		}

		return &walletstypes.MsgSetCreateUserWalletPrice{}, errs.Internal{Cause: err.Error()}
	}

	return toSetCreateUserWalletPriceDomain(wallet), nil
}

// GetAllSetCreateUserWalletPrice - get the given messages from the "overgold_chain_wallets_set_create_user_wallet_price" table
func (r Repository) GetAllSetCreateUserWalletPrice(filter filter.Filter) ([]*walletstypes.MsgSetCreateUserWalletPrice, error) {
	query, args := filter.Build(types.TableOvergoldChainWalletsSetCreateUserWalletPrice, types.FieldCreator, types.FieldAmount)

	var wallets []types.DBSetCreateUserWalletPrice
	if err := r.db.Select(&wallets, query, args...); err != nil {
		return []*walletstypes.MsgSetCreateUserWalletPrice{}, errs.Internal{Cause: err.Error()}
	}

	result := make([]*walletstypes.MsgSetCreateUserWalletPrice, 0, len(wallets))
	for _, w := range wallets {
		result = append(result, toSetCreateUserWalletPriceDomain(w))
	}

	return result, nil
}
