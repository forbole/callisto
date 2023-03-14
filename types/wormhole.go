package types

import wormholetypes "github.com/wormhole-foundation/wormchain/x/wormhole/types"

// WormholeConfig represents the x/wormhole config
type WormholeConfig struct {
	Config *wormholetypes.Config
	Height int64
}

// NewWormholeConfig allows to build a new WormholeConfig instance
func NewWormholeConfig(config *wormholetypes.Config, height int64) *WormholeConfig {
	return &WormholeConfig{
		Config: config,
		Height: height,
	}
}
