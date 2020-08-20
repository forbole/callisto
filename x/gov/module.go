package gov


import (
	"github.com/desmos-labs/juno/parse"
	juno "github.com/desmos-labs/juno/parse/worker"
	"github.com/forbole/bdjuno/x/types"
)

// GovModule represent /x/Gov module
type GovModule struct {}

// Name return the name of the module
func (m GovModule) Name() string {
  return "gov" 
} 

// BlockHandlers return a list of block handler of the module
func (m GovModule) BlockHandlers() []juno.BlockHandler {
  return []juno.BlockHandler{BlockHandler}
} 

// TxHandlers return a list of TxHandlers of the module
func (m GovModule) TxHandlers() []juno.TxHandler{
	return []juno.TxHandler{}
}

// MsgHandlers return a list of MsgHandlers of the module
func (m GovModule)	MsgHandlers() []juno.MsgHandler{
	return []juno.MsgHandler{MsgHandler}
}

// AdditionalOperations return a list of AdditionalOperations of the module
func (m GovModule)AdditionalOperations() []parse.AdditionalOperation{
	return []parse.AdditionalOperation{}
}

// PeriodicOperations return a list of PeriodicOperations of the module
func (m GovModule)PeriodicOperations() []x.PerodicOperation{
	return []x.PerodicOperation{}
}

// GenesisHandlers return a list of GenesisHandlers of the module
func (m GovModule) GenesisHandlers() []juno.GenesisHandler {
	return []juno.GenesisHandler{GenesisHandler}
}

