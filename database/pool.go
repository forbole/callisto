package database

import (
	"encoding/json"
	"fmt"

	"github.com/forbole/bdjuno/v4/types"
)

// SavePoolList saves the given pool infos inside the database
func (db *Db) SavePoolList(pools []types.PoolList) error {
	if len(pools) == 0 {
		return nil
	}

	stmt := `
INSERT INTO pool 
    (id, name, runtime, logo, config, start_key, current_key, current_summary, current_index,
		total_bundles, upload_interval, operating_cost, min_delegation, max_bundle_size,
		disabled, funders, total_funds, protocol, upgrade_plan, current_storage_provider_id, 
		current_compression_id, height)
VALUES `
	var args []interface{}

	for i, pool := range pools {
		p := i * 7

		stmt += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d,$%d, $%d, $%d, $%d, $%d, $%d, $%d,$%d, $%d, $%d, $%d, $%d, $%d, $%d,$%d),",
			p+1, p+2, p+3, p+4, p+5, p+6, p+7, p+8, p+9, p+10, p+11, p+12, p+13, p+14, p+15, p+16, p+17, p+18, p+19, p+20, p+21, p+22)

		fundersBz, err := json.Marshal(pool.Funders)
		if err != nil {
			return fmt.Errorf("error while marshaling pool funders: %s", err)
		}

		protocolBz, err := json.Marshal(pool.Protocol)
		if err != nil {
			return fmt.Errorf("error while marshaling pool protocol: %s", err)
		}

		upgradePlanBz, err := json.Marshal(pool.UpgradePlan)
		if err != nil {
			return fmt.Errorf("error while marshaling pool upgrade plan: %s", err)
		}

		args = append(args,
			pool.Id,
			pool.Name,
			pool.Runtime,
			pool.Logo,
			pool.Config,
			pool.StartKey,
			pool.CurrentKey,
			pool.CurrentSummary,
			fmt.Sprint(pool.CurrentIndex),
			fmt.Sprint(pool.TotalBundles),
			fmt.Sprint(pool.UploadInterval),
			fmt.Sprint(pool.OperatingCost),
			fmt.Sprint(pool.MinDelegation),
			fmt.Sprint(pool.MaxBundleSize),
			pool.Disabled,
			string(fundersBz),
			fmt.Sprint(pool.TotalFunds),
			string(protocolBz),
			string(upgradePlanBz),
			fmt.Sprint(pool.CurrentStorageProviderId),
			fmt.Sprint(pool.CurrentCompressionId),
			pool.Height,
		)
	}

	stmt = stmt[:len(stmt)-1] // Remove trailing ","
	stmt += `
ON CONFLICT (id) DO UPDATE 
	SET name = excluded.name,
		runtime = excluded.runtime,
		logo = excluded.logo,
		config = excluded.config,
		start_key = excluded.start_key,
		current_key = excluded.current_key,
		current_summary = excluded.current_summary,
		current_index = excluded.current_index,
		total_bundles = excluded.total_bundles,
		upload_interval = excluded.upload_interval,
		operating_cost = excluded.operating_cost,
		min_delegation = excluded.min_delegation,
		max_bundle_size = excluded.max_bundle_size,
		disabled = excluded.disabled,
		funders = excluded.funders,
		total_funds = excluded.total_funds,
		protocol = excluded.protocol,
		upgrade_plan = excluded.upgrade_plan,
		current_storage_provider_id = excluded.current_storage_provider_id,
		current_compression_id = excluded.current_compression_id,
		height = excluded.height
WHERE pool.height <= excluded.height`

	_, err := db.SQL.Exec(stmt, args...)
	if err != nil {
		return fmt.Errorf("error while storing pool list: %s", err)
	}

	return nil
}
