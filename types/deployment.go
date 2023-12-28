package types

import deploymenttypes "github.com/akash-network/akash-api/go/node/deployment/v1beta3"

// DeploymentParams represents the x/deployment parameters
type DeploymentParams struct {
	deploymenttypes.Params
	Height int64
}

// NewDeploymentParams allows to build a new DeploymentParams instance
func NewDeploymentParams(params deploymenttypes.Params, height int64) *DeploymentParams {
	return &DeploymentParams{
		Params: params,
		Height: height,
	}
}
