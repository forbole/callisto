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

// handleMsgUpdateAddresses allows to properly handle a MsgUpdateAddresses
func (m *Module) handleMsgUpdateAddresses(tx *juno.Tx, index int, msg *allowed.MsgUpdateAddresses) error {
	// 1) logic for table overgold_allowed_update_addresses
	// 1.1) check if already exists (not found is ok)
	updateAddresses, err := m.allowedRepo.GetAllUpdateAddresses(filter.NewFilter().SetCondition(filter.ConditionAND).
		SetArgument(types.FieldID, msg.Id).
		SetArgument(types.FieldCreator, msg.Creator).
		SetArgument(types.FieldAddress, msg.Address))
	if err != nil {
		if !errors.As(err, &errs.NotFound{}) {
			return err
		}
	}
	if len(updateAddresses) > 0 {
		return errs.AlreadyExists{What: "update_addresses, id: " + strconv.FormatUint(msg.Id, 10)}
	}

	// 1.2) insert to table
	if err := m.allowedRepo.InsertToUpdateAddresses(tx.TxHash, msg); err != nil {
		return err
	}

	// 2) logic for table overgold_allowed_addresses
	// 2.1) check if already exists
	_, err = m.allowedRepo.GetAllAddresses(filter.NewFilter().SetCondition(filter.ConditionAND).
		SetArgument(types.FieldCreator, msg.Creator).
		SetArgument(types.FieldAddress, msg.Address).
		SetArgument(types.FieldID, msg.Id))
	if err != nil {
		if !errors.As(err, &errs.NotFound{}) {
			return err
		}
	}

	// 2.2) update data in table
	return m.allowedRepo.UpdateAddresses(allowed.Addresses{
		Id:      msg.Id,
		Address: msg.Address,
		Creator: msg.Creator,
	})
}
