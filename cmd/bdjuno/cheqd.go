package main

import (
	didtypes "github.com/cheqd/cheqd-node/x/did/types"

	resources "github.com/cheqd/cheqd-node/x/resource/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/forbole/juno/v4/modules/messages"
)

// CheqdAddressesParser represents a MessageAddressesParser for the my custom module
// here, we're using a DID as the address
func CheqdAddressesParser(_ codec.Codec, cosmosMsg sdk.Msg) ([]string, error) {
	switch msg := cosmosMsg.(type) {

	case *didtypes.MsgCreateDidDoc:
		return []string{msg.GetPayload().GetId()}, nil

	case *didtypes.MsgUpdateDidDoc:
		return []string{msg.GetPayload().GetId()}, nil

	case *resources.MsgCreateResource:
		return []string{msg.GetPayload().GetId()}, nil

	case *didtypes.MsgDeactivateDidDoc:
		return []string{msg.GetPayload().GetId()}, nil

	default:
		return nil, messages.MessageNotSupported(cosmosMsg)
	}
}
