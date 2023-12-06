package allowed

import (
	"errors"
	"strings"

	"git.ooo.ua/vipcoin/lib/errs"
	"git.ooo.ua/vipcoin/lib/filter"
	allowed "git.ooo.ua/vipcoin/ovg-chain/x/allowed/types"
	juno "github.com/forbole/juno/v5/types"

	"github.com/forbole/bdjuno/v4/database/types"
)

// handleMsgDeleteByAddresses allows to properly handle a MsgDeleteByAddresses
func (m *Module) handleMsgDeleteByAddresses(tx *juno.Tx, index int, msg *allowed.MsgDeleteByAddresses) error {
	// 1) logic for table overgold_allowed_delete_by_addresses
	// 1.1) check if already exists (not found is ok)
	deleteAddresses, err := m.allowedRepo.GetAllDeleteByAddresses(filter.NewFilter().SetCondition(filter.ConditionAND).
		SetArgument(types.FieldCreator, msg.Creator).
		SetArgument(types.FieldAddress, msg.Address))
	if err != nil {
		if !errors.As(err, &errs.NotFound{}) {
			return err
		}
	}
	if len(deleteAddresses) > 0 {
		return errs.AlreadyExists{What: "delete_by_addresses, address: " + strings.Join(msg.Address, ", ")}
	}

	// 1.2) insert to table
	if err = m.allowedRepo.InsertToDeleteByAddresses(tx.TxHash, msg); err != nil {
		return err
	}

	// 2) logic for table overgold_allowed_addresses
	// 2.1) check if already exists
	_, err = m.allowedRepo.GetAllAddresses(filter.NewFilter().SetCondition(filter.ConditionAND).
		SetArgument(types.FieldCreator, msg.Creator).
		SetArgument(types.FieldAddress, msg.Address))
	if err != nil {
		if !errors.As(err, &errs.NotFound{}) {
			return err
		}
	}

	// 2.2) delete data from table
	return m.allowedRepo.DeleteAddressesByAddress(msg.Address...)
}
