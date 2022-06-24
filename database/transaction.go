package database

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"time"

	"git.ooo.ua/vipcoin/lib/errs"
	"git.ooo.ua/vipcoin/lib/filter"
	cosmos "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	cosmostx "github.com/cosmos/cosmos-sdk/types/tx"
	txtypes "github.com/forbole/juno/v2/types"

	"github.com/forbole/bdjuno/v2/database/types"
)

// GetTransaction - get transaction from database
func (db *Db) GetTransaction(filter filter.Filter) (*txtypes.Tx, error) {
	query, args := filter.SetLimit(1).Build("transaction")

	var result types.TransactionRow
	if err := db.Sqlx.Get(&result, query, args...); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return &txtypes.Tx{}, errs.Internal{Cause: err.Error()}
		}

		return &txtypes.Tx{}, errs.NotFound{What: "transaction"}
	}

	return db.toTxTypesTx(result)
}

// GetTransactions - get transactions from database
func (db *Db) GetTransactions(filter filter.Filter) ([]*txtypes.Tx, error) {
	query, args := filter.Build("transaction")

	var result []types.TransactionRow
	if err := db.Sqlx.Select(&result, query, args...); err != nil {
		return []*txtypes.Tx{}, errs.Internal{Cause: err.Error()}
	}

	if len(result) == 0 {
		return []*txtypes.Tx{}, errs.NotFound{What: "transaction"}
	}

	transactions := make([]*txtypes.Tx, 0, len(result))
	for _, tx := range result {
		transaction, err := db.toTxTypesTx(tx)
		if err != nil {
			return []*txtypes.Tx{}, errs.Internal{Cause: err.Error()}
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

// toTxTypesTx - convert database row to Tx
func (db *Db) toTxTypesTx(tx types.TransactionRow) (*txtypes.Tx, error) {
	var err error
	result, err := txtypes.NewTx(
		&sdk.TxResponse{},
		&cosmostx.Tx{Body: &cosmostx.TxBody{}, AuthInfo: &cosmostx.AuthInfo{}},
	)

	if err != nil {
		return nil, err
	}

	if !tx.Success {
		result.TxResponse.Code = 1
	}

	var anyRaw []json.RawMessage
	if err = json.Unmarshal(tx.Messages, &anyRaw); err != nil {
		return nil, err
	}

	result.Body.Messages = make([]*cosmos.Any, 0, len(anyRaw))
	for _, raw := range anyRaw {
		msg := cosmos.Any{}

		if err := db.EncodingConfig.Marshaler.UnmarshalJSON(raw, &msg); err != nil {
			return nil, err
		}

		result.Body.Messages = append(result.Body.Messages, &msg)
	}

	result.Signatures = make([][]byte, len(tx.Signatures))
	for index, sig := range tx.Signatures {
		if result.Signatures[index], err = base64.StdEncoding.DecodeString(sig); err != nil {
			return nil, err
		}
	}

	var sigInfoRaw []json.RawMessage
	if err = json.Unmarshal(tx.SignerInfos, &sigInfoRaw); err != nil {
		return nil, err
	}

	result.AuthInfo.SignerInfos = make([]*cosmostx.SignerInfo, 0, len(sigInfoRaw))
	for _, sig := range sigInfoRaw {
		sigInfo := cosmostx.SignerInfo{}

		if err = db.EncodingConfig.Marshaler.UnmarshalJSON(sig, &sigInfo); err != nil {
			return nil, err
		}

		result.AuthInfo.SignerInfos = append(result.AuthInfo.SignerInfos, &sigInfo)
	}

	result.AuthInfo.Fee = &cosmostx.Fee{}
	if err = db.EncodingConfig.Marshaler.UnmarshalJSON(tx.Fee, result.AuthInfo.Fee); err != nil {
		return nil, err
	}

	result.Logs = sdk.ABCIMessageLogs{}
	if err = json.Unmarshal(tx.Logs, &result.Logs); err != nil {
		return nil, err
	}

	block, err := db.GetBlock(filter.NewFilter().SetArgument(types.FieldHeight, tx.Height))
	if err != nil {
		return nil, err
	}

	result.TxHash = tx.Hash
	result.Height = tx.Height
	result.Body.Memo = tx.Memo
	result.GasWanted = tx.GasWanted
	result.GasUsed = tx.GasUsed
	result.RawLog = tx.RawLog
	result.Timestamp = block.Timestamp.Format(time.RFC3339)

	return result, nil
}
