/*
 * Copyright 2022 Business Process Technologies. All rights reserved.
 */

package accounts

import (
	"encoding/base64"
	"encoding/hex"

	accountstypes "git.ooo.ua/vipcoin/chain/x/accounts/types"
	extratypes "git.ooo.ua/vipcoin/chain/x/types"
	"github.com/cosmos/cosmos-sdk/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/forbole/bdjuno/v2/database/types"
	"github.com/lib/pq"
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
		afaffiliate := accountstypes.Affiliate{
			Address:     data.Address,
			Affiliation: accountstypes.AffiliationKind(data.AffiliationKind),
			Extras:      fromExtrasDB(data.Extras),
		}

		result = append(result, &afaffiliate)
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
		result = append(result, &extra)
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

// toSetKindsDomain - mapping func to database model
func toSetKindsDomain(kinds types.DBSetKinds) *accountstypes.MsgSetKinds {
	return &accountstypes.MsgSetKinds{
		Creator: kinds.Creator,
		Hash:    kinds.Hash,
		Kinds:   toKindsDomain(kinds.Kinds),
	}
}
