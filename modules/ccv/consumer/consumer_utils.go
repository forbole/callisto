package consumer

import (
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
		providerSelfDelegateAddress, _ := m.db.GetProviderSelfDelegateAddress(providerConsensusAddress)

		if len(providerSelfDelegateAddress) > 0 {
			_, bz, err := bech32.DecodeAndConvert(providerSelfDelegateAddress)
			if err != nil {
				fmt.Errorf("cannot decode %s address: %s", providerSelfDelegateAddress, err)
			}

			consumerSelfDelegateAddress, err := bech32.ConvertAndEncode("neutron", bz)
			if err != nil {
				fmt.Errorf("cannot decode neutron address: %s", err)
			}

			ccvValidators = append(ccvValidators, types.NewCCVValidator(index.ConsensusAddress, consumerSelfDelegateAddress, providerConsensusAddress, providerSelfDelegateAddress, height))
		} else {
			ccvValidators = append(ccvValidators, types.NewCCVValidator(index.ConsensusAddress, "", providerConsensusAddress, "", height))
		}
	}

	return m.db.StoreCCvValidators(ccvValidators)
}
