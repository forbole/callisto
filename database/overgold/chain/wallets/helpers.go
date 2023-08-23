package wallets

import (
	extratypes "git.ooo.ua/vipcoin/chain/x/types"
	walletstypes "git.ooo.ua/vipcoin/chain/x/wallets/types"
	cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"

	"github.com/forbole/bdjuno/v3/database/types"
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

	result = append(result, balance.Balance...)

	return result
}

// toWalletsDatabase - mapping func to database model
func toWalletsDatabase(wallets *walletstypes.Wallet) types.DBWallets {
	return types.DBWallets{
		Address:        wallets.Address,
		AccountAddress: wallets.AccountAddress,
		Kind:           int32(wallets.Kind),
		State:          int32(wallets.State),
		Balance:        toBalanceDB(wallets.Balance),
		Extras:         toExtrasDB(wallets.Extras),
		DefaultStatus:  wallets.Default,
	}
}

// toWalletDomain - mapping func to domain model
func toWalletDomain(wallet types.DBWallets) *walletstypes.Wallet {
	return &walletstypes.Wallet{
		Address:        wallet.Address,
		AccountAddress: wallet.AccountAddress,
		Kind:           walletstypes.WalletKind(wallet.Kind),
		State:          walletstypes.WalletState(wallet.State),
		Balance:        fromBalanceDB(wallet.Balance),
		Extras:         fromExtrasDB(wallet.Extras),
		Default:        wallet.DefaultStatus,
	}
}

// toCreateWalletDatabase - mapping func to database model
func toCreateWalletDatabase(wallet *walletstypes.MsgCreateWallet, transactionHash string) types.DBCreateWallet {
	return types.DBCreateWallet{
		Hash:           transactionHash,
		Creator:        wallet.Creator,
		Address:        wallet.Address,
		AccountAddress: wallet.AccountAddress,
		Kind:           int32(wallet.Kind),
		State:          int32(wallet.State),
		Extras:         toExtrasDB(wallet.Extras),
		AddressPayFrom: wallet.AddressPayFrom,
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
		AddressPayFrom: wallet.AddressPayFrom,
	}
}

// toSetWalletStateDatabase - mapping func to database model
func toSetWalletStateDatabase(wallet *walletstypes.MsgSetWalletState, transactionHash string) types.DBSetWalletState {
	return types.DBSetWalletState{
		Hash:    transactionHash,
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
func toSetDefaultWalletDatabase(wallet *walletstypes.MsgSetDefaultWallet, transactionHash string) types.DBSetDefaultWallet {
	return types.DBSetDefaultWallet{
		Hash:    transactionHash,
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

// toSetExtraDatabase - mapping func to database model
func toSetExtraDatabase(wallet *walletstypes.MsgSetExtra, transactionHash string) types.DBSetExtra {
	return types.DBSetExtra{
		Hash:    transactionHash,
		Creator: wallet.Creator,
		Address: wallet.Address,
		Extras:  toExtrasDB(wallet.Extras),
	}
}

// toSetExtraDomain - mapping func to domain model
func toSetExtraDomain(wallet types.DBSetExtra) *walletstypes.MsgSetExtra {
	return &walletstypes.MsgSetExtra{
		Creator: wallet.Creator,
		Address: wallet.Address,
		Extras:  fromExtrasDB(wallet.Extras),
	}
}

// toWalletsDatabase - mapping func to database model
func toCreateWalletWithBalanceDatabase(
	wallets *walletstypes.MsgCreateWalletWithBalance,
	transactionHash string,
) types.DBCreateWalletWithBalance {
	return types.DBCreateWalletWithBalance{
		Hash:           transactionHash,
		Creator:        wallets.Creator,
		Address:        wallets.Address,
		AccountAddress: wallets.AccountAddress,
		Kind:           int32(wallets.Kind),
		State:          int32(wallets.State),
		Extras:         toExtrasDB(wallets.Extras),
		DefaultStatus:  wallets.Default,
		Balance:        toBalanceDB(wallets.Balance),
	}
}

// toWalletDomain - mapping func to domain model
func toCreateWalletWithBalanceDomain(wallet types.DBCreateWalletWithBalance) *walletstypes.MsgCreateWalletWithBalance {
	return &walletstypes.MsgCreateWalletWithBalance{
		Creator:        wallet.Creator,
		Address:        wallet.Address,
		AccountAddress: wallet.AccountAddress,
		Kind:           walletstypes.WalletKind(wallet.Kind),
		State:          walletstypes.WalletState(wallet.State),
		Extras:         fromExtrasDB(wallet.Extras),
		Default:        wallet.DefaultStatus,
		Balance:        fromBalanceDB(wallet.Balance),
	}
}

// toSetKindDatabase - mapping func to database model
func toSetKindDatabase(wallet *walletstypes.MsgSetWalletKind, transactionHash string) types.DBSetWalletKind {
	return types.DBSetWalletKind{
		Hash:    transactionHash,
		Creator: wallet.Creator,
		Address: wallet.Address,
		Kind:    int32(wallet.Kind),
	}
}

// toSetKindsDomain - mapping func to domain model
func toSetKindDomain(wallet types.DBSetWalletKind) *walletstypes.MsgSetWalletKind {
	return &walletstypes.MsgSetWalletKind{
		Creator: wallet.Creator,
		Address: wallet.Address,
		Kind:    walletstypes.WalletKind(wallet.Kind),
	}
}

// toSetCreateUserWalletPriceDatabase - mapping func to database model
func toSetCreateUserWalletPriceDatabase(wallet *walletstypes.MsgSetCreateUserWalletPrice, transactionHash string) types.DBSetCreateUserWalletPrice {
	return types.DBSetCreateUserWalletPrice{
		Hash:    transactionHash,
		Creator: wallet.Creator,
		Amount:  wallet.Amount,
	}
}

// toSetCreateUserWalletPriceDomain - mapping func to domain model
func toSetCreateUserWalletPriceDomain(wallet types.DBSetCreateUserWalletPrice) *walletstypes.MsgSetCreateUserWalletPrice {
	return &walletstypes.MsgSetCreateUserWalletPrice{
		Creator: wallet.Creator,
		Amount:  wallet.Amount,
	}
}
