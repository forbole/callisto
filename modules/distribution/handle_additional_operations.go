package distribution

import (
	"fmt"
)

// RunAdditionalOperations implements modules.AdditionalOperationsModule
func (m *Module) RunAdditionalOperations() error {
	return checkConfig(m.cfg)
}

func checkConfig(cfg *Config) error {
	if cfg == nil {
		return fmt.Errorf("distribution config is not set but module is enabled")
	}

	return nil
}
