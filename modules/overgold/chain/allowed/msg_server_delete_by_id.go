package allowed

import (
	"errors"
	"strconv"

	"git.ooo.ua/vipcoin/lib/errs"
	"git.ooo.ua/vipcoin/lib/filter"
	allowed "git.ooo.ua/vipcoin/ovg-chain/x/allowed/types"
	juno "github.com/forbole/juno/v5/types"

	"github.com/forbole/bdjuno/v4/database/types"
)

// handleMsgDeleteByID allows to properly handle a MsgDeleteByID
func (m *Module) handleMsgDeleteByID(tx *juno.Tx, index int, msg *allowed.MsgDeleteByID) error {
	// 1) logic for table overgold_allowed_delete_by_id
	// 1.1) check if already exists (not found is ok)
	deleteByIDs, err := m.allowedRepo.GetAllDeleteByID(filter.NewFilter().SetCondition(filter.ConditionAND).
		SetArgument(types.FieldCreator, msg.Creator).
		SetArgument(types.FieldID, msg.Id))
	if err != nil && !errors.Is(err, errs.NotFound{}) {
		return err
	}
	if len(deleteByIDs) > 0 {
		return errs.AlreadyExists{What: "delete_by_id, id: " + strconv.FormatUint(msg.Id, 10)}
	}

	// 1.2) insert to table
	if err = m.allowedRepo.InsertToDeleteByID(tx.TxHash, msg); err != nil {
		return err
	}

	// 2) logic for table overgold_allowed_addresses
	// 2.1) check if already exists
	_, err = m.allowedRepo.GetAllAddresses(filter.NewFilter().SetCondition(filter.ConditionAND).
		SetArgument(types.FieldCreator, msg.Creator).
		SetArgument(types.FieldID, msg.Id))
	if err != nil {
		return err
	}

	// 2.2) delete data from table
	return m.allowedRepo.DeleteAddressesByID(msg.Id)
}
