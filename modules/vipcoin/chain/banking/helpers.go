package banking

import (
	"encoding/json"
	"errors"
	"strconv"

	banking "git.ooo.ua/vipcoin/chain/x/banking/types"
	wallets "git.ooo.ua/vipcoin/chain/x/wallets/types"
	"git.ooo.ua/vipcoin/lib/filter"
	sdk "github.com/cosmos/cosmos-sdk/types"
	juno "github.com/forbole/juno/v2/types"

	dbtypes "github.com/forbole/bdjuno/v2/database/types"
)

// getIssueFromTx - get issue from tx
func getIssueFromTx(tx *juno.Tx, msg *banking.MsgIssue) (*banking.Issue, error) {
	for _, log := range tx.TxResponse.Logs {
		for _, event := range log.Events {
			if event.Type != "vipcoin.chain.banking.Issue" {
				continue
			}

			issue, err := getIssueFromAttribute(event.Attributes)
			if err != nil {
				return nil, err
			}

			if issue.Wallet != msg.Wallet {
				continue
			}

			if issue.Amount != msg.Amount {
				continue
			}

			issue.Extras = msg.Extras

			return issue, nil
		}
	}

	return nil, errors.New("not found")
}

// getSystemTransferFromTx - get system transfer from tx
func getSystemTransferFromTx(tx *juno.Tx, msg *banking.MsgSystemTransfer) (*banking.SystemTransfer, error) {
	for _, log := range tx.TxResponse.Logs {
		for _, event := range log.Events {
			if event.Type != "vipcoin.chain.banking.SystemTransfer" {
				continue
			}

			systemTransfer, err := getSystemTransferFromAttribute(event.Attributes)
			if err != nil {
				return nil, err
			}

			if systemTransfer.WalletFrom != msg.WalletFrom {
				continue
			}

			if systemTransfer.WalletTo != msg.WalletTo {
				continue
			}

			if systemTransfer.Amount != msg.Amount {
				continue
			}

			systemTransfer.Extras = msg.Extras

			return systemTransfer, nil
		}
	}

	return nil, errors.New("not found")
}

// getWithdrawFromTx - get withdraw from tx
func getWithdrawFromTx(tx *juno.Tx, msg *banking.MsgWithdraw) (*banking.Withdraw, error) {
	for _, log := range tx.TxResponse.Logs {
		for _, event := range log.Events {
			if event.Type != "vipcoin.chain.banking.Withdraw" {
				continue
			}

			withdraw, err := getWithdrawFromAttribute(event.Attributes)
			if err != nil {
				return nil, err
			}

			if withdraw.Wallet != msg.Wallet {
				continue
			}

			if withdraw.Amount != msg.Amount {
				continue
			}

			withdraw.Extras = msg.Extras

			return withdraw, nil
		}
	}

	return nil, errors.New("not found")
}

// getPaymentFromTx - get payment from tx
func getPaymentFromTx(tx *juno.Tx, msg *banking.MsgPayment) (*banking.Payment, error) {
	for _, log := range tx.TxResponse.Logs {
		for _, event := range log.Events {
			if event.Type != "vipcoin.chain.banking.Payment" {
				continue
			}

			payment, err := getPaymentTransferFromAttribute(event.Attributes)
			if err != nil {
				return nil, err
			}

			if payment.WalletFrom != msg.WalletFrom {
				continue
			}

			if payment.WalletTo != msg.WalletTo {
				continue
			}

			if payment.BaseTransfer.Amount != msg.Amount {
				continue
			}

			payment.Extras = msg.Extras

			return payment, nil
		}
	}

	return nil, errors.New("not found")
}

// getSystemTransfers - get system transfers from logs
func (m *Module) getSystemTransfers(
	tx *juno.Tx,
	walletFrom string,
) (
	systemReward *banking.SystemTransfer,
	systemRefReward *banking.SystemTransfer,
	err error,
) {
	for _, log := range tx.TxResponse.Logs {
		for _, event := range log.Events {
			if event.Type != "vipcoin.chain.banking.SystemTransfer" {
				continue
			}

			systemReward, systemRefReward, err = m.getSystemTransfersFromAttribute(event.Attributes)
			if err != nil {
				return nil, nil, err
			}

			if systemReward.WalletFrom != walletFrom {
				continue
			}

			if systemRefReward.WalletFrom != walletFrom {
				continue
			}

			return systemReward, systemRefReward, nil
		}

	}

	return nil, nil, errors.New("not found")
}

