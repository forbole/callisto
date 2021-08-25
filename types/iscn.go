package types

iscntypes "https://github.com/likecoin/likecoin-chain/tree/master/x/iscn/types"

// IscnParams represents the x/iscn parameters
type IscnParams struct {
	iscntypes.Params
	Height int64
}

// NewIscnParams allows to build a new IscnParams instance
func NewIscnParams(params iscntypes.Params, height int64) IscnParams {
	return IscnParams{
		Params: params,
		Height: height,
	}
}
