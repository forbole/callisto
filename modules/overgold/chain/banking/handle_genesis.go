package banking

import (
	"encoding/json"

	bankingtypes "git.ooo.ua/vipcoin/chain/x/banking/types"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/rs/zerolog/log"
	tmtypes "github.com/tendermint/tendermint/types"
)

// HandleGenesis implements GenesisModule
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	log.Debug().Str("module", "banking").Msg("parsing genesis")

	var bankingState bankingtypes.GenesisState
	if err := m.cdc.UnmarshalJSON(appState[bankingtypes.ModuleName], &bankingState); err != nil {
		return err
	}

	transfers, err := m.UnpackTransfers(bankingState.Transfers...)
	if err != nil {
		return err
	}

	var systems []*bankingtypes.SystemTransfer
	var payments []*bankingtypes.Payment
	var withdraws []*bankingtypes.Withdraw
	var issues []*bankingtypes.Issue

	for _, transfer := range transfers {
		switch transfer.GetBase().Kind {
		case bankingtypes.TRANSFER_KIND_SYSTEM:
			system := &bankingtypes.SystemTransfer{BaseTransfer: transfer.GetBase()}
			if transferWalletsI, ok := transfer.(bankingtypes.TransferWalletsI); ok {
				system.WalletFrom = transferWalletsI.GetWalletFrom()
				system.WalletTo = transferWalletsI.GetWalletTo()
			}
			systems = append(systems, system)

		case bankingtypes.TRANSFER_KIND_PAYMENT:
			payment := &bankingtypes.Payment{BaseTransfer: transfer.GetBase()}
			if transferWalletsI, ok := transfer.(bankingtypes.TransferWalletsI); ok {
				payment.WalletFrom = transferWalletsI.GetWalletFrom()
				payment.WalletTo = transferWalletsI.GetWalletTo()
			}

			if transferWalletsI, ok := transfer.(bankingtypes.TransferFeeI); ok {
				payment.Fee = transferWalletsI.GetFee()
			}
			payments = append(payments, payment)

		case bankingtypes.TRANSFER_KIND_WITHDRAW:
			withdraw := &bankingtypes.Withdraw{BaseTransfer: transfer.GetBase()}
			if transferWalletI, ok := transfer.(bankingtypes.TransferWalletI); ok {
				withdraw.Wallet = transferWalletI.GetWallet()
			}
			withdraws = append(withdraws, withdraw)

		case bankingtypes.TRANSFER_KIND_ISSUE:
			issue := &bankingtypes.Issue{BaseTransfer: transfer.GetBase()}
			if transferWalletI, ok := transfer.(bankingtypes.TransferWalletI); ok {
				issue.Wallet = transferWalletI.GetWallet()
			}
			issues = append(issues, issue)
		}
	}

	if err := m.bankingRepo.SaveSystemTransfers(systems...); err != nil {
		return err
	}

	if err := m.bankingRepo.SavePayments(payments...); err != nil {
		return err
	}

	if err := m.bankingRepo.SaveWithdraws(withdraws...); err != nil {
		return err
	}

	return m.bankingRepo.SaveIssues(issues...)
}

func (m *Module) UnpackTransfer(any *cdctypes.Any) (bankingtypes.TransferI, error) {
	var transfer bankingtypes.TransferI
	if err := m.cdc.UnpackAny(any, &transfer); err != nil {
		return nil, err
	}
	return transfer, nil
}

func (m *Module) UnpackTransfers(any ...*cdctypes.Any) ([]bankingtypes.TransferI, error) {
	transfers := make([]bankingtypes.TransferI, 0, len(any))
	for _, a := range any {
		transfer, err := m.UnpackTransfer(a)
		if err != nil {
			return nil, err
		}

		transfers = append(transfers, transfer)
	}

	return transfers, nil
}
