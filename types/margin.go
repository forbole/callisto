package types

import (
	margintypes "github.com/Sifchain/sifnode/x/margin/types"
	sdk "github.com/tendermint/tendermint/abci/types"
)

// MarginParams represents the x/margin parameters
type MarginParams struct {
	*margintypes.Params
	Height int64
}

// NewMarginParams allows to build a new MarginParams instance
func NewMarginParams(params *margintypes.Params, height int64) *MarginParams {
	return &MarginParams{
		Params: params,
		Height: height,
	}
}

type MarginEvent struct {
	TxHash     string
	Index      int
	MsgType    string
	Value      sdk.Event
	Addressess []string
	Height     int64
}

func NewMarginEvent(txHash string, index int,
	msgType string, value sdk.Event, addresses []string,
	height int64) *MarginEvent {
	return &MarginEvent{
		TxHash:     txHash,
		Index:      index,
		MsgType:    msgType,
		Value:      value,
		Addressess: addresses,
		Height:     height,
	}
}
