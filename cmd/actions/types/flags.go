package types

import (
	"github.com/spf13/cobra"
)

var (
	FlagRpc      string
	FlagGRpc     string
	FlagInsecure bool
)

// AddNodeFlagsToCmd adds node flags to hasura actions.
func AddNodeFlagsToCmd(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&FlagRpc, "rpc", "https://localhost:26657", "<host>:<port> to the RPC interface for the chain")
	cmd.PersistentFlags().StringVar(&FlagGRpc, "grpc", "https://localhost:9090", "<host>:<port> to the gRPC interface for the chain")
	cmd.PersistentFlags().BoolVar(&FlagInsecure, "insecure", false, "insecure or secure connection")
}
