package types

type ValidatorAddressRow struct {
	ConsensusAddress string `db:"consensus_address"`
}

func NewValidatorAddressRow(consensusAddress string) ValidatorAddressRow {
	return ValidatorAddressRow{
		ConsensusAddress: consensusAddress,
	}
}
