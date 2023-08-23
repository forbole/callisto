package modules

import (
	"github.com/forbole/bdjuno/v3/modules/actions"
	"github.com/forbole/bdjuno/v3/modules/types"

	"github.com/forbole/juno/v3/modules/pruning"
	"github.com/forbole/juno/v3/modules/telemetry"

	"github.com/forbole/bdjuno/v3/modules/slashing"

	"github.com/cosmos/cosmos-sdk/codec"
	"fmt"
	"os"

	accountstypes "git.ooo.ua/vipcoin/chain/x/accounts/types"
	assetstypes "git.ooo.ua/vipcoin/chain/x/assets/types"
	bankingtypes "git.ooo.ua/vipcoin/chain/x/banking/types"
	walletstypes "git.ooo.ua/vipcoin/chain/x/wallets/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	jmodules "github.com/forbole/juno/v3/modules"
	"github.com/forbole/juno/v3/modules/messages"
	"github.com/forbole/juno/v3/modules/registrar"

	"github.com/forbole/bdjuno/v3/utils"

	"github.com/forbole/bdjuno/v3/database"
	"github.com/forbole/bdjuno/v3/modules/auth"
	"github.com/forbole/bdjuno/v3/modules/bank"
	"github.com/forbole/bdjuno/v3/modules/consensus"
	"github.com/forbole/bdjuno/v3/modules/distribution"
	"github.com/forbole/bdjuno/v3/modules/feegrant"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	jmodules "github.com/forbole/juno/v2/modules"
	"github.com/forbole/juno/v2/modules/messages"
	"github.com/forbole/juno/v2/modules/pruning"
	"github.com/forbole/juno/v2/modules/registrar"
	"github.com/forbole/juno/v2/modules/telemetry"
	nodeconfig "github.com/forbole/juno/v2/node/config"
	"github.com/forbole/juno/v2/node/local"
	"github.com/forbole/juno/v2/node/remote"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/forbole/bdjuno/v3/modules/gov"
	"github.com/forbole/bdjuno/v3/modules/mint"
	"github.com/forbole/bdjuno/v3/modules/modules"
	"github.com/forbole/bdjuno/v3/modules/pricefeed"
	"github.com/forbole/bdjuno/v3/modules/staking"
	"github.com/forbole/bdjuno/v2/database"
	"github.com/forbole/bdjuno/v2/modules/auth"
	"github.com/forbole/bdjuno/v2/modules/bank"
	banksource "github.com/forbole/bdjuno/v2/modules/bank/source"
	localbanksource "github.com/forbole/bdjuno/v2/modules/bank/source/local"
	remotebanksource "github.com/forbole/bdjuno/v2/modules/bank/source/remote"
	"github.com/forbole/bdjuno/v2/modules/consensus"
	"github.com/forbole/bdjuno/v2/modules/distribution"
	distrsource "github.com/forbole/bdjuno/v2/modules/distribution/source"
	localdistrsource "github.com/forbole/bdjuno/v2/modules/distribution/source/local"
	remotedistrsource "github.com/forbole/bdjuno/v2/modules/distribution/source/remote"
	"github.com/forbole/bdjuno/v2/modules/gov"
	govsource "github.com/forbole/bdjuno/v2/modules/gov/source"
	localgovsource "github.com/forbole/bdjuno/v2/modules/gov/source/local"
	remotegovsource "github.com/forbole/bdjuno/v2/modules/gov/source/remote"
	"github.com/forbole/bdjuno/v2/modules/mint"
	mintsource "github.com/forbole/bdjuno/v2/modules/mint/source"
	localmintsource "github.com/forbole/bdjuno/v2/modules/mint/source/local"
	remotemintsource "github.com/forbole/bdjuno/v2/modules/mint/source/remote"
	"github.com/forbole/bdjuno/v2/modules/modules"
	"github.com/forbole/bdjuno/v2/modules/overgold"
	overgoldAccountsSource "github.com/forbole/bdjuno/v2/modules/overgold/chain/accounts/source"
	remoteOvergoldAccountsSource "github.com/forbole/bdjuno/v2/modules/overgold/chain/accounts/source/remote"
	overgoldAssetsSource "github.com/forbole/bdjuno/v2/modules/overgold/chain/assets/source"
	remoteOvergoldAssetsSource "github.com/forbole/bdjuno/v2/modules/overgold/chain/assets/source/remote"
	overgoldBankingSource "github.com/forbole/bdjuno/v2/modules/overgold/chain/banking/source"
	remoteOvergoldBankingSource "github.com/forbole/bdjuno/v2/modules/overgold/chain/banking/source/remote"
	overgoldWalletsSource "github.com/forbole/bdjuno/v2/modules/overgold/chain/wallets/source"
	remoteOvergoldWalletsSource "github.com/forbole/bdjuno/v2/modules/overgold/chain/wallets/source/remote"
	"github.com/forbole/bdjuno/v2/modules/pricefeed"
	"github.com/forbole/bdjuno/v2/modules/slashing"
	slashingsource "github.com/forbole/bdjuno/v2/modules/slashing/source"
	localslashingsource "github.com/forbole/bdjuno/v2/modules/slashing/source/local"
	remoteslashingsource "github.com/forbole/bdjuno/v2/modules/slashing/source/remote"
	"github.com/forbole/bdjuno/v2/modules/staking"
	stakingsource "github.com/forbole/bdjuno/v2/modules/staking/source"
	localstakingsource "github.com/forbole/bdjuno/v2/modules/staking/source/local"
	remotestakingsource "github.com/forbole/bdjuno/v2/modules/staking/source/remote"
	"github.com/forbole/bdjuno/v2/utils"
)

