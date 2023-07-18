package consumer

import (
	// "encoding/hex"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/forbole/bdjuno/v4/types"
	"github.com/rs/zerolog/log"
)

// UpdateCcvValidators updates ccv validators operator address for the given height
func (m *Module) UpdateCcvValidators(height int64) error {
	log.Debug().Str("module", "ccvconsumer").Int64("height", height).
		Msg("updating all ccv validators info")
	var ccvValidators []types.CCVValidator

	validatorsDB, err := m.db.GetValidatorsConsensusAddress()
	if err != nil {
		return fmt.Errorf("error while getting validators cons address from db: %s", err)
	}

	for _, index := range validatorsDB {
		_, providerBytes, _ := bech32.DecodeAndConvert(index.ConsensusAddress)
		providerConsensusAddress, _ := sdk.Bech32ifyAddressBytes("cosmosvalcons", providerBytes)
		ccvValidators = append(ccvValidators, types.NewCCVValidator(index.ConsensusAddress, providerConsensusAddress, height))
	}

	return m.db.StoreCCvValidators(ccvValidators)
}