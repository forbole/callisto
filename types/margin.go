package types

import (
	margintypes "github.com/Sifchain/sifnode/x/margin/types"
	abci "github.com/tendermint/tendermint/abci/types"
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
	TxHash  string
	Index   int
	MsgType string
	Value   abci.Event
	Height  int64
}

func NewMarginEvent(txHash string, index int,
	msgType string, value abci.Event,
	height int64) *MarginEvent {
	return &MarginEvent{
		TxHash:  txHash,
		Index:   index,
		MsgType: msgType,
		Value:   value,
		Height:  height,
	}
}
