package modules

import (
	"github.com/forbole/callisto/v4/modules/actions"
	"github.com/forbole/callisto/v4/modules/types"

	"github.com/forbole/juno/v5/modules/pruning"
	"github.com/forbole/juno/v5/modules/telemetry"

	"github.com/forbole/callisto/v4/modules/slashing"

	jmodules "github.com/forbole/juno/v5/modules"
	"github.com/forbole/juno/v5/modules/messages"
	"github.com/forbole/juno/v5/modules/registrar"

	"github.com/forbole/callisto/v4/utils"

	"github.com/forbole/callisto/v4/database"
	"github.com/forbole/callisto/v4/modules/auth"
	"github.com/forbole/callisto/v4/modules/bank"
	"github.com/forbole/callisto/v4/modules/consensus"
	"github.com/forbole/callisto/v4/modules/distribution"
	"github.com/forbole/callisto/v4/modules/feegrant"

	dailyrefetch "github.com/forbole/callisto/v4/modules/daily_refetch"
	"github.com/forbole/callisto/v4/modules/gov"
	messagetype "github.com/forbole/callisto/v4/modules/message_type"
	"github.com/forbole/callisto/v4/modules/mint"
	"github.com/forbole/callisto/v4/modules/modules"
	"github.com/forbole/callisto/v4/modules/pricefeed"
	"github.com/forbole/callisto/v4/modules/staking"
	"github.com/forbole/callisto/v4/modules/upgrade"
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
	dailyRefetchModule := dailyrefetch.NewModule(ctx.Proxy, db)
	distrModule := distribution.NewModule(sources.DistrSource, cdc, db)
	feegrantModule := feegrant.NewModule(cdc, db)
	messagetypeModule := messagetype.NewModule(r.parser, cdc, db)
	mintModule := mint.NewModule(sources.MintSource, cdc, db)
	slashingModule := slashing.NewModule(sources.SlashingSource, cdc, db)
	stakingModule := staking.NewModule(sources.StakingSource, cdc, db)
	govModule := gov.NewModule(sources.GovSource, distrModule, mintModule, slashingModule, stakingModule, cdc, db)
	upgradeModule := upgrade.NewModule(db, stakingModule)

	return []jmodules.Module{
		messages.NewModule(r.parser, cdc, ctx.Database),
		telemetry.NewModule(ctx.JunoConfig),
		pruning.NewModule(ctx.JunoConfig, db, ctx.Logger),

		actionsModule,
		authModule,
		bankModule,
		consensusModule,
		dailyRefetchModule,
		distrModule,
		feegrantModule,
		govModule,
		mintModule,
		messagetypeModule,
		modules.NewModule(ctx.JunoConfig.Chain, db),
		pricefeed.NewModule(ctx.JunoConfig, cdc, db),
		slashingModule,
		stakingModule,
		upgradeModule,
	}
}
