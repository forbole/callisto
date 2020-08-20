package staking


import (
	"github.com/desmos-labs/juno/parse"
	juno "github.com/desmos-labs/juno/parse/worker"
	"github.com/forbole/bdjuno/x/types"
)

// StakingModule represent /x/Staking module
type StakingModule struct {}

// Name return the name of the module
func (m StakingModule) Name() string {
  return "staking" 
} 

// BlockHandlers return a list of block handler of the module
func (m StakingModule) BlockHandlers() []juno.BlockHandler {
  return []juno.BlockHandler{}
} 

// TxHandlers return a list of TxHandlers of the module
func (m StakingModule) TxHandlers() []juno.TxHandler{
	return []juno.TxHandler{}
}

// MsgHandlers return a list of MsgHandlers of the module
func (m StakingModule)	MsgHandlers() []juno.TxHandler{
	return []juno.TxHandler{}
}

// AdditionalOperations return a list of AdditionalOperations of the module
func (m StakingModule)AdditionalOperations() []parse.AdditionalOperation{
	return []parse.AdditionalOperation{}
}

// PeriodicOperations return a list of PeriodicOperations of the module
func (m StakingModule)PeriodicOperations() []x.PerodicOperation{
	return []x.PerodicOperation{PeriodicStakingOperations}
}
