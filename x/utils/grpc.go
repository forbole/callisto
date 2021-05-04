package utils

import (
	"strconv"

	grpctypes "github.com/cosmos/cosmos-sdk/types/grpc"
	"github.com/desmos-labs/juno/client"
	juno "github.com/desmos-labs/juno/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// GetHeightRequestHeader returns the grpc.CallOption to query the state at a given height
func GetHeightRequestHeader(height int64) grpc.CallOption {
	header := metadata.New(map[string]string{
		grpctypes.GRPCBlockHeightHeader: strconv.FormatInt(height, 10),
	})
	return grpc.Header(&header)
}

// MustCreateGrpcConnection creates a new gRPC connection using the provided configuration and panics on error
func MustCreateGrpcConnection(cfg *juno.Config) *grpc.ClientConn {
	grpConnection, err := client.CreateGrpcConnection(cfg)
	if err != nil {
		panic(err)
	}
	return grpConnection
}
