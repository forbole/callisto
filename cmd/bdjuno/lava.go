package main

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	junomessages "github.com/forbole/juno/v3/modules/messages"
	confilcttypes "github.com/lavanet/lava/x/conflict/types"

	pairingtypes "github.com/lavanet/lava/x/pairing/types"
)

// lavaMessageAddressesParser represents a parser able to get the addresses of the involved
var LavaMessageAddressesParser = junomessages.JoinMessageParsers(
	lavaMessageAddressesParser,
)

// lavaMessageAddressesParser represents a MessageAddressesParser for the x/profiles module
func lavaMessageAddressesParser(_ codec.Codec, cosmosMsg sdk.Msg) ([]string, error) {
	switch msg := cosmosMsg.(type) {

	case *pairingtypes.MsgRelayPayment:
		return []string{msg.Creator}, nil
	case *pairingtypes.MsgStakeClient:
		return []string{msg.Creator}, nil
	case *pairingtypes.MsgStakeProvider:
		return []string{msg.Creator}, nil
	case *pairingtypes.MsgUnstakeProvider:
		return []string{msg.Creator}, nil
	case *pairingtypes.MsgUnstakeClient:
		return []string{msg.Creator}, nil

	case *confilcttypes.MsgConflictVoteCommit:
		return []string{msg.Creator}, nil
	case *confilcttypes.MsgDetection:
		return []string{msg.Creator}, nil
	case *confilcttypes.MsgConflictVoteReveal:
		return []string{msg.Creator}, nil

	}

	return nil, junomessages.MessageNotSupported(cosmosMsg)
}
