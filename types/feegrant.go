package types

// FeeGrantAllowance represents the x/feegrant parameters
type FeeGrantAllowance struct {
	Grantee   string
	Granter   string
	Allowance interface{}
	Height    int64
}

// NewFeeGrantAllowance allows to build a new FeeGrantAllowance instance
func NewFeeGrantAllowance(grantee string, granter string, allowance interface{}, height int64) FeeGrantAllowance {
	return FeeGrantAllowance{
		Grantee:   grantee,
		Granter:   granter,
		Allowance: allowance,
		Height:    height,
	}
}
