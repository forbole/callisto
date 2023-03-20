package did

import (
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
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
	case *didtypes.MsgCreateDidDoc:
		return m.handleMsgCreateDidDoc(tx.Height, cosmosMsg)

	case *didtypes.MsgUpdateDidDoc:
		return m.handleMsgUpdateDidDoc(tx.Height, cosmosMsg)

	case *didtypes.MsgDeactivateDidDoc:
		return m.handleMsgDeactivateDidDoc()

	case *resourcetypes.MsgCreateResource:
		return m.handleMsgCreateResource()

	}

	return nil
}

func (m *Module) handleMsgCreateDidDoc(height int64, msg *didtypes.MsgCreateDidDoc) error {
	// payerAddress := tx.FeePayer()
	return m.db.SaveDidDoc(types.NewDidDoc(msg.Payload.Id, msg.Payload.Context,
		msg.Payload.Controller, msg.Payload.VerificationMethod, msg.Payload.Authentication,
		msg.Payload.AssertionMethod, msg.Payload.CapabilityInvocation,
		msg.Payload.CapabilityDelegation, msg.Payload.KeyAgreement,
		msg.Payload.Service, msg.Payload.AlsoKnownAs, msg.Payload.VersionId, height))
}

func (m *Module) handleMsgUpdateDidDoc(height int64, msg *didtypes.MsgUpdateDidDoc) error {
	return m.db.SaveDidDoc(types.NewDidDoc(msg.Payload.Id, msg.Payload.Context,
		msg.Payload.Controller, msg.Payload.VerificationMethod, msg.Payload.Authentication,
		msg.Payload.AssertionMethod, msg.Payload.CapabilityInvocation,
		msg.Payload.CapabilityDelegation, msg.Payload.KeyAgreement,
		msg.Payload.Service, msg.Payload.AlsoKnownAs, msg.Payload.VersionId, height))

}

func (m *Module) handleMsgDeactivateDidDoc() error {
	return nil
}

func (m *Module) handleMsgCreateResource() error {
	return nil
}
