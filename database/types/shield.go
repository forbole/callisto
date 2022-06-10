package types

import "time"

type ShieldPoolRow struct {
	PoolID             int64    `db:"pool_id"`
	Shield             string   `db:"shield"`
	NativeServiceFees  *DbCoins `db:"native_service_fees"`
	ForeignServiceFees *DbCoins `db:"foreign_service_fees"`
	Sponsor            string   `db:"sponsor"`
	SponsorAddr        string   `db:"sponsor_address"`
	Description        string   `db:"description"`
	ShieldLimit        string   `db:"shield_limit"`
	Pause              bool     `db:"pause"`
	Height             int64    `db:"height"`
}

// Equal tells whether v and w represent the same rows
func (v ShieldPoolRow) Equal(w ShieldPoolRow) bool {
	return v.PoolID == w.PoolID &&
		v.Shield == w.Shield &&
		v.NativeServiceFees.Equal(w.NativeServiceFees) &&
		v.ForeignServiceFees.Equal(w.ForeignServiceFees) &&
		v.Sponsor == w.Sponsor &&
		v.SponsorAddr == w.SponsorAddr &&
		v.Description == w.Description &&
		v.ShieldLimit == w.ShieldLimit &&
		v.Pause == w.Pause &&
		v.Height == w.Height
}

// NewShieldPoolRow allows to build a new ShieldPoolRow
func NewShieldPoolRow(
	poolID int64,
	shield string,
	nativeServiceFees DbCoins,
	foreignServiceFees DbCoins,
	sponsor string,
	sponsorAddress string,
	description string,
	shieldLimit string,
	pause bool,
	height int64,
) ShieldPoolRow {
	return ShieldPoolRow{
		PoolID:             poolID,
		Shield:             shield,
		NativeServiceFees:  &nativeServiceFees,
		ForeignServiceFees: &foreignServiceFees,
		Sponsor:            sponsor,
		SponsorAddr:        sponsorAddress,
		Description:        description,
		ShieldLimit:        shieldLimit,
		Pause:              pause,
		Height:             height,
	}
}

// ----------------------------------------------------------------

type ShieldProviderRow struct {
	Address          string      `db:"address"`
	Collateral       int64       `db:"collateral"`
	DelegationBonded int64       `db:"delegation_bonded"`
	NativeRewards    *DbDecCoins `db:"native_rewards"`
	ForeignRewards   *DbDecCoins `db:"foreign_rewards"`
	TotalLocked      int64       `db:"total_locked"`
	Withdrawing      int64       `db:"withdrawing"`
	Height           int64       `db:"height"`
}

// Equal tells whether v and w represent the same rows
func (v ShieldProviderRow) Equal(w ShieldProviderRow) bool {
	return v.Address == w.Address &&
		v.Collateral == w.Collateral &&
		v.DelegationBonded == w.DelegationBonded &&
		v.NativeRewards.Equal(w.NativeRewards) &&
		v.ForeignRewards.Equal(w.ForeignRewards) &&
		v.TotalLocked == w.TotalLocked &&
		v.Withdrawing == w.Withdrawing &&
		v.Height == w.Height
}

// NewShieldProviderRow allows to build a new ShieldProviderRow
func NewShieldProviderRow(
	address string,
	collateral int64,
	delegationBonded int64,
	nativeRewards DbDecCoins,
	foreignRewards DbDecCoins,
	totalLocked int64,
	withdrawing int64,
	height int64,
) ShieldProviderRow {
	return ShieldProviderRow{
		Address:          address,
		Collateral:       collateral,
		DelegationBonded: delegationBonded,
		NativeRewards:    &nativeRewards,
		ForeignRewards:   &foreignRewards,
		TotalLocked:      totalLocked,
		Withdrawing:      withdrawing,
		Height:           height,
	}
}

// ----------------------------------------------------------------

type ShieldPurchaseRow struct {
	PoolID      uint64 `db:"pool_id"`
	FromAddress string `db:"purchaser"`
	Shield      string `db:"shield"`
	Description string `db:"description"`
	Height      int64  `db:"height"`
}

// Equal tells whether v and w represent the same rows
func (v ShieldPurchaseRow) Equal(w ShieldPurchaseRow) bool {
	return v.PoolID == w.PoolID &&
		v.FromAddress == w.FromAddress &&
		v.Shield == w.Shield &&
		v.Description == w.Description &&
		v.Height == w.Height
}

