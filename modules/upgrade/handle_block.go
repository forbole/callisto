package upgrade

import (
	"fmt"

	"github.com/forbole/juno/v3/types"

	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
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
	// refresh validator details
	return nil
}
