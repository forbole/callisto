package utils

import (
	"strconv"

	juno "github.com/desmos-labs/juno/types"

	"github.com/desmos-labs/juno/client"

	grpctypes "github.com/cosmos/cosmos-sdk/types/grpc"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// WatchMethod allows to watch for a method that returns an error.
// It executes the given method in a goroutine, logging any error that might raise.
func WatchMethod(method func() error) {
	go func() {
		err := method()
		if err != nil {
			log.Error().Err(err).Send()
		}
	}()
}

// GetHeightRequestHeader returns the grpc.CallOption to query the state at a given height
func GetHeightRequestHeader(height int64) grpc.CallOption {
	header := metadata.New(map[string]string{
		grpctypes.GRPCBlockHeightHeader: strconv.FormatInt(height, 10),
	})
	return grpc.Header(&header)
}

func MustCreateGrpcConnection(cfg *juno.Config) *grpc.ClientConn {
	grpConnection, err := client.CreateGrpcConnection(cfg)
	if err != nil {
		panic(err)
	}
	return grpConnection
}
