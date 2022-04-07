package banking

import (
	bankingtypes "git.ooo.ua/vipcoin/chain/x/banking/types"
	"git.ooo.ua/vipcoin/lib/filter"

	"github.com/forbole/bdjuno/v2/database/types"
)

// SaveMsgSetTransferExtra - method that create extras to the "vipcoin_chain_banking_set_transfer_extra" table
func (r Repository) SaveMsgSetTransferExtra(extras ...*bankingtypes.MsgSetTransferExtra) error {
	if len(extras) == 0 {
		return nil
	}

	query := `INSERT INTO vipcoin_chain_banking_set_transfer_extra 
		(creator, id, extras) 
		VALUES 
		(:creator, :id, :extras)`

	if _, err := r.db.NamedExec(query, toMsgSetTransferExtrasDatabase(extras...)); err != nil {
		return err
	}

	return nil
}

// GetMsgSetTransferExtra - method that get extras from the "vipcoin_chain_banking_set_transfer_extra" table
func (r Repository) GetMsgSetTransferExtra(filter filter.Filter) ([]*bankingtypes.MsgSetTransferExtra, error) {
	query, args := filter.Build(
		tableMsgSetTransferExtra,
		types.FieldCreator, types.FieldID, types.FieldExtras,
	)

	var result []types.DBSetTransferExtra
	if err := r.db.Select(&result, query, args...); err != nil {
		return []*bankingtypes.MsgSetTransferExtra{}, err
	}

	extras := make([]*bankingtypes.MsgSetTransferExtra, 0, len(result))
	for _, extra := range result {
		extras = append(extras, toMsgSetTransferExtraDomain(extra))
	}

	return extras, nil
}
