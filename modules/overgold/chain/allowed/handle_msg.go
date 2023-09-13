package allowed

import (
	"git.ooo.ua/vipcoin/ovg-chain/x/allowed/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	juno "github.com/forbole/juno/v5/types"
)

// HandleMsg implements MessageModule
func (m *Module) HandleMsg(index int, msg sdk.Msg, tx *juno.Tx) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	// switch allowedMsg := msg.(type) {
	switch msg.(type) {
	case *types.MsgCreateAddresses:
		// TODO: create addresses
		return nil
	case *types.MsgUpdateAddresses:
		// TODO: update addresses
		return nil
	case *types.MsgDeleteAddresses:
		// TODO: delete addresses
		return nil
	default:
		return nil
	}
}
