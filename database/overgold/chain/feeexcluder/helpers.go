package feeexcluder

import (
	"git.ooo.ua/vipcoin/lib/errs"
	"github.com/jmoiron/sqlx"
)

const (
	// genesis state
	tableGenesisState = "overgold_feeexcluder_genesis_state"

	// address
	tableAddress       = "overgold_feeexcluder_address"
	tableCreateAddress = "overgold_feeexcluder_create_address"
	tableDeleteAddress = "overgold_feeexcluder_delete_address"
	tableUpdateAddress = "overgold_feeexcluder_update_address"

	// tariffs
	tableTariffs       = "overgold_feeexcluder_tariffs"
	tableCreateTariffs = "overgold_feeexcluder_create_tariffs"
	tableDeleteTariffs = "overgold_feeexcluder_delete_tariffs"
	tableUpdateTariffs = "overgold_feeexcluder_update_tariffs"

	// many-to-many
	tableM2MGenesisStateAddress    = "overgold_feeexcluder_m2m_genesis_state_address"
	tableM2MGenesisStateDailyStats = "overgold_feeexcluder_m2m_genesis_state_daily_stats"
	tableM2MGenesisStateStats      = "overgold_feeexcluder_m2m_genesis_state_stats"
	tableM2MGenesisStateTariffs    = "overgold_feeexcluder_m2m_genesis_state_tariffs"
	tableM2MTariffFees             = "overgold_feeexcluder_m2m_tariff_fees"
	tableM2MTariffTariffs          = "overgold_feeexcluder_m2m_tariff_tariffs"

	// helpers
	tableDailyStats = "overgold_feeexcluder_daily_stats"
	tableFees       = "overgold_feeexcluder_fees"
	tableStats      = "overgold_feeexcluder_stats"
	tableTariff     = "overgold_feeexcluder_tariff"
)

const (
	layoutDate = "2006-01-02"
)

func commit(tx *sqlx.Tx, err error) {
	if err != nil {
		_ = tx.Rollback()
	} else if err = tx.Commit(); err != nil {
		err = errs.Internal{Cause: err.Error()}
	}
}
