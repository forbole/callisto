package types

import (
	shieldtypes "github.com/certikfoundation/shentu/v2/x/shield/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ShieldPool represents a pool of the shield module at a given height
type ShieldPool struct {
	PoolID         uint64
	FromAddress    string
	Shield         sdk.Coins
	NativeDeposit  sdk.Coins
	ForeignDeposit sdk.Coins
	Sponsor        string
	SponsorAddr    string
	Description    string
	ShieldLimit    sdk.Int
	Pause          bool
	Height         int64
}

// NewShieldPool allows to build a new ShieldPool instance
func NewShieldPool(
	poolID uint64, fromAddress string, shield sdk.Coins, deposit shieldtypes.MixedCoins, sponsor string,
	sponsorAddress string, description string, shieldLimit sdk.Int, pause bool, height int64,
) *ShieldPool {
	return &ShieldPool{
		PoolID:         poolID,
		FromAddress:    fromAddress,
		Shield:         shield,
		NativeDeposit:  deposit.Native,
		ForeignDeposit: deposit.Foreign,
		Sponsor:        sponsor,
		SponsorAddr:    sponsorAddress,
		Description:    description,
		ShieldLimit:    shieldLimit,
		Pause:          pause,
		Height:         height,
	}
}
