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
				return fmt.Errorf("error while decoding provider %s self delegate address: %s", providerSelfDelegateAddress, err)
			}

			consumerSelfDelegateAddress, err := bech32.ConvertAndEncode("neutron", bz)
			if err != nil {
				return fmt.Errorf("error while encoding consumer self delegate address: %s", err)
			}

			consumerOperatorAddress, err := bech32.ConvertAndEncode("neutronvaloper", bz)
			if err != nil {
				return fmt.Errorf("error while encoding consumer operator address: %s", err)
			}

			providerOperatorAddress, err := bech32.ConvertAndEncode("cosmosvaloper", bz)
			if err != nil {
				return fmt.Errorf("error while encoding provider operator address: %s", err)
			}

			err = m.db.StoreCCvValidator(types.NewCCVValidator(index.ConsensusAddress,
				consumerSelfDelegateAddress,
				consumerOperatorAddress,
				providerConsensusAddress,
				providerSelfDelegateAddress,
				providerOperatorAddress,
				height))

			if err != nil {
				return fmt.Errorf("error while storing ccv validator: %s", err)
			}

		} else {
			providerLatestHeight, err := m.db.GetProviderLastBlockHeight()
			if err != nil {
				return fmt.Errorf("error while getting provider last block height: %s", err)
			}

			// query the provider consensus address from provider chain
			providerConsensusAddress, err = m.providerModule.GetValidatorProviderAddr(providerLatestHeight, "neutron-1", providerConsensusAddress)
			if err != nil {
				return fmt.Errorf("error while getting validator provider consensus address: %s", err)
			}

			providerSelfDelegateAddress, err = m.db.GetProviderSelfDelegateAddress(providerConsensusAddress)
			if err != nil {
				return fmt.Errorf("error while getting validator provider self delegate address: %s", err)
			}

			if len(providerSelfDelegateAddress) > 0 {

				_, bz, err := bech32.DecodeAndConvert(providerSelfDelegateAddress)
				if err != nil {
					return fmt.Errorf("error while decoding provider %s address: %s", providerSelfDelegateAddress, err)
				}

				consumerSelfDelegateAddress, err := bech32.ConvertAndEncode("neutron", bz)
				if err != nil {
					return fmt.Errorf("error while encoding consumer self delegate address: %s", err)
				}

				consumerOperatorAddress, err := bech32.ConvertAndEncode("neutronvaloper", bz)
				if err != nil {
					return fmt.Errorf("error while encoding consumer operator address: %s", err)
				}

				providerOperatorAddress, err := bech32.ConvertAndEncode("cosmosvaloper", bz)
				if err != nil {
					return fmt.Errorf("error while encoding provider operator address: %s", err)
				}

				err = m.db.StoreCCvValidator(types.NewCCVValidator(index.ConsensusAddress,
					consumerSelfDelegateAddress,
					consumerOperatorAddress,
					providerConsensusAddress,
					providerSelfDelegateAddress,
					providerOperatorAddress,
					height))
				if err != nil {
					return fmt.Errorf("error while storing ccv validators: %s", err)
				}
			}
		}
	}

	return nil
}
