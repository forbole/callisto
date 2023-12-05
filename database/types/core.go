package types

type (
	// CoreMsgIssue - db model for 'overgold_core_issue'
	CoreMsgIssue struct {
		ID      uint64 `db:"id"`
		TxHash  string `db:"tx_hash"`
		Creator string `db:"creator"`
		Amount  uint64 `db:"amount"`
		Denom   string `db:"denom"`
		Address string `db:"address"`
	}

	// CoreMsgWithdraw - db model for 'overgold_core_withdraw'
	CoreMsgWithdraw struct {
		ID      uint64 `db:"id"`
		TxHash  string `db:"tx_hash"`
		Creator string `db:"creator"`
		Amount  uint64 `db:"amount"`
		Denom   string `db:"denom"`
		Address string `db:"address"`
	}

	// CoreMsgSend - db model for 'overgold_core_send'
	CoreMsgSend struct {
		ID          uint64 `db:"id"`
		TxHash      string `db:"tx_hash"`
		Creator     string `db:"creator"`
		AddressFrom string `db:"address_from"`
		AddressTo   string `db:"address_to"`
		Amount      uint64 `db:"amount"`
		Denom       string `db:"denom"`
	}
)
