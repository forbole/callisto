package accounts

import (
	"encoding/base64"
	"encoding/hex"

	accountstypes "git.ooo.ua/vipcoin/chain/x/accounts/types"
	extratypes "git.ooo.ua/vipcoin/chain/x/types"
	"github.com/cosmos/cosmos-sdk/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/lib/pq"

	"github.com/forbole/bdjuno/v3/database/types"
)

const (
	tableState            = "overgold_chain_accounts_set_state"
	tableKinds            = "overgold_chain_accounts_set_kinds"
	tableRegisterUser     = "overgold_chain_accounts_register_user"
	tableExtra            = "overgold_chain_accounts_set_extra"
	tableAffiliateExtra   = "overgold_chain_accounts_set_affiliate_extra"
	tableAffiliateAddress = "overgold_chain_accounts_set_affiliate_address"
	tableAccountMigrate   = "overgold_chain_accounts_account_migrate"
	tableAffiliates       = "overgold_chain_accounts_affiliates"
	tableAccounts         = "overgold_chain_accounts_accounts"
	tableCreateAccount    = "overgold_chain_accounts_create_account"
	tableAddAffiliate     = "overgold_chain_accounts_add_affiliate"
)

// toAffiliateDatabase - mapping func to database model
func toAffiliateDatabase(affiliate *accountstypes.Affiliate, accHash string) types.DBAffiliates {
	return types.DBAffiliates{
		Address:         affiliate.Address,
		AccountHash:     accHash,
		AffiliationKind: accountstypes.AffiliationKind_value[affiliate.Affiliation.String()],
		Extras:          toExtrasDB(affiliate.Extras),
	}
}

// toAffiliatesDatabase - mapping func to database model
func toAffiliatesDatabase(affiliates []*accountstypes.Affiliate, accHash string) []types.DBAffiliates {
	result := make([]types.DBAffiliates, 0, len(affiliates))
	for _, affiliate := range affiliates {
		result = append(result, toAffiliateDatabase(affiliate, accHash))
	}

	return result
}

// toAffiliatesDomain - mapping func to domain model
func toAffiliatesDomain(affiliates []types.DBAffiliates) []*accountstypes.Affiliate {
	result := make([]*accountstypes.Affiliate, 0, len(affiliates))
	for _, data := range affiliates {
		affiliates := accountstypes.Affiliate{
			Address:     data.Address,
			Affiliation: accountstypes.AffiliationKind(data.AffiliationKind),
			Extras:      fromExtrasDB(data.Extras),
		}

		result = append(result, &affiliates)
	}

	return result
}

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
		extraCopy := extra
		result = append(result, &extraCopy)
	}

	return result
}

// toAccountDatabase - mapping func to database model
func toAccountDatabase(account *accountstypes.Account, cdc codec.Codec) (types.DBAccount, error) {
	var pubKey cryptotypes.PubKey
	if err := cdc.UnpackAny(account.PublicKey, &pubKey); err != nil {
		return types.DBAccount{}, err
	}

	publicKey := hex.EncodeToString([]byte(base64.StdEncoding.EncodeToString(pubKey.Bytes())))

	return types.DBAccount{
		Address:   account.Address,
		Hash:      account.Hash,
		PublicKey: publicKey,
		Kinds:     toKindsDB(account.Kinds),
		State:     int32(account.State),
		Extras:    toExtrasDB(account.Extras),
		Wallets:   account.Wallets,
	}, nil
}

// toAccountDatabase - mapping func to domain model
func toAccountDomain(account types.DBAccount) (*accountstypes.Account, error) {
	pubKey, err := accountstypes.PubKeyFromString(account.PublicKey)
	if err != nil {
		return &accountstypes.Account{}, err
	}

	pubKeyAny, err := accountstypes.PubKeyToAny(pubKey)
	if err != nil {
		return &accountstypes.Account{}, err
	}

	return &accountstypes.Account{
		Address:   account.Address,
		Hash:      account.Hash,
		PublicKey: pubKeyAny,
		Kinds:     toKindsDomain(account.Kinds),
		State:     accountstypes.AccountState(account.State),
		Extras:    fromExtrasDB(account.Extras),
		Wallets:   account.Wallets,
	}, nil
}

// toKindsDB - mapping func to database model
func toKindsDB(kinds []accountstypes.AccountKind) pq.Int32Array {
	result := make(pq.Int32Array, 0, len(kinds))
	for _, kind := range kinds {
		result = append(result, accountstypes.AccountKind_value[kind.String()])
	}

	return result
}