// getSystemTransfersFromAttribute - get system transfers from attributes
func (m *Module) getSystemTransfersFromAttribute(
	attributes []sdk.Attribute,
) (
	systemReward *banking.SystemTransfer,
	systemRefReward *banking.SystemTransfer,
	err error,
) {
	baseString := getAttributeValuesWithKey(attributes, "base_transfer")
	if len(baseString) == 0 {
		return nil, nil, errors.New("base_transfer not found in log")
	}

	var baseVipcoinList []banking.BaseTransfer
	for _, baseString := range baseString {
		var base baseTransfer
		if err := json.Unmarshal([]byte(baseString), &base); err != nil {
			return nil, nil, err
		}

		baseVipcoin, err := base.toVipcoinBaseTransfer()
		if err != nil {
			return nil, nil, err
		}

		baseVipcoinList = append(baseVipcoinList, baseVipcoin)
	}

	wallet := getAttributeValueWithKey(attributes, "wallet_from")
	if wallet == "" {
		return nil, nil, errors.New("wallet_from not found in log")
	}

	var walletFrom string
	if err := json.Unmarshal([]byte(wallet), &walletFrom); err != nil {
		return nil, nil, err
	}

	walletArr := getAttributeValuesWithKey(attributes, "wallet_to")
	if len(walletArr) == 0 {
		return nil, nil, errors.New("wallet_to not found in log")
	}

	var walletToList []string
	for _, walletString := range walletArr {
		var walletTo string
		if err := json.Unmarshal([]byte(walletString), &walletTo); err != nil {
			return nil, nil, err
		}

		walletToList = append(walletToList, walletTo)
	}

	walletsSystemReward, err := m.walletsRepo.GetWallets(
		filter.NewFilter().
			SetArgument(dbtypes.FieldKind, wallets.WALLET_KIND_SYSTEM_REWARD))
	switch {
	case err != nil:
		return nil, nil, err
	case len(walletsSystemReward) != 1:
		return nil, nil, wallets.ErrInvalidKindField
	}

	for _, walletTo := range walletToList {
		if walletTo == walletsSystemReward[0].Address {
			systemReward = &banking.SystemTransfer{
				WalletFrom: walletFrom,
				WalletTo:   walletTo,
			}

			systemReward.BaseTransfer, err = getTransfersWithKind(baseVipcoinList, banking.TRANSFER_KIND_SYSTEM_REWARD)
			if err != nil {
				return nil, nil, err
			}

			continue
		}

		systemRefReward = &banking.SystemTransfer{
			WalletFrom: walletFrom,
			WalletTo:   walletTo,
		}

		systemRefReward.BaseTransfer, err = getTransfersWithKind(baseVipcoinList, banking.TRANSFER_KIND_SYSTEM_REF_REWARD)
		if err != nil {
			return nil, nil, err
		}
	}

	return systemReward, systemRefReward, nil
}

// getPaymentTransferFromAttribute - get payment transfer by kind from attributes
func getPaymentTransferFromAttribute(attributes []sdk.Attribute) (*banking.Payment, error) {
	var base baseTransfer
	if err := json.Unmarshal([]byte(getAttributeValueWithKey(attributes, "base_transfer")), &base); err != nil {
		return nil, err
	}

	var walletFrom string
	if err := json.Unmarshal([]byte(getAttributeValueWithKey(attributes, "wallet_from")), &walletFrom); err != nil {
		return nil, err
	}

	var walletTo string
	if err := json.Unmarshal([]byte(getAttributeValueWithKey(attributes, "wallet_to")), &walletTo); err != nil {
		return nil, err
	}

	var fee string
	if err := json.Unmarshal([]byte(getAttributeValueWithKey(attributes, "fee")), &fee); err != nil {
		return nil, err
	}

	feeUint, err := strconv.ParseUint(fee, 10, 64)
	if err != nil {
		return nil, err
	}

	baseVipcoin, err := base.toVipcoinBaseTransfer()
	if err != nil {
		return nil, err
	}

	paymentTransfer := banking.Payment{
		BaseTransfer: baseVipcoin,
		WalletFrom:   walletFrom,
		WalletTo:     walletTo,
		Fee:          feeUint,
	}

	return &paymentTransfer, nil
}

