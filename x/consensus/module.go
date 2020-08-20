package consensus


import (
	"github.com/desmos-labs/juno/parse"
	juno "github.com/desmos-labs/juno/parse/worker"
	"github.com/forbole/bdjuno/x/types"
)

// ConsensusModule represent /x/Consensus module
type ConsensusModule struct {}

// Name return the name of the module
func (m ConsensusModule) Name() string {
  return "consensus" 
} 

// BlockHandlers return a list of block handler of the module
func (m ConsensusModule) BlockHandlers() []juno.BlockHandler {
  return []juno.BlockHandler{}
} 

// TxHandlers return a list of TxHandlers of the module
func (m ConsensusModule) TxHandlers() []juno.TxHandler{
	return []juno.TxHandler{}
}

// MsgHandlers return a list of MsgHandlers of the module
func (m ConsensusModule)	MsgHandlers() []juno.TxHandler{
	return []juno.TxHandler{}
}

// AdditionalOperations return a list of AdditionalOperations of the module
func (m ConsensusModule)AdditionalOperations() []parse.AdditionalOperation{
	return []parse.AdditionalOperation{}
}

// PeriodicOperations return a list of PeriodicOperations of the module
func (m ConsensusModule)PeriodicOperations() []x.PerodicOperation{
	return []x.PerodicOperation{PeriodicConsensusOperations}
}
