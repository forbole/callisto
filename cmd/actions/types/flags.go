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
	cmd.PersistentFlags().StringVar(&FlagGRPC, "grpc", "https://localhost:9090", "<host>:<port> to the gRPC interface for the chain")
	cmd.PersistentFlags().BoolVar(&FlagInsecure, "insecure", false, "insecure or secure connection")
}
