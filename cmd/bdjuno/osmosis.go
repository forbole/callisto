package main

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	junomessages "github.com/desmos-labs/juno/modules/messages"
	gammtypes "github.com/osmosis-labs/osmosis/x/gamm/types"
)

// osmosisMessageAddressesParser represents a parser able to get the addresses of the involved
// accounts from an Osmosis message
func osmosisMessageAddressesParser(_ codec.Marshaler, cosmosMsg sdk.Msg) ([]string, error) {
	if msg, ok := (cosmosMsg).(*gammtypes.MsgCreatePool); ok {
		return []string{msg.Sender, msg.FuturePoolGovernor}, nil
	}

	return nil, junomessages.MessageNotSupported(cosmosMsg)
}
