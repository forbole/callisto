package margin

import (
	"fmt"

	margintypes "github.com/Sifchain/sifnode/x/margin/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/forbole/bdjuno/v3/types"
	juno "github.com/forbole/juno/v3/types"
)

// HandleMsg implements MessageModule
func (m *Module) HandleMsg(index int, msg sdk.Msg, tx *juno.Tx) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch cosmosMsg := msg.(type) {
	case *margintypes.MsgUpdateParams:
		return m.handleMsgUpdateParams(tx.Height, cosmosMsg)
	}

	return nil
}

// ---------------------------------------------------------------------------------------------------------------------

// handleMsgUpdateParams handles properly a MsgUpdateParams instance by
// saving into the database updated margin params
func (m *Module) handleMsgUpdateParams(height int64, msg *margintypes.MsgUpdateParams) error {

	// Save the params
	err := m.db.SaveMarginParams(types.NewMarginParams(msg.Params, height))
	if err != nil {
		return fmt.Errorf("error while storing updated margin params from MsgUpdateParams: %s", err)
	}

	return nil
}
