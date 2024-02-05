package v5

type MessageRow struct {
	TransactionHash           string `db:"transaction_hash"`
	Index                     int64  `db:"index"`
	Type                      string `db:"type"`
	Value                     string `db:"value"`
	InvolvedAccountsAddresses string `db:"involved_accounts_addresses"`
	Height                    int64  `db:"height"`
}
