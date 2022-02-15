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

			// -- Bank --
			mux.HandleFunc("/account_balance", handlers.AccountBalance)

			// -- Distribution --
			mux.HandleFunc("/delegation_reward", handlers.DelegationReward)
			mux.HandleFunc("/delegator_withdraw_address", handlers.DelegatorWithdrawAddress)
			mux.HandleFunc("/validator_commission_amount", handlers.ValidatorCommissionAmount)

			// -- Staking Delegator --
			mux.HandleFunc("/delegation", handlers.Delegation)
			mux.HandleFunc("/delegation_total", handlers.TotalDelegationAmount)
			mux.HandleFunc("/unbonding_delegation", handlers.UnbondingDelegations)
			mux.HandleFunc("/unbonding_delegation_total", handlers.UnbondingDelegationsTotal)
			mux.HandleFunc("/redelegation", handlers.Redelegation)

			// -- Staking Validator --
			mux.HandleFunc("/validator_delegations", handlers.ValidatorDelegation)
			mux.HandleFunc("/validator_redelegations_from", handlers.ValidatorRedelegationsFrom)
			mux.HandleFunc("/validator_unbonding_delegations", handlers.ValidatorUnbondingDelegations)

			err := http.ListenAndServe(":3000", mux)
			log.Fatal(err)

			return nil
		},
	}
}