// UniqueAddressesParser returns a wrapper around the given parser that removes all duplicated addresses
func UniqueAddressesParser(parser messages.MessageAddressesParser) messages.MessageAddressesParser {
	return func(cdc codec.Codec, msg sdk.Msg) ([]string, error) {
		addresses, err := parser(cdc, msg)
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
	cdc := ctx.EncodingConfig.Marshaler
	db := database.Cast(ctx.Database)

	sources, err := types.BuildSources(ctx.JunoConfig.Node, ctx.EncodingConfig)
	if err != nil {
		panic(err)
	}

	actionsModule := actions.NewModule(ctx.JunoConfig, ctx.EncodingConfig)
	authModule := auth.NewModule(r.parser, cdc, db)
	bankModule := bank.NewModule(r.parser, sources.BankSource, cdc, db)
	consensusModule := consensus.NewModule(db)
	distrModule := distribution.NewModule(sources.DistrSource, cdc, db)
	feegrantModule := feegrant.NewModule(cdc, db)
	mintModule := mint.NewModule(sources.MintSource, cdc, db)
	slashingModule := slashing.NewModule(sources.SlashingSource, cdc, db)
	stakingModule := staking.NewModule(sources.StakingSource, slashingModule, cdc, db)
	govModule := gov.NewModule(sources.GovSource, authModule, distrModule, mintModule, slashingModule, stakingModule, cdc, db)

	overgoldModules := overgold.NewModule(
		cdc,
		db,
		ctx.Proxy,
		ctx.Logger,

		sources.OvergoldAccountsSource,
		sources.OvergoldWalletsSource,
		sources.OvergoldBankingSource,
		sources.OvergoldAssetsSource,
	)

	return []jmodules.Module{
		messages.NewModule(r.parser, cdc, ctx.Database),
		telemetry.NewModule(ctx.JunoConfig),
		pruning.NewModule(ctx.JunoConfig, db, ctx.Logger),

		actionsModule,
		authModule,
		bankModule,
		consensusModule,
		distrModule,
		feegrantModule,
		govModule,
		mintModule,
		modules.NewModule(ctx.JunoConfig.Chain, db),
		pricefeed.NewModule(ctx.JunoConfig, cdc, db),
		slashingModule,
		stakingModule,

		overgoldModules,
	}
}
