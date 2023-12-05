package config

import (
	"cosmossdk.io/simapp/params"
	allowed "git.ooo.ua/vipcoin/ovg-chain/x/allowed/types"
	core "git.ooo.ua/vipcoin/ovg-chain/x/core/types"
	feeexcludertypes "git.ooo.ua/vipcoin/ovg-chain/x/feeexcluder/types"
	referral "git.ooo.ua/vipcoin/ovg-chain/x/referral/types"
	stake "git.ooo.ua/vipcoin/ovg-chain/x/stake/types"
	"github.com/cosmos/cosmos-sdk/std"
	"github.com/cosmos/cosmos-sdk/types/module"
)

// MakeEncodingConfig creates an EncodingConfig to properly handle all the messages
func MakeEncodingConfig(managers []module.BasicManager) func() params.EncodingConfig {
	return func() params.EncodingConfig {
		encodingConfig := params.MakeTestEncodingConfig()
		std.RegisterLegacyAminoCodec(encodingConfig.Amino)
		std.RegisterInterfaces(encodingConfig.InterfaceRegistry)
		manager := mergeBasicManagers(managers)
		manager.RegisterLegacyAminoCodec(encodingConfig.Amino)
		manager.RegisterInterfaces(encodingConfig.InterfaceRegistry)

		// custom modules
		allowed.RegisterInterfaces(encodingConfig.InterfaceRegistry)
		core.RegisterInterfaces(encodingConfig.InterfaceRegistry)
		feeexcludertypes.RegisterInterfaces(encodingConfig.InterfaceRegistry)
		referral.RegisterInterfaces(encodingConfig.InterfaceRegistry)
		stake.RegisterInterfaces(encodingConfig.InterfaceRegistry)

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
