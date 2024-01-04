package wormhole

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	juno "github.com/forbole/juno/v5/types"
	wormholetypes "github.com/wormhole-foundation/wormchain/x/wormhole/types"
)

// HandleMsg implements modules.MessageModule
func (m *Module) HandleMsg(_ int, msg sdk.Msg, tx *juno.Tx) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch cosmosMsg := msg.(type) {
	case *wormholetypes.MsgRegisterAccountAsGuardian:
		return m.HandleMsgRegisterAccountAsGuardian(tx, cosmosMsg)
	case *wormholetypes.MsgCreateAllowlistEntryRequest:
		return m.HandleMsgCreateAllowlistEntryRequest(tx, cosmosMsg)
	case *wormholetypes.MsgDeleteAllowlistEntryRequest:
		return m.HandleMsgDeleteAllowlistEntryRequest(tx, cosmosMsg)
	}

	return nil
}

// HandleMsgRegisterAccountAsGuardian allows to properly handle a MsgRegisterAccountAsGuardian
func (m *Module) HandleMsgRegisterAccountAsGuardian(tx *juno.Tx, msg *wormholetypes.MsgRegisterAccountAsGuardian) error {
	// Get the latest guardian validators list
	guardianValidatorsList, err := m.source.GetGuardianValidatorAll(tx.Height)
	if err != nil {
		return err
	}

	// Save refreshed guardian validators list in db
	return m.db.SaveGuardianValidatorList(guardianValidatorsList, tx.Height)
}

// HandleMsgCreateAllowlistEntryRequest allows to properly handle a MsgCreateAllowlistEntryRequest
func (m *Module) HandleMsgCreateAllowlistEntryRequest(tx *juno.Tx, msg *wormholetypes.MsgCreateAllowlistEntryRequest) error {
	// Get the latest validators allow list
	validatorsAllowList, err := m.source.GetAllowlistAll(tx.Height)
	if err != nil {
		return err
	}

	return m.db.SaveValidatorAllowList(validatorsAllowList, tx.Height)
}

// HandleMsgDeleteAllowlistEntryRequest allows to properly handle a MsgDeleteAllowlistEntryRequest
func (m *Module) HandleMsgDeleteAllowlistEntryRequest(tx *juno.Tx, msg *wormholetypes.MsgDeleteAllowlistEntryRequest) error {
	// Remove validator allow liist entry from database
	return m.db.DeleteValidatorAllowListEntry(msg.Address)
}
