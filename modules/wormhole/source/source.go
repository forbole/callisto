package source

import (
	wormholetypes "github.com/wormhole-foundation/wormchain/x/wormhole/types"
)

type Source interface {
	GetAllowlistAll(height int64) ([]*wormholetypes.ValidatorAllowedAddress, error)
	GetGuardianSetAll(height int64) ([]wormholetypes.GuardianSet, error)
	GetGuardianValidatorAll(height int64) ([]wormholetypes.GuardianValidator, error)
}
