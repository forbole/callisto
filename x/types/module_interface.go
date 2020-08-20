package x

import (
	"github.com/desmos-labs/juno/parse"
	juno "github.com/desmos-labs/juno/parse/worker"
	"github.com/go-co-op/gocron"
)

// PerodicOperation represent a function that allow to do function perodically
// need to pass in the scheduler in order to have the same timer as other perodic operations 
// this return a additionalOperation such that it will do though juno's additionalOperation channel
type PerodicOperation func(scheduler *gocron.Scheduler) parse.AdditionalOperation

// Module is a interface for all modules inside \x
type Module interface {
	Name() string
	BlockHandlers() []juno.BlockHandler
	TxHandlers() []juno.TxHandler
	MsgHandlers() []juno.TxHandler
	AdditionalOperations() []parse.AdditionalOperation
	PeriodicOperations() []PerodicOperation
}