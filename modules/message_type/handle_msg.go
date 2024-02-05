package message_type

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
	utils "github.com/forbole/callisto/v4/modules/utils"
	msgtypes "github.com/forbole/callisto/v4/types"

	"github.com/forbole/juno/v5/types"
)

// HandleMsg represents a message handler that stores the given message inside the proper database table
func (m *Module) HandleMsg(
	index int, msg sdk.Msg, tx *types.Tx) error {
	// Save message type
	err := m.db.SaveMessageType(msgtypes.NewMessageType(
		proto.MessageName(msg),
		utils.GetModuleNameFromTypeURL(proto.MessageName(msg)),
		utils.GetMsgFromTypeURL(proto.MessageName(msg)),
		tx.Height))

	if err != nil {
		return err
	}

	return nil
}
