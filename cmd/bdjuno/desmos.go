package main

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	profilestypes "github.com/desmos-labs/desmos/v4/x/profiles/types"
	junomessages "github.com/forbole/juno/v4/modules/messages"
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
		return []string{msg.Sender, msg.Receiver, msg.NewDTag}, nil

	case *profilestypes.MsgRefuseDTagTransferRequest:
		return []string{msg.Sender, msg.Receiver}, nil

	case *profilestypes.MsgSaveProfile:
		return []string{msg.Creator, msg.DTag, msg.Bio, msg.Nickname, msg.CoverPicture, msg.ProfilePicture}, nil

	case *profilestypes.MsgDeleteProfile:
		return []string{msg.Creator}, nil

	case *profilestypes.MsgLinkApplication:
		return []string{msg.Sender, msg.SourceChannel, msg.SourcePort, msg.CallData}, nil

	case *profilestypes.MsgUnlinkApplication:
		return []string{msg.Signer, msg.Username, msg.Application}, nil

	case *profilestypes.MsgLinkChainAccount:
		return []string{msg.Signer, msg.ChainAddress.TypeUrl, msg.ChainConfig.Name, msg.Proof.String()}, nil

	case *profilestypes.MsgUnlinkChainAccount:
		return []string{msg.Owner, msg.ChainName, msg.Target}, nil

	case *profilestypes.MsgSetDefaultExternalAddress:
		return []string{msg.Signer, msg.ChainName, msg.Target}, nil

	}

	return nil, junomessages.MessageNotSupported(cosmosMsg)
}
