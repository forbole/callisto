package bank

import (
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/lib/pq"

	db "github.com/forbole/bdjuno/v4/database/types"
)

const (
	tableMsgSend      = "msg_send"
	tableMsgMultiSend = "msg_multi_send"
)

type (
	// msgSend represents a single row inside the 'msg_send' table (used only for SELECT)
	msgSend struct {
		ID          uint64         `db:"id"`
		TxHash      string         `db:"tx_hash"`
		FromAddress string         `db:"from_address"`
		ToAddress   string         `db:"to_address"`
		Amount      pq.StringArray `db:"amount"`
	}

	// msgSend represents a single row inside the 'msg_multi_send' table (used only for SELECT)
	msgMultiSend struct {
		ID      uint64         `db:"id"`
		TxHash  string         `db:"tx_hash"`
		Inputs  pq.StringArray `db:"inputs"`
		Outputs pq.StringArray `db:"outputs"`
	}
)

// toGetMsgSendDomain - mapping func to a domain model.
func toGetMsgSendDomain(m msgSend) (bank.MsgSend, error) {
	amount, err := db.FromPqStringArrayToCoins(m.Amount)
	if err != nil {
		return bank.MsgSend{}, err
	}

	return bank.MsgSend{
		FromAddress: m.FromAddress,
		ToAddress:   m.ToAddress,
		Amount:      amount,
	}, nil
}

// toGetMsgSendDomainList - mapping func to a domain list.
func toGetMsgSendDomainList(m []msgSend) ([]bank.MsgSend, error) {
	res := make([]bank.MsgSend, 0, len(m))
	for _, msg := range m {
		r, err := toGetMsgSendDomain(msg)
		if err != nil {
			return nil, err
		}

		res = append(res, r)
	}

	return res, nil
}

// toGetMsgMultiSendDomain - mapping func to a domain model.
func toGetMsgMultiSendDomain(m msgMultiSend) (bank.MsgMultiSend, error) {
	inputs, err := db.FromPqStringArrayToInputs(m.Inputs)
	if err != nil {
		return bank.MsgMultiSend{}, err
	}

	outputs, err := db.FromPqStringArrayToOutputs(m.Outputs)
	if err != nil {
		return bank.MsgMultiSend{}, err
	}

	return bank.MsgMultiSend{
		Inputs:  inputs,
		Outputs: outputs,
	}, nil
}

// toGetMsgMultiSendDomainList - mapping func to a domain list.
func toGetMsgMultiSendDomainList(m []msgMultiSend) ([]bank.MsgMultiSend, error) {
	res := make([]bank.MsgMultiSend, 0, len(m))
	for _, msg := range m {
		r, err := toGetMsgMultiSendDomain(msg)
		if err != nil {
			return nil, err
		}

		res = append(res, r)
	}

	return res, nil
}
