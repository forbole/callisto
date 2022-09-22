package database

import (
	"encoding/json"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/forbole/bdjuno/v3/types"
	"github.com/lib/pq"
	abci "github.com/tendermint/tendermint/abci/types"
	// abci "github.com/tendermint/tendermint/abci/types"
)

// SaveMarginParams allows to store the given params inside the database
func (db *Db) SaveMarginParams(params *types.MarginParams) error {
	paramsBz, err := json.Marshal(&params.Params)
	if err != nil {
		return fmt.Errorf("error while marshaling margin params: %s", err)
	}

	stmt := `
INSERT INTO margin_params (params, height) 
VALUES ($1, $2)
ON CONFLICT (one_row_id) DO UPDATE 
    SET params = excluded.params,
        height = excluded.height
WHERE margin_params.height <= excluded.height`

	_, err = db.Sql.Exec(stmt, string(paramsBz), params.Height)
	if err != nil {
		return fmt.Errorf("error while storing margin params: %s", err)
	}

	return nil
}

// SaveMarginEvent allows to store the given x/margin events inside the database
func (db *Db) SaveMarginEvent(events []types.MarginEvent) error {
	stmt := `INSERT INTO margin_events (transaction_hash, index, type, value, involved_accounts_addresses, height) VALUES`

	if len(events) == 0 {
		return nil
	}

	var marginEvents []interface{}

	for i, event := range events {
		vi := i * 6
		stmt += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d),", vi+1, vi+2, vi+3, vi+4, vi+5, vi+6)
		var eventAttr []abci.EventAttribute
		for _, attr := range event.Value.Attributes {
			value := []byte(strings.Replace(string(attr.Value), " ", ", ", -1))
			eventAttr = append(eventAttr, abci.EventAttribute{Key: attr.Key, Value: value})
		}
		eventObj := abci.Event{Type: event.Value.Type, Attributes: eventAttr}

		ev := sdk.StringifyEvent(eventObj)
		eventBz, err := json.Marshal(&ev.Attributes)
		if err != nil {
			return fmt.Errorf("error while marshaling x/margin events: %s", err)
		}

		marginEvents = append(marginEvents, event.TxHash, event.Index, event.MsgType, string(eventBz), pq.Array(event.Addressess), event.Height)
	}

	stmt = stmt[:len(stmt)-1] // Remove trailing ","
	stmt += `
ON CONFLICT DO NOTHING`

	_, err := db.Sql.Exec(stmt, marginEvents...)
	if err != nil {
		return fmt.Errorf("error while storing x/margin events: %s", err)
	}

	return nil
}
