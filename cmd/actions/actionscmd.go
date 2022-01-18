package actions

import (
	"log"
	"net/http"

	"github.com/forbole/juno/v2/cmd/parse"
	"github.com/spf13/cobra"

	"github.com/forbole/bdjuno/v2/cmd/actions/handlers"
)

// NewActionsCmd returns the Cobra command allowing to activate hasura actions
func NewActionsCmd(parseCfg *parse.Config) *cobra.Command {
	return &cobra.Command{
		Use:     "hasura-actions",
		Short:   "Activate hasura actions",
		PreRunE: parse.ReadConfig(parseCfg),
		RunE: func(cmd *cobra.Command, args []string) error {

			// HTTP server for the handlers
			mux := http.NewServeMux()

			// End points:
			mux.HandleFunc("/account_balance", handlers.AccountBalance)
			mux.HandleFunc("/delegation_reward", handlers.DelegationReward)
			mux.HandleFunc("/delegation", handlers.Delegation)
			mux.HandleFunc("/unbonding_delegations", handlers.UnbondingDelegations)
			mux.HandleFunc("/validator_commission", handlers.ValidatorCommission)
			mux.HandleFunc("/redelegation", handlers.Redelegation)

			err := http.ListenAndServe(":3000", mux)
			log.Fatal(err)

			return nil
		},
	}
}
