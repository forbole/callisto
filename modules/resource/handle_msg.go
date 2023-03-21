package resource

import (
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/forbole/bdjuno/v4/types"
	juno "github.com/forbole/juno/v4/types"
)

// HandleMsg implements MessageModule
func (m *Module) HandleMsg(index int, msg sdk.Msg, tx *juno.Tx) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch cosmosMsg := msg.(type) {
	case *resourcetypes.MsgCreateResource:
		return m.handleMsgCreateResource(tx.Height, cosmosMsg, tx.FeePayer().String())
	default:
		return nil
	}

}

func (m *Module) handleMsgCreateResource(height int64, msg *resourcetypes.MsgCreateResource, feePayer string) error {
	return m.db.SaveResource(types.NewResource(msg.Payload.Id, msg.Payload.CollectionId,
		msg.Payload.Data, msg.Payload.Name, msg.Payload.Version,
		msg.Payload.ResourceType, msg.Payload.AlsoKnownAs, feePayer, height))
}
