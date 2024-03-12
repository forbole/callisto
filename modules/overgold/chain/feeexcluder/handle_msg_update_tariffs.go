package feeexcluder

import (
	"errors"

	"git.ooo.ua/vipcoin/lib/errs"
	"git.ooo.ua/vipcoin/lib/filter"
	"git.ooo.ua/vipcoin/ovg-chain/x/feeexcluder/types"
	juno "github.com/forbole/juno/v5/types"

	db "github.com/forbole/bdjuno/v4/database/types"
)

// handleMsgUpdateTariffs allows to properly handle a message
func (m *Module) handleMsgUpdateTariffs(tx *juno.Tx, _ int, msg *types.MsgUpdateTariffs) error {
	// 1.2) insert to table
	if err := m.feeexcluderRepo.InsertToMsgUpdateTariffs(tx.TxHash, *msg); err != nil {
		return err
	}

	// 2) logic for table overgold_feeexcluder_tariffs
	// 2.1) check if already exists
	tariffsList, err := m.feeexcluderRepo.GetAllTariffs(filter.NewFilter().SetArgument(db.FieldMsgID, msg.Tariff.Id))
	if err != nil {
		if !errors.As(err, &errs.NotFound{}) {
			return err
		}
	}
	if len(tariffsList) != 1 {
		return errs.Internal{Cause: "expected only 1 tariffs"}
	}

	// 2.2) get unique id for tariffs
	tariffsDB, err := m.feeexcluderRepo.GetTariffsDB(filter.NewFilter().SetArgument(db.FieldMsgID, msg.Tariff.Id))
	if err != nil {
		return err
	}

	// 2.3) update data in table
	return m.feeexcluderRepo.UpdateTariffs(nil, tariffsDB.ID, types.Tariffs{
		Denom:   msg.Denom,
		Creator: msg.Creator,
		Tariffs: []*types.Tariff{msg.Tariff},
	})
}
