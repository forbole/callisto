package utils

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	profilestypes "github.com/desmos-labs/desmos/x/profiles/types"
	junomessages "github.com/desmos-labs/juno/modules/messages"
)

// desmosMessageAddressesParser represents a parser able to get the addresses of the involved
// account from a Desmos message
var desmosMessageAddressesParser = junomessages.JoinMessageParsers(
	profilesMessageAddressesParser,
)

// profilesMessageAddressesParser represents a MessageAddressesParser for the x/profiles module
func profilesMessageAddressesParser(_ codec.Marshaler, cosmosMsg sdk.Msg) ([]string, error) {
	switch msg := cosmosMsg.(type) {

	case *profilestypes.MsgSaveProfile:
		return []string{msg.Creator}, nil

	case *profilestypes.MsgDeleteProfile:
		return []string{msg.Creator}, nil

	case *profilestypes.MsgRequestDTagTransfer:
		return []string{msg.Sender, msg.Receiver}, nil

	case *profilestypes.MsgCancelDTagTransfer:
		return []string{msg.Sender, msg.Receiver}, nil

	case *profilestypes.MsgAcceptDTagTransfer:
		return []string{msg.Sender, msg.Receiver}, nil

	case *profilestypes.MsgRefuseDTagTransfer:
		return []string{msg.Sender, msg.Receiver}, nil

	case *profilestypes.MsgCreateRelationship:
		return []string{msg.Sender, msg.Receiver}, nil

	case *profilestypes.MsgDeleteRelationship:
		return []string{msg.User, msg.Counterparty}, nil

	case *profilestypes.MsgBlockUser:
		return []string{msg.Blocker, msg.Blocked}, nil

	case *profilestypes.MsgUnblockUser:
		return []string{msg.Blocker, msg.Blocked}, nil

	}

	return nil, junomessages.MessageNotSupported(cosmosMsg)
}
