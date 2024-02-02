package modules

import (
	"github.com/forbole/bdjuno/v4/modules/actions"
	"github.com/forbole/bdjuno/v4/modules/types"

	"github.com/forbole/juno/v5/modules/pruning"
	"github.com/forbole/juno/v5/modules/telemetry"

	"github.com/forbole/bdjuno/v4/modules/slashing"

	"github.com/forbole/bdjuno/v4/utils"
	jmodules "github.com/forbole/juno/v5/modules"
	"github.com/forbole/juno/v5/modules/messages"
	"github.com/forbole/juno/v5/modules/registrar"

	"github.com/forbole/bdjuno/v4/database"
	"github.com/forbole/bdjuno/v4/modules/auth"
	"github.com/forbole/bdjuno/v4/modules/bank"
	ccvconsumer "github.com/forbole/bdjuno/v4/modules/ccv/consumer"
	ccvprovider "github.com/forbole/bdjuno/v4/modules/ccv/provider"
	"github.com/forbole/bdjuno/v4/modules/consensus"
	dailyrefetch "github.com/forbole/bdjuno/v4/modules/daily_refetch"
	"github.com/forbole/bdjuno/v4/modules/feegrant"
	messagetype "github.com/forbole/bdjuno/v4/modules/message_type"
	"github.com/forbole/bdjuno/v4/modules/modules"
	"github.com/forbole/bdjuno/v4/modules/pricefeed"
	"github.com/forbole/bdjuno/v4/modules/wasm"
	juno "github.com/forbole/juno/v5/types"
)

// UniqueAddressesParser returns a wrapper around the given parser that removes all duplicated addresses
func UniqueAddressesParser(parser messages.MessageAddressesParser) messages.MessageAddressesParser {
	return func(tx *juno.Tx) ([]string, error) {
		addresses, err := parser(tx)
		if err != nil {
			return nil, err
		}

		return utils.RemoveDuplicateValues(addresses), nil
	}
}

// --------------------------------------------------------------------------------------------------------------------

var (
	_ registrar.Registrar = &Registrar{}
)

// Registrar represents the modules.Registrar that allows to register all modules that are supported by BigDipper
type Registrar struct {
	parser messages.MessageAddressesParser
}

// NewRegistrar allows to build a new Registrar instance
func NewRegistrar(parser messages.MessageAddressesParser) *Registrar {
	return &Registrar{
		parser: UniqueAddressesParser(parser),
	}
}

// BuildModules implements modules.Registrar
func (r *Registrar) BuildModules(ctx registrar.Context) jmodules.Modules {
	cdc := ctx.EncodingConfig.Codec
	db := database.Cast(ctx.Database)

	sources, err := types.BuildSources(ctx.JunoConfig.Node, ctx.EncodingConfig)
	if err != nil {
		panic(err)
	}

	actionsModule := actions.NewModule(ctx.JunoConfig, ctx.EncodingConfig)
	authModule := auth.NewModule(r.parser, cdc, db)
	bankModule := bank.NewModule(r.parser, sources.BankSource, cdc, db)
	consensusModule := consensus.NewModule(db)
	ccvConsumerModule := ccvconsumer.NewModule(sources.ProviderSource, cdc, db)
	ccvProviderModule := ccvprovider.NewModule(sources.ProviderSource, cdc, db)
	dailyRefetchModule := dailyrefetch.NewModule(ctx.Proxy, db)
	feegrantModule := feegrant.NewModule(cdc, db)
	messagetypeModule := messagetype.NewModule(r.parser, cdc, db)
	slashingModule := slashing.NewModule(sources.SlashingSource, cdc, db)
	wasmModule := wasm.NewModule(sources.WasmSource, cdc, db)

	return []jmodules.Module{
		actionsModule,
		messages.NewModule(r.parser, cdc, ctx.Database),
		telemetry.NewModule(ctx.JunoConfig),
		pruning.NewModule(ctx.JunoConfig, db, ctx.Logger),
		authModule,
		bankModule,
		consensusModule,
		ccvConsumerModule,
		ccvProviderModule,
		dailyRefetchModule,
		feegrantModule,
		messagetypeModule,
		modules.NewModule(ctx.JunoConfig.Chain, db),
		pricefeed.NewModule(ctx.JunoConfig, cdc, db),
		slashingModule,
		wasmModule,
	}
}
