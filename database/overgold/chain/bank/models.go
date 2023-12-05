package bank

import (
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"

	db "github.com/forbole/bdjuno/v4/database/types"
)

// BLOCK MsgSend

// toMsgSendDomain - mapping func to a domain model.
func toMsgSendDomain(m db.MsgSend) bank.MsgSend {
	return bank.MsgSend{
		FromAddress: m.FromAddress,
		ToAddress:   m.ToAddress,
		Amount:      m.Amount.ToCoins(),
	}
}

// toMsgSendDomainList - mapping func to a domain list.
func toMsgSendDomainList(m []db.MsgSend) []bank.MsgSend {
	res := make([]bank.MsgSend, 0, len(m))
	for _, msg := range m {
		res = append(res, toMsgSendDomain(msg))
	}

	return res
}

// toMsgSendDatabase - mapping func to a database model.
func toMsgSendDatabase(hash string, m bank.MsgSend) db.MsgSend {
	return db.MsgSend{
		TxHash:      hash,
		FromAddress: m.FromAddress,
		ToAddress:   m.ToAddress,
		Amount:      db.NewDbCoins(m.Amount),
	}
}

// BLOCK MsgMultiSend

// toMsgMultiSendDomain - mapping func to a domain model.
func toMsgMultiSendDomain(m db.MsgMultiSend) bank.MsgMultiSend {
	return bank.MsgMultiSend{
		Inputs:  m.Inputs.ToInputList(),
		Outputs: m.Ouputs.ToOutputList(),
	}
}

// toMsgMultiSendDomainList - mapping func to a domain list.
func toMsgMultiSendDomainList(m []db.MsgMultiSend) []bank.MsgMultiSend {
	res := make([]bank.MsgMultiSend, 0, len(m))
	for _, msg := range m {
		res = append(res, toMsgMultiSendDomain(msg))
	}

	return res
}

// toMsgMultiSendDatabase - mapping func to a database model.
func toMsgMultiSendDatabase(hash string, m bank.MsgMultiSend) db.MsgMultiSend {
	return db.MsgMultiSend{
		TxHash: hash,
		Inputs: db.NewDbSendDataListByInputs(m.Inputs),
		Ouputs: db.NewDbSendDataListByOutputs(m.Outputs),
	}
}
