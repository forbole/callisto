package accounts

import (
	"git.ooo.ua/vipcoin/chain/x/accounts/types"
	"git.ooo.ua/vipcoin/lib/filter"
	juno "github.com/forbole/juno/v2/types"

	dbtypes "github.com/forbole/bdjuno/v2/database/types"
)

// handleMsgAccountMigrate allows to properly handle a handleMsgAccountMigrate
func (m *Module) handleMsgAccountMigrate(tx *juno.Tx, index int, msg *types.MsgAccountMigrate) error {
	if err := m.accountRepo.SaveAccountMigrate(msg); err != nil {
		return err
	}

	account, err := m.accountRepo.GetAccounts(filter.NewFilter().SetArgument(dbtypes.FieldHash, msg.Hash))
	if err != nil {
		return err
	}

	if len(account) != 1 {
		return types.ErrInvalidHashField
	}

	publicKey, err := types.PubKeyFromString(msg.PublicKey)
	if err != nil {
		return types.ErrInvalidPublicKeyField
	}

	account[0].PublicKey, err = types.PubKeyToAny(publicKey)
	if err != nil {
		return types.ErrInvalidPublicKeyField
	}

	account[0].Address = msg.Address

	return m.accountRepo.UpdateAccounts(account...)
}
