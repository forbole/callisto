package pricefeed

func (m *Module) GetDenom() (denom string) {
	for _, token := range m.cfg.Tokens {
		for _, unit := range token.Units {
			if unit.Exponent == 0 {
				denom = unit.Denom
			}
		}
	}
	return denom
}
