package bank


import (
	"github.com/desmos-labs/juno/parse"
	juno "github.com/desmos-labs/juno/parse/worker"
	"github.com/forbole/bdjuno/x/types"
)

// BankModule represent /x/Bank module
type BankModule struct {}

// Name return the name of the module
func (m BankModule) Name() string {
  return "bank" 
} 

// BlockHandlers return a list of block handler of the module
func (m BankModule) BlockHandlers() []juno.BlockHandler {
  return []juno.BlockHandler{}
} 

// TxHandlers return a list of TxHandlers of the module
func (m BankModule) TxHandlers() []juno.TxHandler{
	return []juno.TxHandler{}
}

// MsgHandlers return a list of MsgHandlers of the module
func (m BankModule)	MsgHandlers() []juno.MsgHandler{
	return []juno.MsgHandler{MsgHandler}
}

// AdditionalOperations return a list of AdditionalOperations of the module
func (m BankModule)AdditionalOperations() []parse.AdditionalOperation{
	return []parse.AdditionalOperation{}
}

// PeriodicOperations return a list of PeriodicOperations of the module
func (m BankModule)PeriodicOperations() []x.PerodicOperation{
	return []x.PerodicOperation{}
}

// GenesisHandlers return a list of GenesisHandlers of the module
func (m BankModule) GenesisHandlers() []juno.GenesisHandler {
	return []juno.GenesisHandler{}
}
