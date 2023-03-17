package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	shieldtypes "github.com/shentufoundation/shentu/v2/x/shield/types"
)

// ShieldPool represents a pool of the shield module at a given height
type ShieldPool struct {
	PoolID             uint64
	Shield             sdk.Int
	NativeServiceFees  sdk.Coins
	ForeignServiceFees sdk.Coins
	Sponsor            string
	SponsorAddr        string
	Description        string
	ShieldLimit        sdk.Int
	Pause              bool
	Height             int64
}

// NewShieldPool allows to build a new ShieldPool instance
func NewShieldPool(
	poolID uint64, shield sdk.Int, nativeServiceFees sdk.Coins, foreignServiceFees sdk.Coins, sponsor string,
	sponsorAddress string, description string, shieldLimit sdk.Int, pause bool, height int64,
) *ShieldPool {
	return &ShieldPool{
		PoolID:             poolID,
		Shield:             shield,
		NativeServiceFees:  nativeServiceFees,
		ForeignServiceFees: foreignServiceFees,
		Sponsor:            sponsor,
		SponsorAddr:        sponsorAddress,
		Description:        description,
		ShieldLimit:        shieldLimit,
		Pause:              pause,
		Height:             height,
	}
}

// ShieldPurchase represents a purchase of the shield module at a given height
type ShieldPurchase struct {
	PoolID      uint64
	FromAddress string
	Shield      sdk.Int
	Description string
	Height      int64
}

// NewShieldPurchase allows to build a new ShieldPurchase instance
func NewShieldPurchase(
	poolID uint64, fromAddress string, shield sdk.Int, description string, height int64,
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

// ShieldPoolParams represents the pool parameters of the shield module at a given height
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

// ShieldClaimProposalParams represents the claim proposal parameters of the shield module at a given height
type ShieldClaimProposalParams struct {
	Params shieldtypes.ClaimProposalParams
	Height int64
}

// NewShieldClaimProposalParams allows to build a new ShieldPoolParams instance
func NewShieldClaimProposalParams(params shieldtypes.ClaimProposalParams, height int64) *ShieldClaimProposalParams {
	return &ShieldClaimProposalParams{
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

// ShieldWithdraw represents the withdraw of the shield module at a given height
type ShieldWithdraw struct {
	Address        string
	Amount         int64
	CompletionTime time.Time
	Height         int64
}

// NewShieldWithdraw allows to build a new ShieldWithdraw instance
func NewShieldWithdraw(address string, amount int64, completionTime time.Time,
	height int64) *ShieldWithdraw {
	return &ShieldWithdraw{
		Address:        address,
		Amount:         amount,
		CompletionTime: completionTime,
		Height:         height,
	}
}

// ShieldStatus represents the status of the shield module at a given height
type ShieldStatus struct {
	GobalStakingPool            sdk.Int
	CurrentNativeServiceFees    sdk.DecCoins
	CurrentForeignServiceFees   sdk.DecCoins
	RemainingNativeServiceFees  sdk.DecCoins
	RemainingForeignServiceFees sdk.DecCoins
	TotalCollateral             sdk.Int
	TotalShield                 sdk.Int
	TotalWithdrawing            sdk.Int
	Height                      int64
}

// NewShieldStatus allows to build a new ShieldStatus instance
func NewShieldStatus(gobalStakingPool sdk.Int, currentNativeServiceFees sdk.DecCoins,
	currentForeignServiceFees sdk.DecCoins, remainingNativeServiceFees sdk.DecCoins,
	remainingForeignServiceFees sdk.DecCoins, totalCollateral sdk.Int, totalShield sdk.Int,
	totalWithdrawing sdk.Int, height int64) *ShieldStatus {
	return &ShieldStatus{
		GobalStakingPool:            gobalStakingPool,
		CurrentNativeServiceFees:    currentNativeServiceFees,
		CurrentForeignServiceFees:   currentForeignServiceFees,
		RemainingNativeServiceFees:  remainingNativeServiceFees,
		RemainingForeignServiceFees: remainingForeignServiceFees,
		TotalCollateral:             totalCollateral,
		TotalShield:                 totalShield,
		TotalWithdrawing:            totalWithdrawing,
		Height:                      height,
	}
}
