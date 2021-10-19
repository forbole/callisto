package types

import (
	"fmt"
	"time"
)

// AccountRow represents a single row inside the account table
type AccountRow struct {
	Address string `db:"address"`
}

// NewAccountRow allows to easily build a new AccountRow
func NewAccountRow(address string) AccountRow {
	return AccountRow{
		Address: address,
	}
}

// Equal tells whether a and b contain the same data
func (a AccountRow) Equal(b AccountRow) bool {
	return a.Address == b.Address
}

// --------------- For Vesting Accounts ---------------

// ContinuousVestingAccountRow represents a single row inside the vesting_account table
type ContinuousVestingAccountRow struct {
	ID              int       `db:"id"`
	Type            string    `db:"type"`
	Address         string    `db:"address"`
	OriginalVesting *DbCoins  `db:"original_vesting"`
	EndTime         time.Time `db:"end_time"`
	StartTime       time.Time `db:"start_time"`
}

// NewContinuousVestingAccountRow allows to build a new DB ContinuousVestingAccountRow
func NewContinuousVestingAccountRow(
	id int,
	accountType string,
	address string,
	originalVesting DbCoins,
	endTime time.Time,
	startTime time.Time,
) ContinuousVestingAccountRow {
	return ContinuousVestingAccountRow{
		ID:              id,
		Type:            accountType,
		Address:         address,
		OriginalVesting: &originalVesting,
		EndTime:         endTime,
		StartTime:       startTime,
	}
}

// Equal tells whether a and b contain the same data
func (a ContinuousVestingAccountRow) Equal(b ContinuousVestingAccountRow) bool {
	return a.ID == b.ID &&
		a.Type == b.Type &&
		a.Address == b.Address &&
		a.OriginalVesting.Equal(b.OriginalVesting)
}

// DelayedVestingAccountRow represents a single row inside the vesting_account table
type DelayedVestingAccountRow struct {
	ID              int       `db:"id"`
	Type            string    `db:"type"`
	Address         string    `db:"address"`
	OriginalVesting *DbCoins  `db:"original_vesting"`
	EndTime         time.Time `db:"end_time"`
}

// NewDelayedVestingAccountRow allows to build a new DB DelayedVestingAccountRow
func NewDelayedVestingAccountRow(
	id int,
	accountType string,
	address string,
	originalVesting DbCoins,
	endTime time.Time,
) DelayedVestingAccountRow {
	return DelayedVestingAccountRow{
		ID:              id,
		Type:            accountType,
		Address:         address,
		OriginalVesting: &originalVesting,
		EndTime:         endTime,
	}
}

// Equal tells whether a and b contain the same data
func (a DelayedVestingAccountRow) Equal(b DelayedVestingAccountRow) bool {
	return a.ID == b.ID &&
		a.Type == b.Type &&
		a.Address == b.Address &&
		a.OriginalVesting.Equal(b.OriginalVesting)
}

// PeriodicVestingAccountRow represents a single row inside the vesting_account table
type PeriodicVestingAccountRow struct {
	ID              int       `db:"id"`
	Type            string    `db:"type"`
	Address         string    `db:"address"`
	OriginalVesting *DbCoins  `db:"original_vesting"`
	EndTime         time.Time `db:"end_time"`
	StartTime       time.Time `db:"start_time"`
}

// NewPeriodicVestingAccountRow allows to build a new DB PeriodicVestingAccountRow
func NewPeriodicVestingAccountRow(
	id int,
	accountType string,
	address string,
	originalVesting DbCoins,
	endTime time.Time,
	startTime time.Time,
) PeriodicVestingAccountRow {
	return PeriodicVestingAccountRow{
		ID:              id,
		Type:            accountType,
		Address:         address,
		OriginalVesting: &originalVesting,
		EndTime:         endTime,
		StartTime:       startTime,
	}
}

// Equal tells whether a and b contain the same data
func (a PeriodicVestingAccountRow) Equal(b PeriodicVestingAccountRow) bool {
	return a.ID == b.ID &&
		a.Type == b.Type &&
		a.Address == b.Address &&
		a.OriginalVesting.Equal(b.OriginalVesting)
}

// VestingPeriodRow represents a Periodic Vesting Account
type VestingPeriodRow struct {
	VestingAccountID int      `db:"vesting_account_id"`
	PeriodOrder      int      `db:"period_order"`
	Length           string   `db:"length"`
	Amount           *DbCoins `db:"amount"`
}

// NewPeriodicVestingAccountRow allows to build a new DB PeriodicVestingAccountRow
func NewVestingPeriodRow(
	vestingAccountID int,
	periodOrder int,
	length string,
	amount DbCoins,
) VestingPeriodRow {
	return VestingPeriodRow{
		VestingAccountID: vestingAccountID,
		PeriodOrder:      periodOrder,
		Length:           length,
		Amount:           &amount,
	}
}

// Equal tells whether a and b contain the same data
func (a VestingPeriodRow) Equal(b VestingPeriodRow) bool {
	fmt.Println(a, b)
	return a.VestingAccountID == b.VestingAccountID &&
		a.PeriodOrder == b.PeriodOrder &&
		a.Length == b.Length &&
		a.Amount.Equal(b.Amount)
}
