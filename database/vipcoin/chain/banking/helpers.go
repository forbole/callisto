package banking

import (
	"time"

	bankingtypes "git.ooo.ua/vipcoin/chain/x/banking/types"
	extratypes "git.ooo.ua/vipcoin/chain/x/types"

	"github.com/forbole/bdjuno/v2/database/types"
)

const (
	tableTransfers         = "vipcoin_chain_banking_base_transfers"
	tableSystemTransfer    = "vipcoin_chain_banking_system_transfer"
	tablePayment           = "vipcoin_chain_banking_payment"
	tableWithdraw          = "vipcoin_chain_banking_withdraw"
	tableIssue             = "vipcoin_chain_banking_issue"
	tableMsgPayment        = "vipcoin_chain_banking_msg_payment"
	tableMsgSystemTransfer = "vipcoin_chain_banking_msg_system_transfer"
	tableMsgIssue          = "vipcoin_chain_banking_msg_issue"
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
	for index := range extras.Extras {
		result = append(result, &extras.Extras[index])
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
		Timestamp: time.Unix(transfer.Timestamp, 0).UTC(),
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

// toSystemTransfersDatabase - mapping func to database model
func toSystemTransfersDatabase(transfers ...*bankingtypes.SystemTransfer) []types.DBSystemTransfer {
	result := make([]types.DBSystemTransfer, 0, len(transfers))
	for _, t := range transfers {
		result = append(result, toSystemTransferDatabase(t))
	}

	return result

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

// toIssuesDatabase - mapping func to database model
func toIssuesDatabase(issues ...*bankingtypes.Issue) []types.DBIssue {
	result := make([]types.DBIssue, 0, len(issues))
	for _, issue := range issues {
		result = append(result, toIssueDatabase(issue))
	}

	return result
}

// toTransfersDatabase - mapping func to database model
func toTransfersDatabase(transfers ...*bankingtypes.BaseTransfer) []types.DBTransfer {
	result := make([]types.DBTransfer, 0, len(transfers))
	for _, transfer := range transfers {
		result = append(result, toTransferDatabase(transfer))
	}

	return result
}

// toPaymentDatabase - mapping func to database model
func toMsgPaymentDatabase(payments *bankingtypes.MsgPayment) types.DBMsgPayment {
	return types.DBMsgPayment{
		Creator:    payments.Creator,
		WalletFrom: payments.WalletFrom,
		WalletTo:   payments.WalletTo,
		Asset:      payments.Asset,
		Amount:     payments.Amount,
		Extras:     toExtrasDB(payments.Extras),
	}
}

// toPaymentsDatabase - mapping func to database model
func toMsgPaymentsDatabase(payments ...*bankingtypes.MsgPayment) []types.DBMsgPayment {
	result := make([]types.DBMsgPayment, 0, len(payments))
	for _, payment := range payments {
		result = append(result, toMsgPaymentDatabase(payment))
	}

	return result
}

// toPaymentDomain - mapping func to domain model
func toMsgPaymentDomain(payments types.DBMsgPayment) *bankingtypes.MsgPayment {
	return &bankingtypes.MsgPayment{
		Creator:    payments.Creator,
		WalletFrom: payments.WalletFrom,
		WalletTo:   payments.WalletTo,
		Asset:      payments.Asset,
		Amount:     payments.Amount,
		Extras:     fromExtrasDB(payments.Extras),
	}
}

// toSystemTransferDatabase - mapping func to database model
func toMsgSystemTransferDatabase(transfer *bankingtypes.MsgSystemTransfer) types.DBMsgSystemTransfer {
	return types.DBMsgSystemTransfer{
		Creator:    transfer.Creator,
		WalletFrom: transfer.WalletFrom,
		WalletTo:   transfer.WalletTo,
		Asset:      transfer.Asset,
		Amount:     transfer.Amount,
		Extras:     toExtrasDB(transfer.Extras),
	}
}

// toSystemTransfersDatabase - mapping func to database model
func toMsgSystemTransfersDatabase(transfers ...*bankingtypes.MsgSystemTransfer) []types.DBMsgSystemTransfer {
	result := make([]types.DBMsgSystemTransfer, 0, len(transfers))
	for _, transfer := range transfers {
		result = append(result, toMsgSystemTransferDatabase(transfer))
	}

	return result
}

// toSystemTransferDomain - mapping func to domain model
func toMsgSystemTransferDomain(transfer types.DBMsgSystemTransfer) *bankingtypes.MsgSystemTransfer {
	return &bankingtypes.MsgSystemTransfer{
		Creator:    transfer.Creator,
		WalletFrom: transfer.WalletFrom,
		WalletTo:   transfer.WalletTo,
		Asset:      transfer.Asset,
		Amount:     transfer.Amount,
		Extras:     fromExtrasDB(transfer.Extras),
	}
}

// toMsgIssueDatabase - mapping func to database model
func toMsgIssueDatabase(issue *bankingtypes.MsgIssue) types.DBMsgIssue {
	return types.DBMsgIssue{
		Creator: issue.Creator,
		Wallet:  issue.Wallet,
		Asset:   issue.Asset,
		Amount:  issue.Amount,
		Extras:  toExtrasDB(issue.Extras),
	}
}

// toMsgIssuesDatabase - mapping func to database model
func toMsgIssuesDatabase(issues ...*bankingtypes.MsgIssue) []types.DBMsgIssue {
	result := make([]types.DBMsgIssue, 0, len(issues))
	for _, issue := range issues {
		result = append(result, toMsgIssueDatabase(issue))
	}

	return result
}

// toMsgIssueDomain - mapping func to domain model
func toMsgIssueDomain(issue types.DBMsgIssue) *bankingtypes.MsgIssue {
	return &bankingtypes.MsgIssue{
		Creator: issue.Creator,
		Wallet:  issue.Wallet,
		Asset:   issue.Asset,
		Amount:  issue.Amount,
		Extras:  fromExtrasDB(issue.Extras),
	}
}
