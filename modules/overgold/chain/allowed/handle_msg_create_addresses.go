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

// handleMsgCreateAddresses allows to properly handle a MsgCreateAddresses
func (m *Module) handleMsgCreateAddresses(tx *juno.Tx, _ int, msg *allowed.MsgCreateAddresses) error {
	// 1) logic for table overgold_allowed_create_addresses
	// 1.1) check if already exists (not found is ok)
	createAddresses, err := m.allowedRepo.GetAllCreateAddresses(filter.NewFilter().SetCondition(filter.ConditionAND).
		SetArgument(types.FieldCreator, msg.Creator).
		SetArgument(types.FieldAddress, msg.Address))
	if err != nil {
		if !errors.As(err, &errs.NotFound{}) {
			return err
		}
	}
	if len(createAddresses) > 0 {
		return errs.AlreadyExists{What: "create_addresses, address: " + strings.Join(msg.Address, ", ")}
	}

	// 1.2) insert to table
	if err = m.allowedRepo.InsertToCreateAddresses(tx.TxHash, msg); err != nil {
		return err
	}

	// 2) logic for table overgold_allowed_addresses
	// 2.1) check if already exists (not found is ok)
	addresses, err := m.allowedRepo.GetAllAddresses(filter.NewFilter().SetCondition(filter.ConditionAND).
		SetArgument(types.FieldCreator, msg.Creator).
		SetArgument(types.FieldAddress, msg.Address))
	if err != nil {
		if !errors.As(err, &errs.NotFound{}) {
			return err
		}
	}
	if len(addresses) > 0 {
		return errs.AlreadyExists{What: "addresses, address: " + strings.Join(msg.Address, ", ")}
	}

	// 2.2) insert to table
	return m.allowedRepo.InsertToAddresses(allowed.Addresses{
		Address: msg.Address,
		Creator: msg.Creator,
	})
}
