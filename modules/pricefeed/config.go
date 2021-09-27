package pricefeed

import "github.com/forbole/bdjuno/types"

// PricefeedConfig contains the configuration about the pricefeed module
type PricefeedConfig struct {
	Tokens []types.Token `toml:"tokens"`
}

// GetTokens returns the list of tokens for which to get the prices
func (p *PricefeedConfig) GetTokens() []types.Token {
	return p.Tokens
}