// toKindsDomain - mapping func to domain model
func toKindsDomain(kinds pq.Int32Array) []accountstypes.AccountKind {
	result := make([]accountstypes.AccountKind, 0, len(kinds))
	for _, kind := range kinds {
		result = append(result, accountstypes.AccountKind(kind))
	}

	return result
}

// toRegisterUserDatabase - mapping func to database model
func toRegisterUserDatabase(user *accountstypes.MsgRegisterUser, transactionHash string) types.DBRegisterUser {
	return types.DBRegisterUser{
		TransactionHash:       transactionHash,
		Creator:               user.Creator,
		Address:               user.Address,
		Hash:                  user.Hash,
		PublicKey:             user.PublicKey,
		HolderWallet:          user.HolderWallet,
		RefRewardWallet:       user.RefRewardWallet,
		HolderWalletExtras:    toExtrasDB(user.HolderWalletExtras),
		RefRewardWalletExtras: toExtrasDB(user.RefRewardWalletExtras),
		ReferrerHash:          user.ReferrerHash,
	}

}

// toRegisterUserDomain - mapping func to database model
func toRegisterUserDomain(user types.DBRegisterUser) *accountstypes.MsgRegisterUser {
	return &accountstypes.MsgRegisterUser{
		Creator:               user.Creator,
		Address:               user.Address,
		Hash:                  user.Hash,
		PublicKey:             user.PublicKey,
		HolderWallet:          user.HolderWallet,
		RefRewardWallet:       user.RefRewardWallet,
		HolderWalletExtras:    fromExtrasDB(user.HolderWalletExtras),
		RefRewardWalletExtras: fromExtrasDB(user.RefRewardWalletExtras),
		ReferrerHash:          user.ReferrerHash,
	}

}

// toSetKindsDatabase - mapping func to database model
func toSetKindsDatabase(kinds *accountstypes.MsgSetKinds, transactionHash string) types.DBSetKinds {
	return types.DBSetKinds{
		TransactionHash: transactionHash,
		Creator:         kinds.Creator,
		Hash:            kinds.Hash,
		Kinds:           toKindsDB(kinds.Kinds),
	}
}

// toSetKindsDomain - mapping func to database model
func toSetKindsDomain(kinds types.DBSetKinds) *accountstypes.MsgSetKinds {
	return &accountstypes.MsgSetKinds{
		Creator: kinds.Creator,
		Hash:    kinds.Hash,
		Kinds:   toKindsDomain(kinds.Kinds),
	}
}

// toSetAffiliateAddressDatabase - mapping func to database model
func toSetAffiliateAddressDatabase(msg *accountstypes.MsgSetAffiliateAddress, transactionHash string) types.DBSetAffiliateAddress {
	return types.DBSetAffiliateAddress{
		TransactionHash: transactionHash,
		Creator:         msg.Creator,
		Hash:            msg.Hash,
		OldAddress:      msg.OldAddress,
		NewAddress:      msg.NewAddress,
	}
}

// toSetAffiliateAddressDomain - mapping func to database model
func toSetAffiliateAddressDomain(msg types.DBSetAffiliateAddress) *accountstypes.MsgSetAffiliateAddress {
	return &accountstypes.MsgSetAffiliateAddress{
		Creator:    msg.Creator,
		Hash:       msg.Hash,
		OldAddress: msg.OldAddress,
		NewAddress: msg.NewAddress,
	}
}

// toAccountMigrateDatabase - mapping func to database model
func toAccountMigrateDatabase(msg *accountstypes.MsgAccountMigrate, transactionHash string) types.DBAccountMigrate {
	return types.DBAccountMigrate{
		TransactionHash: transactionHash,
		Creator:         msg.Creator,
		Address:         msg.Address,
		Hash:            msg.Hash,
		PublicKey:       msg.PublicKey,
	}
}

// toAccountMigrateDomain - mapping func to database model
func toAccountMigrateDomain(msg types.DBAccountMigrate) *accountstypes.MsgAccountMigrate {
	return &accountstypes.MsgAccountMigrate{
		Creator:   msg.Creator,
		Address:   msg.Address,
		Hash:      msg.Hash,
		PublicKey: msg.PublicKey,
	}
}

