package mint


import (
	"github.com/desmos-labs/juno/parse"
	juno "github.com/desmos-labs/juno/parse/worker"
	"github.com/forbole/bdjuno/x/types"
)

// MintModule represent /x/Mint module
type MintModule struct {}

// Name return the name of the module
func (m MintModule) Name() string {
  return "supply" 
} 

// BlockHandlers return a list of block handler of the module
func (m MintModule) BlockHandlers() []juno.BlockHandler {
  return []juno.BlockHandler{}
} 

// TxHandlers return a list of TxHandlers of the module
func (m MintModule) TxHandlers() []juno.TxHandler{
	return []juno.TxHandler{}
}

// MsgHandlers return a list of MsgHandlers of the module
func (m MintModule)	MsgHandlers() []juno.MsgHandler{
	return []juno.MsgHandler{}
}

// AdditionalOperations return a list of AdditionalOperations of the module
func (m MintModule)AdditionalOperations() []parse.AdditionalOperation{
	return []parse.AdditionalOperation{}
}

// PeriodicOperations return a list of PeriodicOperations of the module
func (m MintModule)PeriodicOperations() []x.PerodicOperation{
	return []x.PerodicOperation{PeriodicMintOperations}
}

// GenesisHandlers return a list of GenesisHandlers of the module
func (m MintModule) GenesisHandlers() []juno.GenesisHandler {
	return []juno.GenesisHandler{}
}
