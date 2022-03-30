/*
 * Copyright 2022 Business Process Technologies. All rights reserved.
 */

package banking

import (
	"time"

	bankingtypes "git.ooo.ua/vipcoin/chain/x/banking/types"
	extratypes "git.ooo.ua/vipcoin/chain/x/types"

	"github.com/forbole/bdjuno/v2/database/types"
)

// toExtrasDB - mapping func to database model
func toExtrasDB(extras []*extratypes.Extra) types.ExtraDB {
	result := make([]extratypes.Extra, 0, len(extras))
	for _, extra := range extras {
		result = append(result, *extra)
	}

	return types.ExtraDB{Extras: result}
}

// fromExtrasDB - mapping func from database model
func fromExtrasDB(extras types.ExtraDB) []*extratypes.Extra {
	result := make([]*extratypes.Extra, 0, len(extras.Extras))
	for _, extra := range extras.Extras {
		result = append(result, &extra)
	}

	return result
}

// toTransferDomain - mapping func to domain model
func toTransferDomain(transfer types.DBTransfer) *bankingtypes.BaseTransfer {
	return &bankingtypes.BaseTransfer{
		Id:        transfer.ID,
		Asset:     transfer.Asset,
		Amount:    transfer.Amount,
		Kind:      bankingtypes.TransferKind(transfer.Kind),
		Extras:    fromExtrasDB(transfer.Extras),
		Timestamp: transfer.Timestamp.Unix(),
		TxHash:    transfer.TxHash,
	}
}

// toTransferDatabase - mapping func to database model
func toTransferDatabase(transfer *bankingtypes.BaseTransfer) types.DBTransfer {
	return types.DBTransfer{
		ID:        transfer.Id,
		Asset:     transfer.Asset,
		Amount:    transfer.Amount,
		Kind:      int32(transfer.Kind),
		Extras:    toExtrasDB(transfer.Extras),
		Timestamp: time.Unix(transfer.Timestamp, 0),
		TxHash:    transfer.TxHash,
	}
}

// toPaymentDomain - mapping func to domain model
func toPaymentDomain(payment types.DBPayment) *bankingtypes.Payment {
	return &bankingtypes.Payment{
		BaseTransfer: bankingtypes.BaseTransfer{
			Id:        payment.ID,
			Asset:     payment.Asset,
			Amount:    payment.Amount,
			Kind:      bankingtypes.TransferKind(payment.Kind),
			Extras:    fromExtrasDB(payment.Extras),
			Timestamp: payment.Timestamp.Unix(),
			TxHash:    payment.TxHash,
		},
		WalletFrom: payment.WalletFrom,
		WalletTo:   payment.WalletTo,
		Fee:        payment.Fee,
	}
}

// toPaymentDatabase - mapping func to database model
func toPaymentDatabase(payment *bankingtypes.Payment) types.DBPayment {
	return types.DBPayment{
		DBTransfer: types.DBTransfer{
			ID:        payment.Id,
			Asset:     payment.Asset,
			Amount:    payment.Amount,
			Kind:      int32(payment.Kind),
			Extras:    toExtrasDB(payment.Extras),
			Timestamp: time.Unix(payment.Timestamp, 0),
			TxHash:    payment.TxHash,
		},
		WalletFrom: payment.WalletFrom,
		WalletTo:   payment.WalletTo,
		Fee:        payment.Fee,
	}
}

// toSystemTransferDomain - mapping func to domain model
func toSystemTransferDomain(transfer types.DBSystemTransfer) *bankingtypes.SystemTransfer {
	return &bankingtypes.SystemTransfer{
		BaseTransfer: bankingtypes.BaseTransfer{
			Id:        transfer.ID,
			Asset:     transfer.Asset,
			Amount:    transfer.Amount,
			Kind:      bankingtypes.TransferKind(transfer.Kind),
			Extras:    fromExtrasDB(transfer.Extras),
			Timestamp: transfer.Timestamp.Unix(),
			TxHash:    transfer.TxHash,
		},
		WalletFrom: transfer.WalletFrom,
		WalletTo:   transfer.WalletTo,
	}
}

// toSystemTransferDatabase - mapping func to database model
func toSystemTransferDatabase(transfer *bankingtypes.SystemTransfer) types.DBSystemTransfer {
	return types.DBSystemTransfer{
		DBTransfer: types.DBTransfer{
			ID:        transfer.Id,
			Asset:     transfer.Asset,
			Amount:    transfer.Amount,
			Kind:      int32(transfer.Kind),
			Extras:    toExtrasDB(transfer.Extras),
			Timestamp: time.Unix(transfer.Timestamp, 0),
			TxHash:    transfer.TxHash,
		},
		WalletFrom: transfer.WalletFrom,
		WalletTo:   transfer.WalletTo,
	}
}

// toWithdrawDomain - mapping func to domain model
func toWithdrawDomain(withdraw types.DBWithdraw) *bankingtypes.Withdraw {
	return &bankingtypes.Withdraw{
		BaseTransfer: bankingtypes.BaseTransfer{
			Id:        withdraw.ID,
			Asset:     withdraw.Asset,
			Amount:    withdraw.Amount,
			Kind:      bankingtypes.TransferKind(withdraw.Kind),
			Extras:    fromExtrasDB(withdraw.Extras),
			Timestamp: withdraw.Timestamp.Unix(),
			TxHash:    withdraw.TxHash,
		},
		Wallet: withdraw.Wallet,
	}
}

// toWithdrawDatabase - mapping func to database model
func toWithdrawDatabase(withdraw *bankingtypes.Withdraw) types.DBWithdraw {
	return types.DBWithdraw{
		DBTransfer: types.DBTransfer{
			ID:        withdraw.Id,
			Asset:     withdraw.Asset,
			Amount:    withdraw.Amount,
			Kind:      int32(withdraw.Kind),
			Extras:    toExtrasDB(withdraw.Extras),
			Timestamp: time.Unix(withdraw.Timestamp, 0),
			TxHash:    withdraw.TxHash,
		},
		Wallet: withdraw.Wallet,
	}
}

// toIssueDomain - mapping func to domain model
func toIssueDomain(issue types.DBIssue) *bankingtypes.Issue {
	return &bankingtypes.Issue{
		BaseTransfer: bankingtypes.BaseTransfer{
			Id:        issue.ID,
			Asset:     issue.Asset,
			Amount:    issue.Amount,
			Kind:      bankingtypes.TransferKind(issue.Kind),
			Extras:    fromExtrasDB(issue.Extras),
			Timestamp: issue.Timestamp.Unix(),
			TxHash:    issue.TxHash,
		},
		Wallet: issue.Wallet,
	}
}

// toIssueDatabase - mapping func to database model
func toIssueDatabase(issue *bankingtypes.Issue) types.DBIssue {
	return types.DBIssue{
		DBTransfer: types.DBTransfer{
			ID:        issue.Id,
			Asset:     issue.Asset,
			Amount:    issue.Amount,
			Kind:      int32(issue.Kind),
			Extras:    toExtrasDB(issue.Extras),
			Timestamp: time.Unix(issue.Timestamp, 0),
			TxHash:    issue.TxHash,
		},
		Wallet: issue.Wallet,
	}
}
