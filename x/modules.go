package modules

import (
	"github.com/desmos-labs/juno/parse/worker"
	djuno "github.com/desmos-labs/juno/types"
	"github.com/rs/zerolog/log"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

var (
	modules []ModuleInfo
)

// Fetcher represents a generic function that accepts all the data that a block handler has, and returns an
// interface which identifies the fetched data for a module.
// If some error raises during the fetching, it is returned instead.
type Fetcher = func(
	block *coretypes.ResultBlock, txs []djuno.Tx, vals *tmctypes.ResultValidators, w worker.Worker,
) (interface{}, error)

// Handler represents a method that accepts the data that comes from a Fetcher and the Worker instance,
// and handles that data correctly.
// If an error is raised during the data handling, it is returned instead.
type Handler = func(data interface{}, w worker.Worker) error

// ModuleInfo represents the information that a module should provide in order to be properly registered
// to be later handled during each block creation.
type ModuleInfo struct {
	Name    string
	Fetcher Fetcher
	Handler Handler
}

// NewModule allows to easily create a new module that can be passed to RegisterModule
func NewModule(name string, fetcher Fetcher, handler Handler) ModuleInfo {
	return ModuleInfo{
		Name:    name,
		Fetcher: fetcher,
		Handler: handler,
	}
}

// RegisterModule allows to register a new module so that it can be properly handled later.
func RegisterModule(module ModuleInfo) {
	modules = append(modules, module)
}

// HandleModules handles properly all the modules that have been registered though the usage of RegisterModule.
// For each module, calls firstly its Fetcher and then with the returned data it calls the proper Handler.
// Errors are handled as soon as they are created, with the first error raised that is immediately returned.
func HandleModules(
	block *coretypes.ResultBlock, txs []djuno.Tx, vals *tmctypes.ResultValidators, w worker.Worker,
) error {
	for _, module := range modules {
		data, err := module.Fetcher(block, txs, vals, w)
		if err != nil {
			return err
		}

		err = module.Handler(data, w)
		if err != nil {
			return err
		}

		log.Debug().Str("module_name", module.Name).Msg("module handled correctly")
	}
	return nil
}