// toSetExtraDatabase - mapping func to database model
func toSetExtraDatabase(msg *accountstypes.MsgSetExtra, transactionHash string) types.DBSetAccountExtra {
	return types.DBSetAccountExtra{
		TransactionHash: transactionHash,
		Creator:         msg.Creator,
		Hash:            msg.Hash,
		Extras:          toExtrasDB(msg.Extras),
	}
}

// toSetExtraDomain - mapping func to database model
func toSetExtraDomain(msg types.DBSetAccountExtra) *accountstypes.MsgSetExtra {
	return &accountstypes.MsgSetExtra{
		Creator: msg.Creator,
		Hash:    msg.Hash,
		Extras:  fromExtrasDB(msg.Extras),
	}
}

// toSetAffiliateExtraDatabase - mapping func to database model
func toSetAffiliateExtraDatabase(msg *accountstypes.MsgSetAffiliateExtra, transactionHash string) types.DBSetAffiliateExtra {
	return types.DBSetAffiliateExtra{
		TransactionHash: transactionHash,
		Creator:         msg.Creator,
		AccountHash:     msg.AccountHash,
		AffiliationHash: msg.AffiliationHash,
		Extras:          toExtrasDB(msg.Extras),
	}
}

// toSetAffiliateExtraDomain - mapping func to database model
func toSetAffiliateExtraDomain(msg types.DBSetAffiliateExtra) *accountstypes.MsgSetAffiliateExtra {
	return &accountstypes.MsgSetAffiliateExtra{
		Creator:         msg.Creator,
		AccountHash:     msg.AccountHash,
		AffiliationHash: msg.AffiliationHash,
		Extras:          fromExtrasDB(msg.Extras),
	}
}

// toSetStateDatabase - mapping func to database model
func toSetStateDatabase(msg *accountstypes.MsgSetState, transactionHash string) types.DBSetState {
	return types.DBSetState{
		TransactionHash: transactionHash,
		Creator:         msg.Creator,
		Hash:            msg.Hash,
		State:           int32(msg.State),
		Reason:          msg.Reason,
	}
}

// toSetStateDomain - mapping func to database model
func toSetStateDomain(msg types.DBSetState) *accountstypes.MsgSetState {
	return &accountstypes.MsgSetState{
		Creator: msg.Creator,
		Hash:    msg.Hash,
		State:   accountstypes.AccountState(msg.State),
		Reason:  msg.Reason,
	}
}

// toCreateAccountDatabase - mapping func to database model
func toCreateAccountDatabase(msg *accountstypes.MsgCreateAccount, transactionHash string) types.DBCreateAccount {
	return types.DBCreateAccount{
		TransactionHash: transactionHash,
		Creator:         msg.Creator,
		Hash:            msg.Hash,
		Address:         msg.Address,
		PublicKey:       msg.PublicKey,
		Kinds:           toKindsDB(msg.Kinds),
		State:           int32(msg.State),
		Extras:          toExtrasDB(msg.Extras),
	}
}

// toSetStateDomain - mapping func to database model
func toCreateAccountDomain(msg types.DBCreateAccount) *accountstypes.MsgCreateAccount {
	return &accountstypes.MsgCreateAccount{
		Creator:   msg.Creator,
		Hash:      msg.Hash,
		Address:   msg.Address,
		PublicKey: msg.PublicKey,
		Kinds:     toKindsDomain(msg.Kinds),
		State:     accountstypes.AccountState(msg.State),
		Extras:    fromExtrasDB(msg.Extras),
	}
}

// toAddAffiliateDatabase - mapping func to database model
func toAddAffiliateDatabase(msg *accountstypes.MsgAddAffiliate, transactionHash string) types.DBAddAffiliate {
	return types.DBAddAffiliate{
		TransactionHash: transactionHash,
		Creator:         msg.Creator,
		AccountHash:     msg.AccountHash,
		AffiliationHash: msg.AffiliationHash,
		Affiliation:     int32(msg.Affiliation),
		Extras:          toExtrasDB(msg.Extras),
	}
}

// toAddAffiliateDomain - mapping func to database model
func toAddAffiliateDomain(msg types.DBAddAffiliate) *accountstypes.MsgAddAffiliate {
	return &accountstypes.MsgAddAffiliate{
		Creator:         msg.Creator,
		AccountHash:     msg.AccountHash,
		AffiliationHash: msg.AffiliationHash,
		Affiliation:     accountstypes.AffiliationKind(msg.Affiliation),
		Extras:          fromExtrasDB(msg.Extras),
	}
}
