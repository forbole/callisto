package types

import "time"

// Token represents a valid token inside the chain
type Token struct {
	Name  string      `yaml:"name"`
	Units []TokenUnit `yaml:"units"`
}

func NewToken(name string, units []TokenUnit) Token {
	return Token{
		Name:  name,
		Units: units,
	}
}

// TokenUnit represents a unit of a token
type TokenUnit struct {
	Denom    string   `yaml:"denom"`
	Exponent int      `yaml:"exponent"`
	Aliases  []string `yaml:"aliases,omitempty"`
	PriceID  string   `yaml:"price_id,omitempty"`
}

func NewTokenUnit(denom string, exponent int, aliases []string, priceID string) TokenUnit {
	return TokenUnit{
		Denom:    denom,
		Exponent: exponent,
		Aliases:  aliases,
		PriceID:  priceID,
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
