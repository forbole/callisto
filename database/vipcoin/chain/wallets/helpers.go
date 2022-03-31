/*
 * Copyright 2022 Business Process Technologies. All rights reserved.
 */

package wallets

import (
	extratypes "git.ooo.ua/vipcoin/chain/x/types"
	walletstypes "git.ooo.ua/vipcoin/chain/x/wallets/types"

	cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"

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

// toBalanceDB - mapping func to database model
func toBalanceDB(balance cosmos_sdk_types.Coins) types.BalanceDB {
	return types.BalanceDB{Balance: balance}
}

// fromBalanceDB - mapping func from database model
func fromBalanceDB(balance types.BalanceDB) cosmos_sdk_types.Coins {
	result := make(cosmos_sdk_types.Coins, 0, len(balance.Balance))
	for _, balance := range balance.Balance {
		result = append(result, balance)
	}

	return result
}

// toWalletsDatabase - mapping func to database model
func toWalletsDatabase(wallets *walletstypes.Wallet) (types.DBWallets, error) {
	return types.DBWallets{
		Address:        wallets.Address,
		AccountAddress: wallets.AccountAddress,
		Kind:           int32(wallets.Kind),
		State:          int32(wallets.State),
		Balance:        toBalanceDB(wallets.Balance),
		Extras:         toExtrasDB(wallets.Extras),
		DefaultStatus:  wallets.Default,
	}, nil
}

// toWalletDomain - mapping func to domain model
func toWalletDomain(wallet types.DBWallets) (*walletstypes.Wallet, error) {
	return &walletstypes.Wallet{
		Address:        wallet.Address,
		AccountAddress: wallet.AccountAddress,
		Kind:           walletstypes.WalletKind(wallet.Kind),
		State:          walletstypes.WalletState(wallet.State),
		Balance:        fromBalanceDB(wallet.Balance),
		Extras:         fromExtrasDB(wallet.Extras),
		Default:        wallet.DefaultStatus,
	}, nil
}

// toCreateWalletDatabase - mapping func to database model
func toCreateWalletDatabase(wallet *walletstypes.MsgCreateWallet) types.DBCreateWallet {
	return types.DBCreateWallet{
		Creator:        wallet.Creator,
		Address:        wallet.Address,
		AccountAddress: wallet.AccountAddress,
		Kind:           int32(wallet.Kind),
		State:          int32(wallet.State),
		Extras:         toExtrasDB(wallet.Extras),
	}
}

// toCreateWalletDomain - mapping func to database model
func toCreateWalletDomain(wallet types.DBCreateWallet) *walletstypes.MsgCreateWallet {
	return &walletstypes.MsgCreateWallet{
		Creator:        wallet.Creator,
		Address:        wallet.Address,
		AccountAddress: wallet.AccountAddress,
		Kind:           walletstypes.WalletKind(wallet.Kind),
		State:          walletstypes.WalletState(wallet.State),
		Extras:         fromExtrasDB(wallet.Extras),
	}
}

// toSetWalletStateDatabase - mapping func to database model
func toSetWalletStateDatabase(wallet *walletstypes.MsgSetWalletState) types.DBSetWalletState {
	return types.DBSetWalletState{
		Creator: wallet.Creator,
		Address: wallet.Address,
		State:   int32(wallet.State),
	}
}

// toSetWalletStateDomain - mapping func to domain model
func toSetWalletStateDomain(wallet types.DBSetWalletState) *walletstypes.MsgSetWalletState {
	return &walletstypes.MsgSetWalletState{
		Creator: wallet.Creator,
		Address: wallet.Address,
		State:   walletstypes.WalletState(wallet.State),
	}
}

// toSetDefaultWalletDatabase - mapping func to database model
func toSetDefaultWalletDatabase(wallet *walletstypes.MsgSetDefaultWallet) types.DBSetDefaultWallet {
	return types.DBSetDefaultWallet{
		Creator: wallet.Creator,
		Address: wallet.Address,
	}
}

// toSetDefaultWalletDomain - mapping func to domain model
func toSetDefaultWalletDomain(wallet types.DBSetDefaultWallet) *walletstypes.MsgSetDefaultWallet {
	return &walletstypes.MsgSetDefaultWallet{
		Creator: wallet.Creator,
		Address: wallet.Address,
	}
}
