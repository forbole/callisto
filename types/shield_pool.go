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
	poolID uint64, fromAddress string, shield sdk.Coins, nativeDeposit sdk.Coins, foreignDeposit sdk.Coins, sponsor string,
	sponsorAddress string, description string, shieldLimit sdk.Int, pause bool, height int64,
) *ShieldPool {
	return &ShieldPool{
		PoolID:         poolID,
		FromAddress:    fromAddress,
		Shield:         shield,
		NativeDeposit:  nativeDeposit,
		ForeignDeposit: foreignDeposit,
		Sponsor:        sponsor,
		SponsorAddr:    sponsorAddress,
		Description:    description,
		ShieldLimit:    shieldLimit,
		Pause:          pause,
		Height:         height,
	}
}

// ShieldPurchase represents a purchase of the shield module at a given height
type ShieldPurchase struct {
	PoolID      uint64
	FromAddress string
	Shield      sdk.Coins
	Description string
	Height      int64
}

// NewShieldPurchase allows to build a new ShieldPurchase instance
func NewShieldPurchase(
	poolID uint64, fromAddress string, shield sdk.Coins, description string, height int64,
) *ShieldPurchase {
	return &ShieldPurchase{
		PoolID:      poolID,
		FromAddress: fromAddress,
		Shield:      shield,
		Description: description,
		Height:      height,
	}
}

// --------------------------------------------------------------------------------------------------------------------

// ShieldPoolParams represents the parameters of the shield module at a given height
type ShieldPoolParams struct {
	Params shieldtypes.PoolParams
	Height int64
}

// NewSlashingParams allows to build a new ShieldPoolParams instance
func NewShieldPoolParams(params shieldtypes.PoolParams, height int64) *ShieldPoolParams {
	return &ShieldPoolParams{
		Params: params,
		Height: height,
	}
}

// --------------------------------------------------------------------------------------------------------------------

// ShieldProvider represents the provider of the shield module at a given height
type ShieldProvider struct {
	Address          string
	Collateral       int64
	DelegationBonded int64
	NativeRewards    sdk.DecCoins
	ForeignRewards   sdk.DecCoins
	TotalLocked      int64
	Withdrawing      int64
	Height           int64
}

// NewShieldProvider allows to build a new ShieldProvider instance
func NewShieldProvider(address string, collateral int64, delegationBonded int64,
	nativeRewards sdk.DecCoins, foreignRewards sdk.DecCoins, totalLocked int64,
	withdrawing int64, height int64) *ShieldProvider {
	return &ShieldProvider{
		Address:          address,
		Collateral:       collateral,
		DelegationBonded: delegationBonded,
		NativeRewards:    nativeRewards,
		ForeignRewards:   foreignRewards,
		TotalLocked:      totalLocked,
		Withdrawing:      withdrawing,
		Height:           height,
	}
}
