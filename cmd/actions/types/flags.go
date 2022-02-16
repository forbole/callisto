package types

import (
	"github.com/spf13/cobra"
)

var (
	FlagRPC      string
	FlagGRPC     string
	FlagInsecure bool
)

// AddNodeFlagsToCmd adds node flags to hasura actions.
func AddNodeFlagsToCmd(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&FlagRPC, "rpc", "http://127.0.0.1:26657", "RPC listen address. Port required")
	cmd.PersistentFlags().StringVar(&FlagGRPC, "grpc", "http://127.0.0.1:9090", "GRPC listen address. Port required")

	cmd.PersistentFlags().BoolVar(&FlagInsecure, "insecure", false, "Allow insecure connections")
}
