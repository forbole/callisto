package main

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	poststypes "github.com/desmos-labs/desmos/v4/x/posts/types"
	profilestypes "github.com/desmos-labs/desmos/v4/x/profiles/types"
	reactionstypes "github.com/desmos-labs/desmos/v4/x/reactions/types"
	relationshiptypes "github.com/desmos-labs/desmos/v4/x/relationships/types"
	reportstypes "github.com/desmos-labs/desmos/v4/x/reports/types"
	subspacestypes "github.com/desmos-labs/desmos/v4/x/subspaces/types"

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

	case *poststypes.MsgCreatePost:
		return []string{msg.Author}, nil
	case *poststypes.MsgEditPost:
		return []string{msg.Editor}, nil
	case *poststypes.MsgDeletePost:
		return []string{msg.Signer}, nil
	case *poststypes.MsgAddPostAttachment:
		return []string{msg.Editor}, nil
	case *poststypes.MsgRemovePostAttachment:
		return []string{msg.Editor}, nil
	case *poststypes.MsgAnswerPoll:
		return []string{msg.Signer}, nil

	case *profilestypes.MsgRequestDTagTransfer:
		return []string{msg.Sender, msg.Receiver}, nil
	case *profilestypes.MsgCancelDTagTransferRequest:
		return []string{msg.Sender, msg.Receiver}, nil
	case *profilestypes.MsgAcceptDTagTransferRequest:
		return []string{msg.Sender, msg.Receiver}, nil
	case *profilestypes.MsgRefuseDTagTransferRequest:
		return []string{msg.Sender, msg.Receiver}, nil

	case *reactionstypes.MsgAddReaction:
		return []string{msg.User}, nil
	case *reactionstypes.MsgRemoveReaction:
		return []string{msg.User}, nil
	case *reactionstypes.MsgAddRegisteredReaction:
		return []string{msg.User}, nil
	case *reactionstypes.MsgEditRegisteredReaction:
		return []string{msg.User}, nil
	case *reactionstypes.MsgRemoveRegisteredReaction:
		return []string{msg.User}, nil
	case *reactionstypes.MsgSetReactionsParams:
		return []string{msg.User}, nil

	case *relationshiptypes.MsgCreateRelationship:
		return []string{msg.Signer, msg.Counterparty}, nil
	case *relationshiptypes.MsgDeleteRelationship:
		return []string{msg.Signer, msg.Counterparty}, nil
	case *relationshiptypes.MsgBlockUser:
		return []string{msg.Blocker, msg.Blocked}, nil
	case *relationshiptypes.MsgUnblockUser:
		return []string{msg.Blocker, msg.Blocked}, nil

	case *reportstypes.MsgCreateReport:
		return []string{msg.Reporter}, nil
	case *reportstypes.MsgDeleteReport:
		return []string{msg.Signer}, nil
	case *reportstypes.MsgSupportStandardReason:
		return []string{msg.Signer}, nil
	case *reportstypes.MsgAddReason:
		return []string{msg.Signer}, nil
	case *reportstypes.MsgRemoveReason:
		return []string{msg.Signer}, nil

	case *subspacestypes.MsgCreateSubspace:
		return []string{msg.Creator}, nil
	case *subspacestypes.MsgEditSubspace:
		return []string{msg.Signer}, nil
	case *subspacestypes.MsgDeleteSubspace:
		return []string{msg.Signer}, nil
	case *subspacestypes.MsgCreateSection:
		return []string{msg.Creator}, nil
	case *subspacestypes.MsgEditSection:
		return []string{msg.Editor}, nil
	case *subspacestypes.MsgMoveSection:
		return []string{msg.Signer}, nil
	case *subspacestypes.MsgDeleteSection:
		return []string{msg.Signer}, nil
	case *subspacestypes.MsgCreateUserGroup:
		return append([]string{msg.Creator}, msg.InitialMembers...), nil
	case *subspacestypes.MsgEditUserGroup:
		return []string{msg.Signer}, nil
	case *subspacestypes.MsgMoveUserGroup:
		return []string{msg.Signer}, nil
	case *subspacestypes.MsgDeleteUserGroup:
		return []string{msg.Signer}, nil
	case *subspacestypes.MsgAddUserToUserGroup:
		return []string{msg.Signer, msg.User}, nil
	case *subspacestypes.MsgRemoveUserFromUserGroup:
		return []string{msg.Signer, msg.User}, nil
	case *subspacestypes.MsgSetUserPermissions:
		return []string{msg.Signer, msg.User}, nil
	}

	return nil, junomessages.MessageNotSupported(cosmosMsg)
}
