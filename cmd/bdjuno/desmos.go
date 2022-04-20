package main

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	profilestypes "github.com/desmos-labs/desmos/v3/x/profiles/types"
	relationshiptypes "github.com/desmos-labs/desmos/v3/x/relationships/types"

	junomessages "github.com/forbole/juno/v3/modules/messages"
)

// desmosMessageAddressesParser represents a parser able to get the addresses of the involved
// account from a Desmos message
var desmosMessageAddressesParser = junomessages.JoinMessageParsers(
	profilesMessageAddressesParser,
)

// profilesMessageAddressesParser represents a MessageAddressesParser for the x/profiles module
func profilesMessageAddressesParser(_ codec.Codec, cosmosMsg sdk.Msg) ([]string, error) {
	switch msg := cosmosMsg.(type) {

	case *profilestypes.MsgRequestDTagTransfer:
		return []string{msg.Sender, msg.Receiver}, nil

	case *profilestypes.MsgCancelDTagTransferRequest:
		return []string{msg.Sender, msg.Receiver}, nil

	case *profilestypes.MsgAcceptDTagTransferRequest:
		return []string{msg.Sender, msg.Receiver}, nil

	case *profilestypes.MsgRefuseDTagTransferRequest:
		return []string{msg.Sender, msg.Receiver}, nil

	case *relationshiptypes.MsgCreateRelationship:
		return []string{msg.Signer, msg.Counterparty}, nil

	case *relationshiptypes.MsgDeleteRelationship:
		return []string{msg.Signer, msg.Counterparty}, nil

	case *relationshiptypes.MsgBlockUser:
		return []string{msg.Blocker, msg.Blocked}, nil

	case *relationshiptypes.MsgUnblockUser:
		return []string{msg.Blocker, msg.Blocked}, nil

	}

	return nil, junomessages.MessageNotSupported(cosmosMsg)
}