package utils

import juno "github.com/desmos-labs/juno/types"

// IsModuleEnabled returns true if the module having the given name is enabled inside the provided configuration
func IsModuleEnabled(cfg juno.Config, moduleName string) bool {
	for _, module := range cfg.GetCosmosConfig().GetModules() {
		if module == moduleName {
			return true
		}
	}
	return false
}
