package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/types/module"
	actionstypes "github.com/forbole/bdjuno/v2/cmd/actions/types"
	"github.com/forbole/bdjuno/v2/database"
	"github.com/forbole/bdjuno/v2/modules"
	"github.com/forbole/bdjuno/v2/types/config"
	"github.com/forbole/juno/v2/cmd/parse"
	"github.com/forbole/juno/v2/modules/messages"

	nodeconfig "github.com/forbole/juno/v2/node/config"
	"github.com/forbole/juno/v2/node/remote"
)

func getCtxAndSources() (*parse.Context, *modules.Sources, error) {
	parseCfg := parse.NewConfig().
		WithDBBuilder(database.Builder).
		WithEncodingConfigBuilder(config.MakeEncodingConfig([]module.BasicManager{
			simapp.ModuleBasics,
		})).
		WithRegistrar(modules.NewRegistrar(messages.JoinMessageParsers(
			messages.CosmosMessageAddressesParser,
		)))

	parseCtx, err := parse.GetParsingContext(parseCfg)
	if err != nil {
		return nil, nil, err
	}

	sources, err := modules.BuildSources(getNode(), parseCtx.EncodingConfig)
	if err != nil {
		return nil, nil, err
	}

	return parseCtx, sources, nil
}

func errorHandler(w http.ResponseWriter, err error) {
	errorObject := actionstypes.GraphQLError{
		Message: err.Error(),
	}
	errorBody, _ := json.Marshal(errorObject)
	w.WriteHeader(http.StatusBadRequest)
	w.Write(errorBody)
}

func getNode() nodeconfig.Config {
	var node nodeconfig.Config
	if actionstypes.FlagRpc == "" || actionstypes.FlagGRpc == "" {
		node = nodeconfig.DefaultConfig()
	}

	node = nodeconfig.NewConfig(
		nodeconfig.TypeRemote,
		remote.NewDetails(
			remote.NewRPCConfig("hasura-actions", actionstypes.FlagRpc, 100),
			remote.NewGrpcConfig(actionstypes.FlagGRpc, actionstypes.FlagInsecure),
		),
	)
	return node
}
