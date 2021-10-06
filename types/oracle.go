package types

import oracletypes "github.com/bandprotocol/chain/v2/x/oracle/types"

// OracleParams represents the x/oracle parameters
type OracleParams struct {
	oracletypes.Params
	Height int64
}

// NewOracleParams allows to build a new OracleParams instance
func NewOracleParams(params oracletypes.Params, height int64) OracleParams {
	return OracleParams{
		Params: params,
		Height: height,
	}
}
