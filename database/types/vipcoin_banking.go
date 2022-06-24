package types

import (
	"time"
)

type (
	// DBTransfer represents a single row inside the "vipcoin_chain_banking_base_transfers" table
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

	// DBSystemTransfer represents a single row inside the "vipcoin_chain_banking_msg_system_transfer" table
	DBMsgSystemTransfer struct {
		Hash       string  `db:"transaction_hash"`
		Creator    string  `db:"creator"`
		WalletFrom string  `db:"wallet_from"`
		WalletTo   string  `db:"wallet_to"`
		Asset      string  `db:"asset"`
		Amount     uint64  `db:"amount"`
		Extras     ExtraDB `db:"extras"`
	}

	// DBSystemRewardTransfer represents a single row inside the "vipcoin_chain_banking_system_msg_reward_transfer" table
	DBMsgSystemRewardTransfer struct {
		Hash       string  `db:"transaction_hash"`
		Creator    string  `db:"creator"`
		WalletFrom string  `db:"wallet_from"`
		WalletTo   string  `db:"wallet_to"`
		Asset      string  `db:"asset"`
		Amount     uint64  `db:"amount"`
		Extras     ExtraDB `db:"extras"`
	}

	// DBPayment represents a single row inside the "vipcoin_chain_banking_msg_payment" table
	DBMsgPayment struct {
		Hash       string  `db:"transaction_hash"`
		Creator    string  `db:"creator"`
		WalletFrom string  `db:"wallet_from"`
		WalletTo   string  `db:"wallet_to"`
		Asset      string  `db:"asset"`
		Amount     uint64  `db:"amount"`
		Extras     ExtraDB `db:"extras"`
	}

	// DBWithdraw represents a single row inside the "vipcoin_chain_banking_msg_withdraw" table
	DBMsgWithdraw struct {
		Hash    string  `db:"transaction_hash"`
		Creator string  `db:"creator"`
		Wallet  string  `db:"wallet"`
		Asset   string  `db:"asset"`
		Amount  uint64  `db:"amount"`
		Extras  ExtraDB `db:"extras"`
	}

	// DBIssue represents a single row inside the "vipcoin_chain_banking_msg_issue" table
	DBMsgIssue struct {
		Hash    string  `db:"transaction_hash"`
		Creator string  `db:"creator"`
		Wallet  string  `db:"wallet"`
		Asset   string  `db:"asset"`
		Amount  uint64  `db:"amount"`
		Extras  ExtraDB `db:"extras"`
	}

	// DBSetTransferExtra represents a single row inside the "vipcoin_chain_banking_set_transfer_extra" table
	DBSetTransferExtra struct {
		Hash    string  `db:"transaction_hash"`
		Creator string  `db:"creator"`
		Id      uint64  `db:"id"`
		Extras  ExtraDB `db:"extras"`
	}

	// DBSetRewardManagerAddress represents a single row inside the "vipcoin_chain_banking_set_reward_manager_address" table
	DBSetRewardManagerAddress struct {
		Hash    string `db:"transaction_hash"`
		Creator string `db:"creator"`
		Address string `db:"address"`
	}
)
