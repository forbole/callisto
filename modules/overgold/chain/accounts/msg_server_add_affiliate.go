package accounts

import (
	"git.ooo.ua/vipcoin/chain/x/accounts/types"
	"git.ooo.ua/vipcoin/lib/filter"
	juno "github.com/forbole/juno/v3/types"

	dbtypes "github.com/forbole/bdjuno/v3/database/types"
)

// handleMsgAccountMigrate allows to properly handle a handleMsgAccountMigrate
func (m *Module) handleMsgAddAffiliate(tx *juno.Tx, index int, msg *types.MsgAddAffiliate) error {
	accountArr, err := m.accountRepo.GetAccounts(filter.NewFilter().SetArgument(dbtypes.FieldHash, msg.AccountHash))
	switch {
	case err != nil:
		return err
	case len(accountArr) != 1:
		return types.ErrInvalidHashField
	}

	affiliateArr, err := m.accountRepo.GetAccounts(filter.NewFilter().SetArgument(dbtypes.FieldHash, msg.AffiliationHash))
	switch {
	case err != nil:
		return err
	case len(affiliateArr) != 1:
		return types.ErrInvalidHashField
	}

	account := accountArr[0]
	affiliate := affiliateArr[0]

	var accIsSet bool // The account affiliates record is already set
	for _, a := range account.Affiliates {
		if a.Address == affiliate.Address {
			accIsSet = true
		}
	}

	var affIsSet bool // The affiliate account affiliates record is already set
	for _, a := range affiliate.Affiliates {
		if a.Address == account.Address {
			affIsSet = true
		}
	}

	if !accIsSet {
		newAffiliate := &types.Affiliate{
			Address:     affiliate.Address,
			Affiliation: msg.Affiliation,
			Extras:      msg.Extras,
		}
		account.Affiliates = append(account.Affiliates, newAffiliate)

		if err := m.accountRepo.UpdateAccounts(account); err != nil {
			return err
		}
	}

	if !affIsSet {
		// 1) set affHash to acc as Referrer, so affiliate holds account as referral
		affiliateRelationType := types.AFFILIATION_KIND_REFERRAL

		// 2) set affHash to acc as Referral, so affiliate holds account as referrer
		if msg.Affiliation == types.AFFILIATION_KIND_REFERRAL {
			affiliateRelationType = types.AFFILIATION_KIND_REFERRER
		}
		newAffiliate := &types.Affiliate{
			Address:     account.Address,
			Affiliation: affiliateRelationType,
			Extras:      msg.Extras,
		}

		affiliate.Affiliates = append(affiliate.Affiliates, newAffiliate)
		if err := m.accountRepo.UpdateAccounts(affiliate); err != nil {
			return err
		}
	}

	return m.accountRepo.SaveAddAffiliate(msg, tx.TxHash)
}
