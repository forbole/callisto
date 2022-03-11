package actions

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/rs/zerolog/log"

	"github.com/forbole/juno/v3/cmd/parse"
	"github.com/forbole/juno/v3/node/builder"
	nodeconfig "github.com/forbole/juno/v3/node/config"
	"github.com/forbole/juno/v3/node/remote"

	"github.com/spf13/cobra"

	"github.com/forbole/bdjuno/v2/cmd/actions/handlers"
	actionstypes "github.com/forbole/bdjuno/v2/cmd/actions/types"
	"github.com/forbole/bdjuno/v2/modules"
)

const (
	flagGRPC            = "grpc"
	flagRPC             = "rpc"
	flagSecure          = "secure"
	flagPort            = "port"
	flagPortPrometheus  = "prometheus-port"
	flagEablePrometheus = "enable-prometheus"
)

var (
	waitGroup sync.WaitGroup
)

// NewActionsCmd returns the Cobra command allowing to activate hasura actions
func NewActionsCmd(parseCfg *parse.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "hasura-actions",
		Short:   "Activate hasura actions",
		PreRunE: parse.ReadConfig(parseCfg),
		RunE: func(cmd *cobra.Command, args []string) error {
			parseCtx, err := parse.GetParsingContext(parseCfg)
			if err != nil {
				return err
			}

			// Get the flags values
			rpc, _ := cmd.Flags().GetString(flagRPC)
			gRPC, _ := cmd.Flags().GetString(flagGRPC)
			secure, _ := cmd.Flags().GetBool(flagSecure)
			port, _ := cmd.Flags().GetUint(flagPort)
			prometheusPort, _ := cmd.Flags().GetUint(flagPortPrometheus)
			enablePrometheus, _ := cmd.Flags().GetBool(flagEablePrometheus)

			log.Info().Str(flagRPC, rpc).Str(flagGRPC, gRPC).Bool(flagSecure, secure).
				Msg("Listening to incoming Hasura actions requests....")

			// Build a custom node config to make sure it's remote
			// TODO: Is this really necessary? Can't we use the default one?
			nodeCfg := nodeconfig.NewConfig(
				nodeconfig.TypeRemote,
				remote.NewDetails(
					remote.NewRPCConfig("hasura-actions", rpc, 100),
					remote.NewGrpcConfig(gRPC, !secure),
				),
			)

			// Build the node
			node, err := builder.BuildNode(nodeCfg, parseCtx.EncodingConfig)
			if err != nil {
				return err
			}

			// Build the sources
			sources, err := modules.BuildSources(nodeCfg, parseCtx.EncodingConfig)
			if err != nil {
				return err
			}

			// Build the worker
			context := actionstypes.NewContext(node, sources)
			worker := actionstypes.NewActionsWorker(context)

			// Register the endpoints

			// -- Bank --
			worker.RegisterHandler("/account_balance", handlers.AccountBalanceHandler)

			// -- Distribution --
			worker.RegisterHandler("/delegation_reward", handlers.DelegationRewardHandler)
			worker.RegisterHandler("/delegator_withdraw_address", handlers.DelegatorWithdrawAddressHandler)
			worker.RegisterHandler("/validator_commission_amount", handlers.ValidatorCommissionAmountHandler)

			// -- Staking Delegator --
			worker.RegisterHandler("/delegation", handlers.DelegationHandler)
			worker.RegisterHandler("/delegation_total", handlers.TotalDelegationAmountHandler)
			worker.RegisterHandler("/unbonding_delegation", handlers.UnbondingDelegationsHandler)
			worker.RegisterHandler("/unbonding_delegation_total", handlers.UnbondingDelegationsTotal)
			worker.RegisterHandler("/redelegation", handlers.RedelegationHandler)

			// -- Staking Validator --
			worker.RegisterHandler("/validator_delegations", handlers.ValidatorDelegation)
			worker.RegisterHandler("/validator_redelegations_from", handlers.ValidatorRedelegationsFromHandler)
			worker.RegisterHandler("/validator_unbonding_delegations", handlers.ValidatorUnbondingDelegationsHandler)

			// Listen for and trap any OS signal to gracefully shutdown and exit
			trapSignal(parseCtx)

			// Start the worker
			waitGroup.Add(1)
			go worker.Start(port)

			// Start Prometheus
			if enablePrometheus {
				go actionstypes.StartPrometheus(prometheusPort)
			}

			// Block main process (signal capture will call WaitGroup's Done)
			waitGroup.Wait()
			return nil
		},
	}

	cmd.Flags().String(flagRPC, "http://127.0.0.1:26657", "RPC listen address. Port required")
	cmd.Flags().String(flagGRPC, "http://127.0.0.1:9090", "GRPC listen address. Port required")
	cmd.Flags().Bool(flagSecure, false, "Activate secure connections")
	cmd.Flags().Uint(flagPort, 3000, "Port to be used to expose the service")
	cmd.Flags().Uint(flagPortPrometheus, 3001, "Port to be used to run hasura prometheus monitoring")
	cmd.Flags().Bool(flagEablePrometheus, false, "Enable prometheus monitoring")

	return cmd
}

// trapSignal will listen for any OS signal and invoke Done on the main
// WaitGroup allowing the main process to gracefully exit.
func trapSignal(parseCtx *parse.Context) {
	var sigCh = make(chan os.Signal)

	signal.Notify(sigCh, syscall.SIGTERM)
	signal.Notify(sigCh, syscall.SIGINT)

	go func() {
		sig := <-sigCh
		parseCtx.Logger.Info("caught signal; shutting down...", "signal", sig.String())
		defer parseCtx.Node.Stop()
		defer parseCtx.Database.Close()
		defer waitGroup.Done()
	}()
}
