package global

import (
	"fmt"

	globaltypes "github.com/KYVENetwork/chain/x/global/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/forbole/bdjuno/v4/types"
	juno "github.com/forbole/juno/v5/types"
)

// HandleMsg implements modules.MessageModule
func (m *Module) HandleMsg(index int, msg sdk.Msg, tx *juno.Tx) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch msg.(type) {
	case *globaltypes.MsgUpdateParams:
		return m.handleMsgUpdateParams(tx)
	default:
		return nil
	}

}

// handleMsgUpdateParams allows to properly handle a MsgUpdateParams
func (m *Module) handleMsgUpdateParams(tx *juno.Tx) error {
	params, err := m.source.Params(tx.Height)
	if err != nil {
		return fmt.Errorf("error while getting global params: %s", err)
	}

	return m.db.SaveGlobalParams(
		types.NewGlobalParams(params, tx.Height))
}
