package upgrade

import (
	"fmt"

	"github.com/forbole/juno/v5/types"

	tmctypes "github.com/cometbft/cometbft/rpc/core/types"
)

// HandleBlock implements modules.Module
func (m *Module) HandleBlock(
	b *tmctypes.ResultBlock, _ *tmctypes.ResultBlockResults, _ []*types.Tx, _ *tmctypes.ResultValidators,
) error {
	err := m.refreshDataUponSoftwareUpgrade(b.Block.Height)
	if err != nil {
		return fmt.Errorf("error while refreshing data upon software upgrade: %s", err)
	}

	return nil
}

func (m *Module) refreshDataUponSoftwareUpgrade(height int64) error {
	exist, err := m.db.CheckSoftwareUpgradePlan(height)
	if err != nil {
		return fmt.Errorf("error while checking software upgrade plan existence: %s", err)
	}
	if !exist {
		return nil
	}

	// Refresh validator infos
	err = m.stakingModule.RefreshAllValidatorInfos(height)
	if err != nil {
		return fmt.Errorf("error while refreshing validator infos upon software upgrade: %s", err)
	}

	// Delete plan after refreshing data
	err = m.db.TruncateSoftwareUpgradePlan(height)
	if err != nil {
		return fmt.Errorf("error while truncating software upgrade plan: %s", err)
	}

	return nil
}
