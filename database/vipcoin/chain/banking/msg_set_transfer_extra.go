package banking

import (
	bankingtypes "git.ooo.ua/vipcoin/chain/x/banking/types"
	"git.ooo.ua/vipcoin/lib/errs"
	"git.ooo.ua/vipcoin/lib/filter"

	"github.com/forbole/bdjuno/v2/database/types"
)

// SaveMsgSetTransferExtra - method that create extras to the "vipcoin_chain_banking_set_transfer_extra" table
func (r Repository) SaveMsgSetTransferExtra(extras *bankingtypes.MsgSetTransferExtra, transactionHash string) error {
	query := `INSERT INTO vipcoin_chain_banking_set_transfer_extra 
		(transaction_hash, creator, id, extras) 
		VALUES 
		(:transaction_hash, :creator, :id, :extras)`

	if _, err := r.db.NamedExec(query, toMsgSetTransferExtraDatabase(extras, transactionHash)); err != nil {
		return errs.Internal{Cause: err.Error()}
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
		return []*bankingtypes.MsgSetTransferExtra{}, errs.Internal{Cause: err.Error()}
	}

	extras := make([]*bankingtypes.MsgSetTransferExtra, 0, len(result))
	for _, extra := range result {
		extras = append(extras, toMsgSetTransferExtraDomain(extra))
	}

	return extras, nil
}
