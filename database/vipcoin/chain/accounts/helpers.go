package accounts

import (
	"encoding/base64"
	"encoding/hex"

	accountstypes "git.ooo.ua/vipcoin/chain/x/accounts/types"
	extratypes "git.ooo.ua/vipcoin/chain/x/types"
	"github.com/cosmos/cosmos-sdk/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/lib/pq"

	"github.com/forbole/bdjuno/v2/database/types"
)

const (
	tableState            = "vipcoin_chain_accounts_set_state"
	tableKinds            = "vipcoin_chain_accounts_set_kinds"
	tableRegisterUser     = "vipcoin_chain_accounts_register_user"
	tableExtra            = "vipcoin_chain_accounts_set_extra"
	tableAffiliateExtra   = "vipcoin_chain_accounts_set_affiliate_extra"
	tableAffiliateAddress = "vipcoin_chain_accounts_set_affiliate_address"
	tableAccountMigrate   = "vipcoin_chain_accounts_account_migrate"
	tableAffiliates       = "vipcoin_chain_accounts_affiliates"
	tableAccounts         = "vipcoin_chain_accounts_accounts"
	tableCreateAccount    = "vipcoin_chain_accounts_create_account"
	tableAddAffiliate     = "vipcoin_chain_accounts_add_affiliate"
)

