package auth

import (
	"github.com/desmos-labs/juno/parse"
	juno "github.com/desmos-labs/juno/parse/worker"
	x "github.com/forbole/bdjuno/x/types"
)

// Module represent /x/Auth module
type Module struct{}

// Name return the name of the module
func (m Module) Name() string {
	return "auth"
}

// BlockHandlers return a list of block handler of the module
func (m Module) BlockHandlers() []juno.BlockHandler {
	return []juno.BlockHandler{}
}

// TxHandlers return a list of TxHandlers of the module
func (m Module) TxHandlers() []juno.TxHandler {
	return []juno.TxHandler{}
}

// MsgHandlers return a list of MsgHandlers of the module
func (m Module) MsgHandlers() []juno.MsgHandler {
	return []juno.MsgHandler{}
}

// AdditionalOperations return a list of AdditionalOperations of the module
func (m Module) AdditionalOperations() []parse.AdditionalOperation {
	return []parse.AdditionalOperation{}
}

// PeriodicOperations return a list of PeriodicOperations of the module
func (m Module) PeriodicOperations() []x.PerodicOperation {
	return []x.PerodicOperation{PeriodicAuthOperations}
}

// GenesisHandlers return a list of GenesisHandlers of the module
func (m Module) GenesisHandlers() []juno.GenesisHandler {
	return []juno.GenesisHandler{GenesisHandler}
}
