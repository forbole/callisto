package types

import (
	"database/sql"
	"time"

	"github.com/shopspring/decimal"
)

// Database types
type (
	// FeeExcluderAddress represents a single row inside the overgold_feeexcluder_address
	FeeExcluderAddress struct {
		ID      uint64        `db:"id"`     // unique id as primary key
		MsgID   sql.NullInt64 `db:"msg_id"` // message id
		Address string        `db:"address"`
		Creator string        `db:"creator"`
	}

	// FeeExcluderDailyStats represents a single row inside the overgold_feeexcluder_daily_stats
	FeeExcluderDailyStats struct {
		CountWithFee  int32   `db:"count_with_fee"`
		CountNoFee    int32   `db:"count_no_fee"`
		ID            uint64  `db:"id"`     // unique id as primary key
		MsgID         uint64  `db:"msg_id"` // daily stats id from message
		AmountWithFee DbCoins `db:"amount_with_fee"`
		AmountNoFee   DbCoins `db:"amount_no_fee"`
		Fee           DbCoins `db:"fee"`
	}

	// FeeExcluderStats represents a single row inside the overgold_feeexcluder_stats
	FeeExcluderStats struct {
		ID           string    `db:"id"`
		DailyStatsID uint64    `db:"daily_stats_id"`
		Date         time.Time `db:"date"`
	}

	// FeeExcluderFees represents a single row inside the overgold_feeexcluder_fees
	FeeExcluderFees struct {
		NoRefReward bool            `db:"no_ref_reward"`
		ID          uint64          `db:"id"`     // unique id as primary key
		MsgID       uint64          `db:"msg_id"` // fees id from message
		MinAmount   uint64          `db:"min_amount"`
		AmountFrom  uint64          `db:"amount_from"`
		Fee         decimal.Decimal `db:"fee"`
		RefReward   decimal.Decimal `db:"ref_reward"`
		StakeReward decimal.Decimal `db:"stake_reward"`
		Creator     string          `db:"creator"`
	}

	// FeeExcluderTariff represents a single row inside the overgold_feeexcluder_tariff
	FeeExcluderTariff struct {
		ID            uint64 `db:"id"`     // unique id as primary key
		MsgID         uint64 `db:"msg_id"` // tariff id from message
		Amount        uint64 `db:"amount"`
		MinRefBalance uint64 `db:"min_ref_balance"`
		Denom         string `db:"denom"`
	}

	// FeeExcluderM2MTariffFees represents a single row inside the overgold_feeexcluder_m2m_tariff_fees
	FeeExcluderM2MTariffFees struct {
		TariffID uint64 `db:"tariff_id"`
		FeesID   uint64 `db:"fees_id"`
	}

	// FeeExcluderTariffs represents a single row inside the overgold_feeexcluder_tariffs
	FeeExcluderTariffs struct {
		ID      uint64 `db:"id"` // unique id as primary key
		Denom   string `db:"denom"`
		Creator string `db:"creator"`
	}

	// FeeExcluderM2MTariffTariffs represents a single row inside the overgold_feeexcluder_m2m_tariff_tariffs
	FeeExcluderM2MTariffTariffs struct {
		TariffID  uint64 `db:"tariff_id"`
		TariffsID uint64 `db:"tariffs_id"`
	}

	// FeeExcluderGenesisState represents a single row inside the overgold_feeexcluder_genesis_state
	FeeExcluderGenesisState struct {
		ID              uint64 `db:"id"`
		AddressCount    uint64 `db:"address_count"`
		DailyStatsCount uint64 `db:"daily_stats_count"`
	}

	// FeeExcluderM2MGenesisStateAddress represents a single row inside the overgold_feeexcluder_m2m_genesis_state_address
	FeeExcluderM2MGenesisStateAddress struct {
		GenesisStateID uint64 `db:"genesis_state_id"`
		AddressID      uint64 `db:"address_id"`
	}

	// FeeExcluderM2MGenesisStateDailyStats represents a single row inside the overgold_feeexcluder_m2m_genesis_state_daily_stats
	FeeExcluderM2MGenesisStateDailyStats struct {
		GenesisStateID uint64 `db:"genesis_state_id"`
		DailyStatsID   uint64 `db:"daily_stats_id"`
	}

	// FeeExcluderM2MGenesisStateStats represents a single row inside the overgold_feeexcluder_m2m_genesis_state_stats
	FeeExcluderM2MGenesisStateStats struct {
		GenesisStateID uint64 `db:"genesis_state_id"`
		StatsID        string `db:"stats_id"`
	}

	// FeeExcluderM2MGenesisStateTariffs represents a single row inside the overgold_feeexcluder_m2m_genesis_state_tariffs
	FeeExcluderM2MGenesisStateTariffs struct {
		GenesisStateID uint64 `db:"genesis_state_id"`
		TariffsID      uint64 `db:"tariffs_id"`
	}

	// FeeExcluderCreateAddress represents a single row inside the overgold_feeexcluder_create_address
	FeeExcluderCreateAddress struct {
		ID      uint64 `db:"id"`
		TxHash  string `db:"tx_hash"`
		Creator string `db:"creator"`
		Address string `db:"address"`
	}

	// FeeExcluderUpdateAddress represents a single row inside the overgold_feeexcluder_update_address
	FeeExcluderUpdateAddress struct {
		ID      uint64 `db:"id"`
		TxHash  string `db:"tx_hash"`
		Creator string `db:"creator"`
		Address string `db:"address"`
	}

	// FeeExcluderDeleteAddress represents a single row inside the overgold_feeexcluder_delete_address
	FeeExcluderDeleteAddress struct {
		ID      uint64 `db:"id"`
		TxHash  string `db:"tx_hash"`
		Creator string `db:"creator"`
	}

	// FeeExcluderCreateTariffs represents a single row inside the overgold_feeexcluder_create_tariffs
	FeeExcluderCreateTariffs struct {
		ID       uint64 `db:"id"`
		TariffID uint64 `db:"tariff_id"`
		TxHash   string `db:"tx_hash"`
		Creator  string `db:"creator"`
		Denom    string `db:"denom"`
	}

	// FeeExcluderUpdateTariffs represents a single row inside the overgold_feeexcluder_update_tariffs
	FeeExcluderUpdateTariffs struct {
		ID       uint64 `db:"id"`
		TariffID uint64 `db:"tariff_id"`
		TxHash   string `db:"tx_hash"`
		Creator  string `db:"creator"`
		Denom    string `db:"denom"`
	}

	// FeeExcluderDeleteTariffs represents a single row inside the overgold_feeexcluder_delete_tariffs
	FeeExcluderDeleteTariffs struct {
		ID       uint64 `db:"id"`
		TariffID uint64 `db:"tariff_id"`
		FeesID   uint64 `db:"fees_id"`
		TxHash   string `db:"tx_hash"`
		Creator  string `db:"creator"`
		Denom    string `db:"denom"`
	}
)

