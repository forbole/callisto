package wallets

import (
	typeswallets "git.ooo.ua/vipcoin/chain/x/wallets/types"
	"git.ooo.ua/vipcoin/lib/filter"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/forbole/juno/v2/types"
	juno "github.com/forbole/juno/v2/types"

	dbtypes "github.com/forbole/bdjuno/v2/database/types"
)

// HandleMsg implements MessageModule
func (m *Module) HandleMsg(index int, msg sdk.Msg, tx *types.Tx) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch walletMsg := msg.(type) {
	case *typeswallets.MsgSetWalletState:
		return m.handleMsgSetStates(walletMsg)
	case *typeswallets.MsgCreateWallet:
		return m.handleMsgCreateWallet(tx, index, walletMsg)
	case *typeswallets.MsgSetDefaultWallet:
		return m.handleMsgSetDefaultWallet(walletMsg)
	case *typeswallets.MsgSetExtra:
		return m.handleMsgSetExtra(walletMsg)
	case *typeswallets.MsgCreateWalletWithBalance:
		return m.MsgCreateWalletWithBalance(walletMsg)
	default:
		return nil
	}
}

// handleMsgCreateWallet allows to properly handle a handleMsgCreateWallet
func (m *Module) handleMsgCreateWallet(tx *juno.Tx, index int, msg *typeswallets.MsgCreateWallet) error {
	if err := m.walletsRepo.SaveCreateWallet(msg); err != nil {
		return err
	}

	newWallet := typeswallets.Wallet{
		Address:        msg.Address,
		AccountAddress: msg.AccountAddress,
		Kind:           msg.Kind,
		State:          msg.State,
		Extras:         msg.Extras,
	}

	return m.walletsRepo.SaveWallets(&newWallet)
}

// handleMsgSetStates allows to properly handle a MsgSetState
func (m *Module) handleMsgSetStates(msg *typeswallets.MsgSetWalletState) error {
	if err := m.walletsRepo.SaveStates(msg); err != nil {
		return err
	}

	wallets, err := m.walletsRepo.GetWallets(filter.NewFilter().SetArgument(dbtypes.FieldAddress, msg.Address))
	if err != nil {
		return err
	}

	if len(wallets) != 1 {
		return typeswallets.ErrInvalidAddressField
	}

	wallets[0].State = msg.State

	return m.walletsRepo.SaveWallets(wallets...)
}

// handleMsgSetDefaultWallet allows to properly handle a MsgSetDefaultWallet
func (m *Module) handleMsgSetDefaultWallet(msg *typeswallets.MsgSetDefaultWallet) error {
	if err := m.walletsRepo.SaveDefaultWallets(msg); err != nil {
		return err
	}

	// TODO: look at chain/x/wallets/keeper/msg_server_set_default_wallet.go

	address, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return typeswallets.ErrInvalidAddressField
	}

	targetWallet, err := m.walletsRepo.GetWallets(filter.NewFilter().SetArgument(dbtypes.FieldAddress, address))
	if err != nil {
		return err
	}

	if targetWallet[0].Address == "" {
		return typeswallets.ErrNotFoundWallet
	}

	account, err := m.accountsRepo.GetAccounts(filter.NewFilter().SetArgument(dbtypes.FieldAddress, targetWallet[0].AccountAddress))
	if err != nil {
		return err
	}

	for _, walletAddr := range account[0].Wallets {
		w, err := m.walletsRepo.GetWallets(filter.NewFilter().SetArgument(dbtypes.FieldAddress, walletAddr))
		if err != nil {
			return err
		}
		if !w[0].Default {
			continue
		}
		w[0].Default = false
		if err := m.walletsRepo.SaveWallets(w[0]); err != nil {
			return err
		}
		break
	}

	targetWallet[0].Default = true

	return m.walletsRepo.SaveWallets(targetWallet[0])
}

// handleMsgSetExtra allows to properly handle a MsgSetExtra
func (m *Module) handleMsgSetExtra(msg *typeswallets.MsgSetExtra) error {
	if err := m.walletsRepo.SaveExtras(msg); err != nil {
		return err
	}

	wallets, err := m.walletsRepo.GetWallets(filter.NewFilter().SetArgument(dbtypes.FieldAddress, msg.Address))
	if err != nil {
		return err
	}

	if len(wallets) != 1 {
		return typeswallets.ErrInvalidAddressField
	}

	wallets[0].Extras = msg.Extras

	return m.walletsRepo.SaveWallets(wallets...)
}

// MsgCreateWalletWithBalance allows to properly handle a handleMsgCreateWallet
func (m *Module) MsgCreateWalletWithBalance(msg *typeswallets.MsgCreateWalletWithBalance) error {
	if err := m.walletsRepo.SaveCreateWalletWithBalance(msg); err != nil {
		return err
	}

	wallets, err := m.walletsRepo.GetWallets(filter.NewFilter().SetArgument(dbtypes.FieldAddress, msg.Address))
	if err != nil {
		return err
	}

	if len(wallets) != 1 {
		return typeswallets.ErrInvalidAddressField
	}

	wallets[0].AccountAddress = msg.AccountAddress
	wallets[0].Kind = msg.Kind
	wallets[0].State = msg.State
	wallets[0].Balance = msg.Balance
	wallets[0].Extras = msg.Extras
	wallets[0].Default = msg.Default
	wallets[0].Balance = msg.Balance

	return m.walletsRepo.SaveWallets(wallets...)
}
