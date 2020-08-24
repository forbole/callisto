package x

import (
	"time"

	"github.com/desmos-labs/juno/parse"
	juno "github.com/desmos-labs/juno/parse/worker"
	"github.com/forbole/bdjuno/x/utils"
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
	MsgHandlers() []juno.MsgHandler
	AdditionalOperations() []parse.AdditionalOperation
	PeriodicOperations() []PerodicOperation
	GenesisHandlers() []juno.GenesisHandler
}

// RegisterModules register modules that enabled and pass that to juno handlers
func RegisterModules(modules []Module) {
	scheduler := gocron.NewScheduler(time.UTC)
	var modulesName []string
	for _, module := range modules {
		for _, handler := range module.GenesisHandlers() {
			juno.RegisterGenesisHandler(handler)
		}
		for _, handler := range module.BlockHandlers() {
			juno.RegisterBlockHandler(handler)
		}
		for _, handler := range module.TxHandlers() {
			juno.RegisterTxHandler(handler)
		}
		for _, handler := range module.MsgHandlers() {
			juno.RegisterMsgHandler(handler)
		}
		for _, handler := range module.AdditionalOperations() {
			parse.RegisterAdditionalOperation(handler)
		}
		for _, handler := range module.PeriodicOperations() {
			parse.RegisterAdditionalOperation(handler(scheduler))
		}
		modulesName = append(modulesName, module.Name())
	}
	parse.RegisterAdditionalOperation(utils.WatchModules(modulesName))
	scheduler.StartAsync()
}
