package pricefeed

import (
	"github.com/desmos-labs/juno/parse"
	juno "github.com/desmos-labs/juno/parse/worker"
	x "github.com/forbole/bdjuno/x/types"
)

// PriceFeedModule represent /x/PriceFeed module
type PriceFeedModule struct{}

// Name return the name of the module
func (m PriceFeedModule) Name() string {
	return "supply"
}

// BlockHandlers return a list of block handler of the module
func (m PriceFeedModule) BlockHandlers() []juno.BlockHandler {
	return []juno.BlockHandler{}
}

// TxHandlers return a list of TxHandlers of the module
func (m PriceFeedModule) TxHandlers() []juno.TxHandler {
	return []juno.TxHandler{}
}
// MsgHandlers return a list of MsgHandlers of the module
func (m PriceFeedModule) MsgHandlers() []juno.MsgHandler {
	return []juno.MsgHandler{}
}

// AdditionalOperations return a list of AdditionalOperations of the module
func (m PriceFeedModule) AdditionalOperations() []parse.AdditionalOperation {
	return []parse.AdditionalOperation{}
}

// PeriodicOperations return a list of PeriodicOperations of the module
func (m PriceFeedModule) PeriodicOperations() []x.PerodicOperation {
	return []x.PerodicOperation{PeriodicPriceFeedOperations}
}

// GenesisHandlers return a list of GenesisHandlers of the module
func (m PriceFeedModule) GenesisHandlers() []juno.GenesisHandler {
	return []juno.GenesisHandler{}
}
