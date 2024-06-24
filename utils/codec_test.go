package utils

import (
	"testing"

	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/stretchr/testify/require"
)

func TestUnpackMessage(t *testing.T) {
	msgBz := []byte("{\"@type\":\"/cosmos.gov.v1beta1.MsgVote\",\"proposal_id\":\"797\",\"voter\":\"osmo1qk74379nc34ynvc4gmw2zrwlscn62mx5wps5xj\",\"option\":\"VOTE_OPTION_YES\"}")
	cdc := GetCodec()
	cosmosMsg := UnpackMessage(cdc, msgBz, &govtypesv1beta1.MsgVote{})
	require.Equal(t, &govtypesv1beta1.MsgVote{
		ProposalId: 797,
		Voter:      "osmo1qk74379nc34ynvc4gmw2zrwlscn62mx5wps5xj",
		Option:     govtypesv1beta1.OptionYes,
	}, cosmosMsg)
}
