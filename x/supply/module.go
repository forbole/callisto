package supply


import (
	"github.com/desmos-labs/juno/parse"
	juno "github.com/desmos-labs/juno/parse/worker"
	"github.com/forbole/bdjuno/x/types"
)

// SupplyModule represent /x/Supply module
type SupplyModule struct {}

// Name return the name of the module
func (m SupplyModule) Name() string {
  return "supply" 
} 

// BlockHandlers return a list of block handler of the module
func (m SupplyModule) BlockHandlers() []juno.BlockHandler {
  return []juno.BlockHandler{}
} 

// TxHandlers return a list of TxHandlers of the module
func (m SupplyModule) TxHandlers() []juno.TxHandler{
	return []juno.TxHandler{}
}

// MsgHandlers return a list of MsgHandlers of the module
func (m SupplyModule)	MsgHandlers() []juno.TxHandler{
	return []juno.TxHandler{}
}

// AdditionalOperations return a list of AdditionalOperations of the module
func (m SupplyModule)AdditionalOperations() []parse.AdditionalOperation{
	return []parse.AdditionalOperation{}
}

// PeriodicOperations return a list of PeriodicOperations of the module
func (m SupplyModule)PeriodicOperations() []x.PerodicOperation{
	return []x.PerodicOperation{PeriodicSupplyOperations}
}
