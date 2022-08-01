package accounts

import (
	"git.ooo.ua/vipcoin/chain/x/accounts/types"
	wallets "git.ooo.ua/vipcoin/chain/x/wallets/types"
	"git.ooo.ua/vipcoin/lib/filter"
	dbtypes "github.com/forbole/bdjuno/v2/database/types"
	juno "github.com/forbole/juno/v2/types"
)

// handleMsgRegisterUser allows to properly handle a handleMsgRegisterUser
func (m *Module) handleMsgRegisterUser(tx *juno.Tx, index int, msg *types.MsgRegisterUser) error {
	publicKey, err := types.PubKeyFromString(msg.PublicKey)
	if err != nil {
		return types.ErrInvalidPublicKeyField
	}

	publicKeyAny, err := types.PubKeyToAny(publicKey)
	if err != nil {
		return types.ErrInvalidPublicKeyField
	}

	if msg.ReferrerHash != "" {
		affiliateAcc, err := m.accountRepo.GetAccounts(filter.NewFilter().SetArgument(dbtypes.FieldHash, msg.ReferrerHash))
		switch {
		case err != nil:
			return err
		case len(affiliateAcc) != 1:
			msg.ReferrerHash = ""
		}
	}

	newAcc := types.Account{
		Hash:      msg.Hash,
		Address:   msg.Address,
		Kinds:     []types.AccountKind{types.ACCOUNT_KIND_GENERAL},
		State:     types.ACCOUNT_STATE_ACTIVE,
		PublicKey: publicKeyAny,
		Wallets:   []string{msg.HolderWallet, msg.RefRewardWallet},
	}

	if err := m.accountRepo.SaveAccounts(&newAcc); err != nil {
		return err
	}

	if msg.ReferrerHash != "" {
		affiliateAcc, err := m.accountRepo.GetAccounts(filter.NewFilter().SetArgument(dbtypes.FieldHash, msg.ReferrerHash))
		switch {
		case err != nil:
			return err
		case len(affiliateAcc) != 1:
			return types.ErrInvalidHashField
		}

		affiliate := affiliateAcc[0]
		// check if referrer address is valid
		if affiliate.Address != msg.Address {
			// set Referrer to created account
			a := &types.Affiliate{
				Address:     affiliate.Address,
				Affiliation: types.AFFILIATION_KIND_REFERRER,
			}
			newAcc.Affiliates = append(newAcc.Affiliates, a)
			if err := m.accountRepo.UpdateAccounts(&newAcc); err != nil {
				return err
			}

			// add referral to the Referrer account
			gotThisReferral := false
			for _, aff := range affiliate.Affiliates {
				if aff.Address == msg.Address {
					gotThisReferral = true
				}
			}

			if !gotThisReferral {
				a := &types.Affiliate{
					Address:     msg.Address,
					Affiliation: types.AFFILIATION_KIND_REFERRAL,
				}
				affiliate.Affiliates = append(affiliate.Affiliates, a)
				if err := m.accountRepo.UpdateAccounts(affiliate); err != nil {
					return err
				}
			}
		}
	}

	// create wallets
	holder := wallets.Wallet{
		Address:        msg.HolderWallet,
		AccountAddress: msg.Address,
		Kind:           wallets.WALLET_KIND_HOLDER,
		State:          wallets.WALLET_STATE_ACTIVE,
		Extras:         msg.HolderWalletExtras,
		Default:        true,
	}

	if err := m.walletsRepo.SaveWallets(&holder); err != nil {
		return err
	}

	refReward := wallets.Wallet{
		Address:        msg.RefRewardWallet,
		AccountAddress: msg.Address,
		Kind:           wallets.WALLET_KIND_REFERRER_REWARD,
		State:          wallets.WALLET_STATE_ACTIVE,
		Extras:         msg.RefRewardWalletExtras,
		Default:        false,
	}

	if err := m.walletsRepo.SaveWallets(&refReward); err != nil {
		return err
	}

	return m.accountRepo.SaveRegisterUser(msg, tx.TxHash)
}
