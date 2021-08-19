package types

import "time"

// Token represents a valid token inside the chain
type Token struct {
	Name  string      `toml:"name"`
	Units []TokenUnit `toml:"units"`
}

func NewToken(name string, units []TokenUnit) Token {
	return Token{
		Name:  name,
		Units: units,
	}
}

// TokenUnit represents a unit of a token
type TokenUnit struct {
	Denom    string   `toml:"denom"`
	Exponent int      `toml:"exponent"`
	Aliases  []string `toml:"aliases"`
}

func NewTokenUnit(denom string, exponent int, aliases []string) TokenUnit {
	return TokenUnit{
		Denom:    denom,
		Exponent: exponent,
		Aliases:  aliases,
	}
}

// TokenPrice represents the price at a given moment in time of a token unit
type TokenPrice struct {
	UnitName  string
	Price     float64
	MarketCap int64
	Timestamp time.Time
}

// NewTokenPrice returns a new TokenPrice instance containing the given data
func NewTokenPrice(unitName string, price float64, marketCap int64, timestamp time.Time) TokenPrice {
	return TokenPrice{
		UnitName:  unitName,
		Price:     price,
		MarketCap: marketCap,
		Timestamp: timestamp,
	}
}
