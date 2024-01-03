package pool

import (
	"encoding/json"
	"fmt"

	tmtypes "github.com/cometbft/cometbft/types"
	"github.com/forbole/bdjuno/v4/types"

	pooltypes "github.com/KYVENetwork/chain/x/pool/types"
	"github.com/rs/zerolog/log"
)

// HandleGenesis implements modules.Module
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	log.Debug().Str("module", "pool").Msg("parsing genesis")

	// Read the genesis state
	var genState pooltypes.GenesisState
	err := m.cdc.UnmarshalJSON(appState[pooltypes.ModuleName], &genState)
	if err != nil {
		return fmt.Errorf("error while reading pool genesis data: %s", err)
	}

	var pools []types.PoolList
	for _, pool := range genState.PoolList {
		pools = append(pools, types.NewPoolList(
			pool.Id,
			pool.Name,
			pool.Runtime,
			pool.Logo,
			pool.Config,
			pool.StartKey,
			pool.CurrentKey,
			pool.CurrentSummary,
			pool.CurrentIndex,
			pool.TotalBundles,
			pool.UploadInterval,
			pool.InflationShareWeight,
			pool.MinDelegation,
			pool.MaxBundleSize,
			pool.Disabled,
			pool.Protocol,
			pool.UpgradePlan,
			pool.CurrentStorageProviderId,
			pool.CurrentCompressionId,
			doc.InitialHeight,
		))
	}

	// Save pool list
	err = m.db.SavePoolList(pools)
	if err != nil {
		return fmt.Errorf("error while storing genesis pools list: %s", err)
	}

	return nil
}