// getIssueFromAttribute - get issue from attributes
func getIssueFromAttribute(attributes []sdk.Attribute) (*banking.Issue, error) {
	var base baseTransfer
	if err := json.Unmarshal([]byte(getAttributeValueWithKey(attributes, "base_transfer")), &base); err != nil {
		return nil, err
	}

	var wallet string
	if err := json.Unmarshal([]byte(getAttributeValueWithKey(attributes, "wallet")), &wallet); err != nil {
		return nil, err
	}

	baseVipcoin, err := base.toVipcoinBaseTransfer()
	if err != nil {
		return nil, err
	}

	issue := banking.Issue{
		BaseTransfer: baseVipcoin,
		Wallet:       wallet,
	}

	return &issue, nil
}

// getSystemTransferFromAttribute - get system transfer from attributes
func getSystemTransferFromAttribute(attributes []sdk.Attribute) (*banking.SystemTransfer, error) {
	var base baseTransfer
	if err := json.Unmarshal([]byte(getAttributeValueWithKey(attributes, "base_transfer")), &base); err != nil {
		return nil, err
	}

	var walletTo string
	if err := json.Unmarshal([]byte(getAttributeValueWithKey(attributes, "wallet_to")), &walletTo); err != nil {
		return nil, err
	}

	var walletFrom string
	if err := json.Unmarshal([]byte(getAttributeValueWithKey(attributes, "wallet_from")), &walletFrom); err != nil {
		return nil, err
	}

	baseVipcoin, err := base.toVipcoinBaseTransfer()
	if err != nil {
		return nil, err
	}

	systemTransfer := banking.SystemTransfer{
		BaseTransfer: baseVipcoin,
		WalletFrom:   walletFrom,
		WalletTo:     walletTo,
	}

	return &systemTransfer, nil
}

// getWithdrawFromAttribute - get withdraw from attributes
func getWithdrawFromAttribute(attributes []sdk.Attribute) (*banking.Withdraw, error) {
	var base baseTransfer
	if err := json.Unmarshal([]byte(getAttributeValueWithKey(attributes, "base_transfer")), &base); err != nil {
		return nil, err
	}

	var wallet string
	if err := json.Unmarshal([]byte(getAttributeValueWithKey(attributes, "wallet")), &wallet); err != nil {
		return nil, err
	}

	baseVipcoin, err := base.toVipcoinBaseTransfer()
	if err != nil {
		return nil, err
	}

	issue := banking.Withdraw{
		BaseTransfer: baseVipcoin,
		Wallet:       wallet,
	}

	return &issue, nil
}

// getAttributeValueWithKey - get attribute value with key
func getAttributeValueWithKey(attributes []sdk.Attribute, key string) string {
	for _, attribute := range attributes {
		if attribute.Key == key {
			return attribute.Value
		}
	}

	return ""
}

// getAttributeValuesWithKey - get attribute values with key
func getAttributeValuesWithKey(attributes []sdk.Attribute, key string) []string {
	var result []string
	for _, attribute := range attributes {
		if attribute.Key == key {
			result = append(result, attribute.Value)
		}
	}

	return result
}

// getTransfersWithKind - get transfers with kind
func getTransfersWithKind(transfers []banking.BaseTransfer, kind banking.TransferKind) (banking.BaseTransfer, error) {
	for _, transfer := range transfers {
		if transfer.Kind == kind {
			return transfer, nil
		}
	}

	return banking.BaseTransfer{}, errors.New("not found")
}