// toAffiliatesDatabase - mapping func to database model
func toAffiliatesDatabase(affiliate *accountstypes.Affiliate) types.DBAffiliates {
	return types.DBAffiliates{
		Address:         affiliate.Address,
		AffiliationKind: accountstypes.AffiliationKind_value[affiliate.Affiliation.String()],
		Extras:          toExtrasDB(affiliate.Extras),
	}

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
func toAccountDatabase(account *accountstypes.Account, cdc codec.Marshaler) (types.DBAccount, error) {
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

// parseID - helper function for convert
func parseID(data []int64) []interface{} {
	result := make([]interface{}, 0, len(data))
	for _, id := range data {
		result = append(result, id)
	}

	return result
}

// toRegisterUserDatabase - mapping func to database model
func toRegisterUserDatabase(user *accountstypes.MsgRegisterUser) types.DBRegisterUser {
	return types.DBRegisterUser{
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

// toRegisterUsersDatabase - mapping func to database model
func toRegisterUsersDatabase(msg ...*accountstypes.MsgRegisterUser) []types.DBRegisterUser {
	result := make([]types.DBRegisterUser, 0, len(msg))
	for _, user := range msg {
		result = append(result, toRegisterUserDatabase(user))
	}

	return result
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
func toSetKindsDatabase(kinds *accountstypes.MsgSetKinds) types.DBSetKinds {
	return types.DBSetKinds{
		Creator: kinds.Creator,
		Hash:    kinds.Hash,
		Kinds:   toKindsDB(kinds.Kinds),
	}
}

// toKindsArrDatabase - mapping func to database model
func toKindsArrDatabase(msg ...*accountstypes.MsgSetKinds) []types.DBSetKinds {
	result := make([]types.DBSetKinds, 0, len(msg))
	for _, kind := range msg {
		result = append(result, toSetKindsDatabase(kind))
	}

	return result
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
func toSetAffiliateAddressDatabase(msg *accountstypes.MsgSetAffiliateAddress) types.DBSetAffiliateAddress {
	return types.DBSetAffiliateAddress{
		Creator:    msg.Creator,
		Hash:       msg.Hash,
		OldAddress: msg.OldAddress,
		NewAddress: msg.NewAddress,
	}
}

// toSetAffiliatesAddressDatabase - mapping func to database model
func toSetAffiliatesAddressDatabase(msg ...*accountstypes.MsgSetAffiliateAddress) []types.DBSetAffiliateAddress {
	result := make([]types.DBSetAffiliateAddress, 0, len(msg))
	for _, affiliate := range msg {
		result = append(result, toSetAffiliateAddressDatabase(affiliate))
	}

	return result
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
func toAccountMigrateDatabase(msg *accountstypes.MsgAccountMigrate) types.DBAccountMigrate {
	return types.DBAccountMigrate{
		Creator:   msg.Creator,
		Address:   msg.Address,
		Hash:      msg.Hash,
		PublicKey: msg.PublicKey,
	}
}

// toAccountsMigrateDatabase - mapping func to database model
func toAccountsMigrateDatabase(msg ...*accountstypes.MsgAccountMigrate) []types.DBAccountMigrate {
	result := make([]types.DBAccountMigrate, 0, len(msg))
	for _, account := range msg {
		result = append(result, toAccountMigrateDatabase(account))
	}

	return result
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
func toSetExtraDatabase(msg *accountstypes.MsgSetExtra) types.DBSetAccountExtra {
	return types.DBSetAccountExtra{
		Creator: msg.Creator,
		Hash:    msg.Hash,
		Extras:  toExtrasDB(msg.Extras),
	}
}

// toSetExtrasDatabase - mapping func to database model
func toSetExtrasDatabase(msg ...*accountstypes.MsgSetExtra) []types.DBSetAccountExtra {
	result := make([]types.DBSetAccountExtra, 0, len(msg))
	for _, extra := range msg {
		result = append(result, toSetExtraDatabase(extra))
	}

	return result
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
func toSetAffiliateExtraDatabase(msg *accountstypes.MsgSetAffiliateExtra) types.DBSetAffiliateExtra {
	return types.DBSetAffiliateExtra{
		Creator:         msg.Creator,
		AccountHash:     msg.AccountHash,
		AffiliationHash: msg.AffiliationHash,
		Extras:          toExtrasDB(msg.Extras),
	}
}

// toSetAffiliatesExtraDatabase - mapping func to database model
func toSetAffiliatesExtraDatabase(msg ...*accountstypes.MsgSetAffiliateExtra) []types.DBSetAffiliateExtra {
	result := make([]types.DBSetAffiliateExtra, 0, len(msg))
	for _, extra := range msg {
		result = append(result, toSetAffiliateExtraDatabase(extra))
	}

	return result
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
func toSetStateDatabase(msg *accountstypes.MsgSetState) types.DBSetState {
	return types.DBSetState{
		Creator: msg.Creator,
		Hash:    msg.Hash,
		State:   int32(msg.State),
		Reason:  msg.Reason,
	}
}

// toSetStatesDatabase - mapping func to database model
func toSetStatesDatabase(msg ...*accountstypes.MsgSetState) []types.DBSetState {
	result := make([]types.DBSetState, 0, len(msg))
	for _, state := range msg {
		result = append(result, toSetStateDatabase(state))
	}

	return result
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
func toCreateAccountDatabase(msg *accountstypes.MsgCreateAccount) types.DBCreateAccount {
	return types.DBCreateAccount{
		Creator:   msg.Creator,
		Hash:      msg.Hash,
		Address:   msg.Address,
		PublicKey: msg.PublicKey,
		Kinds:     toKindsDB(msg.Kinds),
		State:     int32(msg.State),
		Extras:    toExtrasDB(msg.Extras),
	}
}

// toSetStatesDatabase - mapping func to database model
func toCreateAccountsDatabase(msg ...*accountstypes.MsgCreateAccount) []types.DBCreateAccount {
	result := make([]types.DBCreateAccount, 0, len(msg))
	for _, account := range msg {
		result = append(result, toCreateAccountDatabase(account))
	}

	return result
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
func toAddAffiliateDatabase(msg *accountstypes.MsgAddAffiliate) types.DBAddAffiliate {
	return types.DBAddAffiliate{
		Creator:         msg.Creator,
		AccountHash:     msg.AccountHash,
		AffiliationHash: msg.AffiliationHash,
		Affiliation:     int32(msg.Affiliation),
		Extras:          toExtrasDB(msg.Extras),
	}
}

// toAddAffiliatesDatabase - mapping func to database model
func toAddAffiliatesDatabase(msg ...*accountstypes.MsgAddAffiliate) []types.DBAddAffiliate {
	result := make([]types.DBAddAffiliate, 0, len(msg))
	for _, affiliate := range msg {
		result = append(result, toAddAffiliateDatabase(affiliate))
	}

	return result
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
