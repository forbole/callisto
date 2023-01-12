package accounts

import "git.ooo.ua/vipcoin/chain/x/accounts/types"

// updateAffiliateAddress - helper func for update affiliate address
func updateAffiliateAddress(affiliates []*types.Affiliate, msg *types.MsgSetAffiliateAddress) error {
	for index, affiliate := range affiliates {
		if affiliate.Address == msg.OldAddress {
			affiliates[index].Address = msg.NewAddress
			return nil
		}
	}

	return types.ErrInvalidOldAddressField
}
