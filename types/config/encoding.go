package config

import (
	junoparams "github.com/forbole/juno/v5/types/params"

	"github.com/cosmos/cosmos-sdk/std"
	"github.com/cosmos/cosmos-sdk/types/module"
)

// MakeEncodingConfig creates an EncodingConfig to properly handle all the messages
func MakeEncodingConfig(managers []module.BasicManager) func() junoparams.EncodingConfig {
	return func() junoparams.EncodingConfig {
		encodingConfig := junoparams.MakeTestEncodingConfig()
		std.RegisterLegacyAminoCodec(encodingConfig.Amino)
		std.RegisterInterfaces(encodingConfig.InterfaceRegistry)
		manager := mergeBasicManagers(managers)
		manager.RegisterLegacyAminoCodec(encodingConfig.Amino)
		manager.RegisterInterfaces(encodingConfig.InterfaceRegistry)
		return encodingConfig
	}
}

// mergeBasicManagers merges the given managers into a single module.BasicManager
func mergeBasicManagers(managers []module.BasicManager) module.BasicManager {
	var union = module.BasicManager{}
	for _, manager := range managers {
		for k, v := range manager {
			union[k] = v
		}
	}
	return union
}
