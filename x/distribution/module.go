package distribution


import (
	"github.com/desmos-labs/juno/parse"
	juno "github.com/desmos-labs/juno/parse/worker"
	"github.com/forbole/bdjuno/x/types"
)

// DistributionModule represent /x/Distribution module
type DistributionModule struct {}

// Name return the name of the module
func (m DistributionModule) Name() string {
  return "distribution" 
} 

// BlockHandlers return a list of block handler of the module
func (m DistributionModule) BlockHandlers() []juno.BlockHandler {
  return []juno.BlockHandler{}
} 

// TxHandlers return a list of TxHandlers of the module
func (m DistributionModule) TxHandlers() []juno.TxHandler{
	return []juno.TxHandler{}
}

// MsgHandlers return a list of MsgHandlers of the module
func (m DistributionModule)	MsgHandlers() []juno.MsgHandler{
	return []juno.MsgHandler{}
}

// AdditionalOperations return a list of AdditionalOperations of the module
func (m DistributionModule)AdditionalOperations() []parse.AdditionalOperation{
	return []parse.AdditionalOperation{}
}

// PeriodicOperations return a list of PeriodicOperations of the module
func (m DistributionModule)PeriodicOperations() []x.PerodicOperation{
	return []x.PerodicOperation{PeriodicDistributionOperations}
}


// GenesisHandlers return a list of GenesisHandlers of the module
func (m DistributionModule) GenesisHandlers() []juno.GenesisHandler {
	return []juno.GenesisHandler{}
}
