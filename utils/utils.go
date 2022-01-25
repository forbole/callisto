package utils

import (
	"context"
	"fmt"
	"strconv"

	grpctypes "github.com/cosmos/cosmos-sdk/types/grpc"
	"github.com/forbole/juno/v2/cmd/parse"
	"google.golang.org/grpc/metadata"
)

// RemoveDuplicateValues removes the duplicated values from the given slice
func RemoveDuplicateValues(slice []string) []string {
	keys := make(map[string]bool)
	var list []string

	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// GetHeightRequestContext adds the height to the context for queries
func GetHeightRequestContext(context context.Context, height int64) context.Context {
	return metadata.AppendToOutgoingContext(
		context,
		grpctypes.GRPCBlockHeightHeader,
		strconv.FormatInt(height, 10),
	)
}

// GetHeight uses the lastest height when the input height is empty from graphql request
func GetHeight(parseCtx *parse.Context, inputHeight int64) (int64, error) {
	if inputHeight == 0 {
		latestHeight, err := parseCtx.Node.LatestHeight()
		if err != nil {
			return 0, fmt.Errorf("error while getting chain latest block height: %s", err)
		}
		return latestHeight, nil
	}

	return inputHeight, nil
}
