package pricefeed

func GetDenom() (denom string) {
	for _, token := range PricefeedCfg.Tokens {
		for _, unit := range token.Units {
			if unit.Exponent == 0 {
				denom = unit.Denom
			}
		}
	}
	return denom
}
