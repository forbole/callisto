package main

import (
	didMsgs "github.com/cheqd/cheqd-node/x/did/types"

	resourceMsgs "github.com/cheqd/cheqd-node/x/resource/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/forbole/juno/v4/modules/messages"
)

// CheqdAddressesParser represents a MessageAddressesParser for the my custom module
// here, we're using a DID as the address
func CheqdAddressesParser(_ codec.Codec, cosmosMsg sdk.Msg) ([]string, error) {
	switch msg := cosmosMsg.(type) {

	case *didMsgs.MsgCreateDidDoc:
		return []string{msg.GetPayload().GetId()}, nil

	case *didMsgs.MsgUpdateDidDoc:
		return []string{msg.GetPayload().GetId()}, nil

	case *resourceMsgs.MsgCreateResource:
		return []string{msg.GetPayload().GetId()}, nil

	default:
		return nil, messages.MessageNotSupported(cosmosMsg)
	}
}
