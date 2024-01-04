package feeexcluder

import (
	"errors"
	"strconv"

	"git.ooo.ua/vipcoin/lib/errs"
	"git.ooo.ua/vipcoin/lib/filter"
	"git.ooo.ua/vipcoin/ovg-chain/x/feeexcluder/types"
	juno "github.com/forbole/juno/v5/types"

	db "github.com/forbole/bdjuno/v4/database/types"
)

// handleMsgUpdateAddress allows to properly handle a message
func (m *Module) handleMsgUpdateAddress(tx *juno.Tx, _ int, msg *types.MsgUpdateAddress) error {
	// 1) logic for table overgold_feeexcluder_update_address
	// 1.1) check if already exists (not found is ok)
	updAddresses, err := m.feeexcluderRepo.GetAllMsgUpdateAddress(filter.NewFilter().SetCondition(filter.ConditionAND).
		SetArgument(db.FieldID, msg.Id).
		SetArgument(db.FieldAddress, msg.Address))
	if err != nil {
		if !errors.As(err, &errs.NotFound{}) {
			return err
		}
	}
	if len(updAddresses) > 0 {
		return errs.AlreadyExists{What: "update_addresses, id: " + strconv.FormatUint(msg.Id, 10)}
	}

	// 1.2) insert to table
	if err = m.feeexcluderRepo.InsertToMsgUpdateAddress(tx.TxHash, *msg); err != nil {
		return err
	}

	// 2) logic for table overgold_feeexcluder_address
	// 2.1) check if already exists
	addressList, err := m.feeexcluderRepo.GetAllAddress(filter.NewFilter().SetCondition(filter.ConditionAND).
		SetArgument(db.FieldCreator, msg.Creator).
		SetArgument(db.FieldAddress, msg.Address).
		SetArgument(db.FieldMsgID, msg.Id))
	if err != nil {
		if !errors.As(err, &errs.NotFound{}) {
			return err
		}
	}
	if len(addressList) != 1 {
		return errs.Internal{Cause: "expected only 1 address"}
	}

	// 2.2) update data in table
	return m.feeexcluderRepo.UpdateAddress(nil, addressList[0].Id, types.Address{
		Id:      msg.Id,
		Address: msg.Address,
		Creator: msg.Creator,
	})
}
