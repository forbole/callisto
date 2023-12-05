package types

type (
	// MsgSend represents a single row inside the 'msg_send' table
	MsgSend struct {
		ID          uint64  `db:"id"`
		TxHash      string  `db:"tx_hash"`
		FromAddress string  `db:"from_address"`
		ToAddress   string  `db:"to_address"`
		Amount      DbCoins `db:"amount"`
	}

	// MsgMultiSend represents a single row inside the 'msg_multi_send' table
	MsgMultiSend struct {
		ID     uint64         `db:"id"`
		TxHash string         `db:"tx_hash"`
		Inputs DbSendDataList `db:"inputs"`
		Ouputs DbSendDataList `db:"outputs"`
	}
)