// GetAddressIDs - returns a list of ids. TODO: use as a method for new type []FeeExcluderM2MGenesisStateAddress
func GetAddressIDs(list []FeeExcluderM2MGenesisStateAddress) []uint64 {
	ids := make([]uint64, 0, len(list))
	for _, stats := range list {
		ids = append(ids, stats.AddressID)
	}
	return ids
}

// GetStatsIDs - returns a list of ids. TODO: use as a method for new type []FeeExcluderM2MGenesisStateStats
func GetStatsIDs(list []FeeExcluderM2MGenesisStateStats) []string {
	ids := make([]string, 0, len(list))
	for _, stats := range list {
		ids = append(ids, stats.StatsID)
	}
	return ids
}

// GetDailyStatsIDs - returns a list of ids. TODO: use as a method for new type []FeeExcluderM2MGenesisStateDailyStats
func GetDailyStatsIDs(list []FeeExcluderM2MGenesisStateDailyStats) []uint64 {
	ids := make([]uint64, 0, len(list))
	for _, stats := range list {
		ids = append(ids, stats.DailyStatsID)
	}
	return ids
}

// GetTariffsIDs - returns a list of ids. TODO: use as a method for new type []FeeExcluderM2MGenesisStateTariffs
func GetTariffsIDs(list []FeeExcluderM2MGenesisStateTariffs) []uint64 {
	ids := make([]uint64, 0, len(list))
	for _, stats := range list {
		ids = append(ids, stats.TariffsID)
	}
	return ids
}
