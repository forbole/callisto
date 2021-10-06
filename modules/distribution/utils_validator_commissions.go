package distribution

import (
	"github.com/forbole/bdjuno/v2/types"

	"github.com/rs/zerolog/log"
)

// updateValidatorsCommissionAmounts updates the validators commissions amounts
func (m *Module) updateValidatorsCommissionAmounts(height int64) {
	log.Debug().Str("module", "distribution").
		Int64("height", height).
		Msg("updating validators commissions")

	validators, err := m.db.GetValidators()
	if err != nil {
		log.Error().Str("module", "distribution").Err(err).
			Int64("height", height).
			Msg("error while getting validators")
		return
	}

	if len(validators) == 0 {
		// No validators, just skip
		return
	}

	// Get all the commissions
	for _, validator := range validators {
		go m.updateValidatorCommissionAmount(height, validator)
	}
}

func (m *Module) updateValidatorCommissionAmount(height int64, validator types.Validator) {
	commission, err := m.source.ValidatorCommission(validator.GetOperator(), height)
	if err != nil {
		log.Error().Str("module", "distribution").Err(err).
			Int64("height", height).
			Str("validator", validator.GetOperator()).
			Msg("error while getting validator commission")
	}

	err = m.db.SaveValidatorCommissionAmount(types.NewValidatorCommissionAmount(
		validator.GetOperator(),
		validator.GetSelfDelegateAddress(),
		commission,
		height,
	))
	if err != nil {
		log.Error().Str("module", "distribution").Err(err).
			Int64("height", height).
			Str("validator", validator.GetOperator()).
			Msg("error while saving validator commission amounts")
	}
}
