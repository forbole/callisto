/*
 * Copyright 2022 Business Process Technologies. All rights reserved.
 */

package types

import (
	"github.com/cosmos/cosmos-sdk/types"
)

type (
	// DBWallets represents a single row inside the "vipcoin_chain_wallets_wallets" table
	DBWallets struct {
		Address        string    `db:"address"`
		AccountAddress string    `db:"account_address"`
		Kind           int32     `db:"kind"`
		State          int32     `db:"state"`
		Balance        BalanceDB `db:"balance"`
		Extras         ExtraDB   `db:"extras"`
		DefaultStatus  bool      `db:"default_status"`
	}

	// DBSetWalletKind represents a single row inside the "vipcoin_chain_wallets_set_wallet_kind" table
	DBSetWalletKind struct {
		Creator string `db:"creator"`
		Address string `db:"address"`
		Kind    int32  `db:"kind"`
	}

	// DBSetWalletState represents a single row inside the "vipcoin_chain_wallets_set_wallet_state" table
	DBSetWalletState struct {
		Creator string `db:"creator"`
		Address string `db:"address"`
		State   int32  `db:"state"`
	}

	// DBCreateWallet represents a single row inside the "vipcoin_chain_wallets_create_wallet" table
	DBCreateWallet struct {
		Creator        string `db:"creator"`
		Address        string `db:"address"`
		AccountAddress string `db:"account_address"`
		Kind           int32  `db:"kind"`
		State          int32  `db:"state"`
	}

	// DBCreateWalletWithBalance represents a single row inside the "vipcoin_chain_wallets_create_wallet_with_balance" table
	DBCreateWalletWithBalance struct {
		Creator        string  `db:"creator"`
		Address        string  `db:"address"`
		AccountAddress string  `db:"account_address"`
		Kind           int32   `db:"kind"`
		State          int32   `db:"state"`
		Extras         ExtraDB `db:"extras"`
		DefaultStatus  bool    `db:"default_status"`
		Balance        ExtraDB `db:"balance"`
	}

	// DBSetDefaultWallet represents a single row inside the "vipcoin_chain_wallets_set_default_wallet" table
	DBSetDefaultWallet struct {
		Creator string `db:"creator"`
		Address string `db:"address"`
	}

	// DBSetExtra represents a single row inside the "vipcoin_chain_wallets_set_extra" table
	DBSetExtra struct {
		Creator string  `db:"creator"`
		Address string  `db:"address"`
		Extras  ExtraDB `db:"extras"`
	}

	// BalanceDB helpers type
	BalanceDB struct {
		Balance types.Coins
	}
)