// NewShieldPurchaseRow allows to build a new ShieldPurchaseRow
func NewShieldPurchaseRow(
	poolID uint64,
	fromAddress string,
	shield string,
	description string,
	height int64,
) ShieldPurchaseRow {
	return ShieldPurchaseRow{
		PoolID:      poolID,
		FromAddress: fromAddress,
		Shield:      shield,
		Description: description,
		Height:      height,
	}
}

// ----------------------------------------------------------------

type ShieldWithdrawRow struct {
	Address        string    `db:"address"`
	Amount         int64     `db:"amount"`
	CompletionTime time.Time `db:"completion_time"`
	Height         int64     `db:"height"`
}

// Equal tells whether v and w represent the same rows
func (v ShieldWithdrawRow) Equal(w ShieldWithdrawRow) bool {
	return v.Address == w.Address &&
		v.Amount == w.Amount &&
		v.CompletionTime.Equal(w.CompletionTime) &&
		v.Height == w.Height
}

// NewShieldWithdrawRow allows to build a new ShieldWithdrawRow
func NewShieldWithdrawRow(
	address string,
	amount int64,
	completionTime time.Time,
	height int64,
) ShieldWithdrawRow {
	return ShieldWithdrawRow{
		Address:        address,
		Amount:         amount,
		CompletionTime: completionTime,
		Height:         height,
	}
}

// ----------------------------------------------------------------

type ShieldStatusRow struct {
	OneRowID                    bool        `db:"one_row_id"`
	GobalStakingPool            int64       `db:"global_staking_pool"`
	CurrentNativeServiceFees    *DbDecCoins `db:"current_native_service_fees"`
	CurrentForeignServiceFees   *DbDecCoins `db:"current_foreign_service_fees"`
	RemainingNativeServiceFees  *DbDecCoins `db:"remaining_native_service_fees"`
	RemainingForeignServiceFees *DbDecCoins `db:"remaining_foreign_service_fees"`
	TotalCollateral             int64       `db:"total_collateral"`
	TotalShield                 int64       `db:"total_shield"`
	TotalWithdrawing            int64       `db:"total_withdrawing"`
	Height                      int64       `db:"height"`
}

// Equal tells whether v and w represent the same rows
func (v ShieldStatusRow) Equal(w ShieldStatusRow) bool {
	return v.GobalStakingPool == w.GobalStakingPool &&
		v.CurrentNativeServiceFees.Equal(w.CurrentNativeServiceFees) &&
		v.CurrentForeignServiceFees.Equal(w.CurrentForeignServiceFees) &&
		v.RemainingNativeServiceFees.Equal(w.RemainingNativeServiceFees) &&
		v.RemainingForeignServiceFees.Equal(w.RemainingForeignServiceFees) &&
		v.TotalCollateral == w.TotalCollateral &&
		v.TotalShield == w.TotalShield &&
		v.TotalWithdrawing == w.TotalWithdrawing &&
		v.Height == w.Height
}

// NewShieldStatusRow allows to build a new ShieldStatusRow
func NewShieldStatusRow(
	gobalStakingPool int64,
	currentNativeServiceFees DbDecCoins,
	currentForeignServiceFees DbDecCoins,
	remainingNativeServiceFees DbDecCoins,
	remainingForeignServiceFees DbDecCoins,
	totalCollateral int64,
	totalShield int64,
	totalWithdrawing int64,
	height int64,
) ShieldStatusRow {
	return ShieldStatusRow{
		OneRowID:                    true,
		GobalStakingPool:            gobalStakingPool,
		CurrentNativeServiceFees:    &currentNativeServiceFees,
		CurrentForeignServiceFees:   &currentForeignServiceFees,
		RemainingNativeServiceFees:  &remainingNativeServiceFees,
		RemainingForeignServiceFees: &remainingForeignServiceFees,
		TotalCollateral:             totalCollateral,
		TotalShield:                 totalShield,
		TotalWithdrawing:            totalWithdrawing,
		Height:                      height,
	}
}

// ----------------------------------------------------------------

type ShieldPoolParamsRow struct {
	OneRowID bool   `db:"one_row_id"`
	Params   string `db:"params"`
	Height   int64  `db:"height"`
}

// Equal reports whether m and n represent the same table rows.
func (m ShieldPoolParamsRow) Equal(n ShieldPoolParamsRow) bool {
	return m.Params == n.Params &&
		m.Height == n.Height
}

// NewShieldPoolParamsRow allows to build a new ShieldPoolParamsRow
func NewShieldPoolParamsRow(params string, height int64) ShieldPoolParamsRow {
	return ShieldPoolParamsRow{
		OneRowID: true,
		Params:   params,
		Height:   height,
	}
}
