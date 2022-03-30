package types

import (
	"time"
)

type (
	// DBTransfer represents a single row inside the "vipcoin_chain_banking_transfers" table
	DBTransfer struct {
		ID        uint64    `db:"id"`
		Asset     string    `db:"asset"`
		Amount    uint64    `db:"amount"`
		Kind      int32     `db:"kind"`
		Extras    ExtraDB   `db:"extras"`
		Timestamp time.Time `db:"timestamp"`
		TxHash    string    `db:"tx_hash"`
	}

	// DBPayment represents a single row inside the "vipcoin_chain_banking_payment" table
	DBPayment struct {
		DBTransfer
		WalletFrom string `db:"wallet_from"`
		WalletTo   string `db:"wallet_to"`
		Fee        uint64 `db:"fee"`
	}

	// DBSystemTransfer represents a single row inside the "vipcoin_chain_banking_system_transfer" table
	DBSystemTransfer struct {
		DBTransfer
		WalletFrom string `db:"wallet_from"`
		WalletTo   string `db:"wallet_to"`
	}

	// DBWithdraw represents a single row inside the "vipcoin_chain_banking_withdraw" table
	DBWithdraw struct {
		DBTransfer
		Wallet string `db:"wallet"`
	}

	// DBIssue represents a single row inside the "vipcoin_chain_banking_issue" table
	DBIssue struct {
		DBTransfer
		Wallet string `db:"wallet"`
	}
)
