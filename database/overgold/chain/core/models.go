package core

import (
	"strconv"

	"git.ooo.ua/vipcoin/lib/errs"
	"git.ooo.ua/vipcoin/ovg-chain/x/core/types"

	db "github.com/forbole/bdjuno/v4/database/types"
)

// BLOCK MsgIssue

// toMsgIssueDomain - mapping func to a domain model.
func toMsgIssueDomain(m db.CoreMsgIssue) types.MsgIssue {
	return types.MsgIssue{
		Creator: m.Creator,
		Amount:  strconv.FormatUint(m.Amount, 10),
		Denom:   m.Denom,
		Address: m.Address,
	}
}

// toMsgIssueDomainList - mapping func to a domain list.
func toMsgIssueDomainList(m []db.CoreMsgIssue) []types.MsgIssue {
	res := make([]types.MsgIssue, 0, len(m))
	for _, msg := range m {
		res = append(res, toMsgIssueDomain(msg))
	}

	return res
}

// toMsgIssueDatabase - mapping func to a database model.
func toMsgIssueDatabase(hash string, m types.MsgIssue) (db.CoreMsgIssue, error) {
	amount, err := strconv.ParseUint(m.Amount, 10, 64)
	if err != nil {
		return db.CoreMsgIssue{}, errs.Internal{Cause: err.Error()}
	}

	return db.CoreMsgIssue{
		TxHash:  hash,
		Creator: m.Creator,
		Amount:  amount,
		Denom:   m.Denom,
		Address: m.Address,
	}, nil
}

// BLOCK MsgWithdraw

// toMsgWithdrawDomain - mapping func to a domain model.
func toMsgWithdrawDomain(m db.CoreMsgWithdraw) types.MsgWithdraw {
	return types.MsgWithdraw{
		Creator: m.Creator,
		Amount:  strconv.FormatUint(m.Amount, 10),
		Denom:   m.Denom,
		Address: m.Address,
	}
}

// toMsgWithdrawDomainList - mapping func to a domain list.
func toMsgWithdrawDomainList(m []db.CoreMsgWithdraw) []types.MsgWithdraw {
	res := make([]types.MsgWithdraw, 0, len(m))
	for _, msg := range m {
		res = append(res, toMsgWithdrawDomain(msg))
	}

	return res
}

// toMsgWithdrawDatabase - mapping func to a database model.
func toMsgWithdrawDatabase(hash string, m types.MsgWithdraw) (db.CoreMsgWithdraw, error) {
	amount, err := strconv.ParseUint(m.Amount, 10, 64)
	if err != nil {
		return db.CoreMsgWithdraw{}, errs.Internal{Cause: err.Error()}
	}

	return db.CoreMsgWithdraw{
		TxHash:  hash,
		Creator: m.Creator,
		Amount:  amount,
		Denom:   m.Denom,
		Address: m.Address,
	}, nil
}

// BLOCK MsgSend

// toMsgSendDomain - mapping func to a domain model.
func toMsgSendDomain(m db.CoreMsgSend) types.MsgSend {
	return types.MsgSend{
		Creator: m.Creator,
		From:    m.AddressFrom,
		To:      m.AddressTo,
		Amount:  strconv.FormatUint(m.Amount, 10),
		Denom:   m.Denom,
	}
}

// toMsgSendDomainList - mapping func to a domain list.
func toMsgSendDomainList(m []db.CoreMsgSend) []types.MsgSend {
	res := make([]types.MsgSend, 0, len(m))
	for _, msg := range m {
		res = append(res, toMsgSendDomain(msg))
	}

	return res
}

// toMsgSendDatabase - mapping func to a database model.
func toMsgSendDatabase(hash string, m types.MsgSend) (db.CoreMsgSend, error) {
	amount, err := strconv.ParseUint(m.Amount, 10, 64)
	if err != nil {
		return db.CoreMsgSend{}, errs.Internal{Cause: err.Error()}
	}

	return db.CoreMsgSend{
		TxHash:      hash,
		Creator:     m.Creator,
		AddressFrom: m.From,
		AddressTo:   m.To,
		Amount:      amount,
		Denom:       m.Denom,
	}, nil
}
