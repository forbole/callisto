package stakers

import (
	"fmt"

	stakersquerytypes "github.com/KYVENetwork/chain/x/query/types"

	"github.com/forbole/bdjuno/v4/modules/staking/keybase"
	"github.com/forbole/bdjuno/v4/types"
)

// UpdateProtocolValidators refreshes protocol validators info in database
func (m *Module) UpdateProtocolValidatorsInfo(height int64) error {
	protocolValidators, err := m.source.Stakers(height)
	if err != nil {
		return fmt.Errorf("error while querying protocol validators list %s", err)
	}

	err = m.UpdateProtocolValidators(protocolValidators, height)
	if err != nil {
		return fmt.Errorf("error while updating protocol validators list %s", err)
	}

	err = m.UpdateProtocolValidatorsDelegation(protocolValidators, height)
	if err != nil {
		return fmt.Errorf("error while updating protocol validators delegation %s", err)
	}

	err = m.UpdateProtocolValidatorsPools(protocolValidators, height)
	if err != nil {
		return fmt.Errorf("error while updating protocol validators pools %s", err)
	}

	return nil
}

func (m *Module) UpdateProtocolValidators(protocolValidators []stakersquerytypes.FullStaker, height int64) error {
	var protocolVals []types.ProtocolValidator

	for _, val := range protocolValidators {
		protocolVals = append(protocolVals, types.NewProtocolValidator(val.Address, height))
	}

	err := m.db.SaveProtocolValidators(protocolVals)
	if err != nil {
		return fmt.Errorf("error while saving protocol validators: %s", err)
	}

	return nil
}

func (m *Module) UpdateProtocolValidatorsDelegation(protocolValidators []stakersquerytypes.FullStaker, height int64) error {
	var protocolValsDelegation []types.ProtocolValidatorDelegation

	for _, val := range protocolValidators {
		protocolValsDelegation = append(protocolValsDelegation, types.NewProtocolValidatorDelegation(val.Address, val.SelfDelegation, val.TotalDelegation, val.DelegatorCount, height))
	}

	err := m.db.SaveProtocolValidatorsDelegation(protocolValsDelegation)
	if err != nil {
		return fmt.Errorf("error while saving protocol validators delegation info: %s", err)
	}

	return nil
}

func (m *Module) UpdateProtocolValidatorsPools(protocolValidators []stakersquerytypes.FullStaker, height int64) error {
	var protocolValsPools []types.ProtocolValidatorPool

	for _, val := range protocolValidators {
		if val.Pools != nil {
			for _, pool := range val.Pools {
				protocolValsPools = append(protocolValsPools, types.NewProtocolValidatorPool(val.Address, pool.Valaddress, pool.Balance, pool.Pool.Name, height))
			}
		}
	}

	err := m.db.SaveProtocolValidatorsPool(protocolValsPools)
	if err != nil {
		return fmt.Errorf("error while saving protocol validators pool info: %s", err)
	}
	return nil
}

// ---------------------------------------------------------------------------------------------------------------------

func (m *Module) UpdateProtocolValidatorsCommission(height int64) error {
	protocolValidators, err := m.source.Stakers(height)
	if err != nil {
		return fmt.Errorf("error while querying protocol validators list %s", err)
	}

	var protocolValsCommission []types.ProtocolValidatorCommission

	for _, val := range protocolValidators {
		protocolValsCommission = append(protocolValsCommission, types.NewProtocolValidatorCommission(val.Address, val.Metadata.Commission, val.Metadata.PendingCommissionChange, val.SelfDelegation, height))
	}

	err = m.db.SaveProtocolValidatorsCommission(protocolValsCommission)
	if err != nil {
		return fmt.Errorf("error while saving protocol validators commission: %s", err)
	}

	return nil
}

func (m *Module) UpdateProtocolValidatorsDescription(height int64) error {
	protocolValidators, err := m.source.Stakers(height)
	if err != nil {
		return fmt.Errorf("error while querying protocol validators list %s", err)
	}

	var protocolValsDescription []types.ProtocolValidatorDescription

	for _, val := range protocolValidators {
		avatarURL, err := keybase.GetAvatarURL(val.Metadata.Identity)
		if err != nil {
			return fmt.Errorf("error while getting Avatar URL: %s", err)
		}
		protocolValsDescription = append(protocolValsDescription, types.NewProtocolValidatorDescription(val.Address, val.Metadata, avatarURL, height))
	}

	err = m.db.SaveProtocolValidatorsDescription(protocolValsDescription)
	if err != nil {
		return fmt.Errorf("error while saving protocol validators description: %s", err)
	}

	return nil
}
