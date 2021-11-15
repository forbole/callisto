package slashing

import (
	"fmt"

	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	juno "github.com/forbole/juno/v2/types"

	"github.com/rs/zerolog/log"
	abci "github.com/tendermint/tendermint/abci/types"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// HandleBlock implements BlockModule
func (m *Module) HandleBlock(
	block *tmctypes.ResultBlock, results *tmctypes.ResultBlockResults, _ []*juno.Tx, _ *tmctypes.ResultValidators,
) error {
	// Update the signing infos
	err := m.updateSigningInfo(block.Block.Height)
	if err != nil {
		return fmt.Errorf("error while updating signing info: %s", err)
	}

	// Update the delegations of the slashed validators
	err = m.updateSlashedDelegations(block.Block.Height, results.BeginBlockEvents)
	if err != nil {
		return fmt.Errorf("error while updating slashes: %s", err)
	}

	return nil
}

// updateSigningInfo reads from the LCD the current staking pool and stores its value inside the database
func (m *Module) updateSigningInfo(height int64) error {
	log.Debug().Str("module", "slashing").Int64("height", height).Msg("updating signing info")

	signingInfos, err := m.getSigningInfos(height)
	if err != nil {
		return err
	}

	return m.db.SaveValidatorsSigningInfos(signingInfos)
}

// updateSlashedDelegations updates all the delegations of the slashed validators
func (m *Module) updateSlashedDelegations(height int64, beginBlockEvents []abci.Event) error {
	events := juno.FindEventsByType(beginBlockEvents, slashingtypes.EventTypeSlash)

	for _, event := range events {
		addressAttr, err := juno.FindAttributeByKey(event, slashingtypes.AttributeKeyAddress)
		if err != nil {
			return err
		}

		consAddr := string(addressAttr.Value)
		valOperAddr, err := m.db.GetValidatorOperatorAddress(consAddr)
		if err != nil {
			return fmt.Errorf("error while getting validator operator address; make sure the slashing module is listed after the staking module: %s", err)
		}

		err = m.stakingModule.RefreshValidatorDelegations(height, valOperAddr.String())
		if err != nil {
			return fmt.Errorf("error while refreshing validator delegations for validator %s: %s", consAddr, err)
		}
	}

	return nil
}
