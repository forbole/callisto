package types

import (
	"time"

	shieldtypes "github.com/certikfoundation/shentu/v2/x/shield/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ShieldPool represents a pool of the shield module at a given height
type ShieldPool struct {
	PoolID         uint64
	FromAddress    string
	Shield         sdk.Int
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
	poolID uint64, fromAddress string, shield sdk.Int, nativeDeposit sdk.Coins, foreignDeposit sdk.Coins, sponsor string,
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

// ShieldPurchaseList represents a purchase of the shield module at a given height
type ShieldPurchaseList struct {
	PurchaseID         uint64
	PoolID             uint64
	Purchaser          string
	DeletionTime       time.Time
	ProtectionEndTime  time.Time
	ForeignServiceFees sdk.DecCoins
	NativeServiceFees  sdk.DecCoins
	Shield             sdk.Int
	Description        string
	Height             int64
}

// NewShieldPurchaseList allows to build a new ShieldPurchaseList instance
func NewShieldPurchaseList(
	purchaseID uint64, poolID uint64, purchaser string, deletionTime time.Time, protectionEndTime time.Time,
	foreignServiceFees sdk.DecCoins, nativeServiceFees sdk.DecCoins, shield sdk.Int, description string, height int64,
) *ShieldPurchaseList {
	return &ShieldPurchaseList{
		PoolID:             poolID,
		Purchaser:          purchaser,
		DeletionTime:       deletionTime,
		ProtectionEndTime:  protectionEndTime,
		PurchaseID:         purchaseID,
		ForeignServiceFees: foreignServiceFees,
		NativeServiceFees:  nativeServiceFees,
		Shield:             shield,
		Description:        description,
		Height:             height,
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

// ShieldInfo represents the base info of the shield module at a given height
type ShieldInfo struct {
	GobalStakingPool            string
	LastUpdateTime              time.Time
	NextPoolID                  uint64
	NextPurchaseID              uint64
	OriginalStaking             []shieldtypes.OriginalStaking
	ProposalIDReimbursementPair []shieldtypes.ProposalIDReimbursementPair
	ShieldAdmin                 string
	ShieldStakingRate           string
	StakeForShields             []shieldtypes.ShieldStaking
	TotalClaimed                sdk.Int
	TotalCollateral             sdk.Int
	TotalShield                 sdk.Int
	TotalWithdrawing            sdk.Int
	Height                      int64
}

// NewShieldInfo allows to build a new ShieldInfo instance
func NewShieldInfo(gobalStakingPool string, lastUpdateTime time.Time, nextPoolID uint64, nextPurchaseID uint64,
	originalStaking []shieldtypes.OriginalStaking, proposalIDReimbursementPair []shieldtypes.ProposalIDReimbursementPair,
	shieldAdmin string, shieldStakingRate string, stakeForShields []shieldtypes.ShieldStaking, totalClaimed sdk.Int,
	totalCollateral sdk.Int, totalShield sdk.Int, totalWithdrawing sdk.Int, height int64) *ShieldInfo {
	return &ShieldInfo{
		GobalStakingPool:            gobalStakingPool,
		LastUpdateTime:              lastUpdateTime,
		NextPoolID:                  nextPoolID,
		NextPurchaseID:              nextPurchaseID,
		OriginalStaking:             originalStaking,
		ProposalIDReimbursementPair: proposalIDReimbursementPair,
		ShieldAdmin:                 shieldAdmin,
		ShieldStakingRate:           shieldStakingRate,
		StakeForShields:             stakeForShields,
		TotalClaimed:                totalClaimed,
		TotalCollateral:             totalCollateral,
		TotalShield:                 totalShield,
		TotalWithdrawing:            totalWithdrawing,
		Height:                      height,
	}
}

type ShieldServiceFees struct {
	ForeignServiceFees          sdk.DecCoins
	NativeServiceFees           sdk.DecCoins
	RemainingForeignServiceFees sdk.DecCoins
	RemainingNativeServiceFees  sdk.DecCoins
	Height                      int64
}

// NewShieldServiceFees allows to build a new ShieldServiceFees instance
func NewShieldServiceFees(foreignServiceFees sdk.DecCoins, nativeServiceFees sdk.DecCoins, remainingForeignServiceFees sdk.DecCoins, remainingNativeServiceFees sdk.DecCoins, height int64) *ShieldServiceFees {
	return &ShieldServiceFees{
		ForeignServiceFees:          foreignServiceFees,
		NativeServiceFees:           nativeServiceFees,
		RemainingForeignServiceFees: remainingForeignServiceFees,
		RemainingNativeServiceFees:  remainingNativeServiceFees,
		Height:                      height,
	}
}
