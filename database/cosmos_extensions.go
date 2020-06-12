package database

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
)

// SaveValidators allows the bulk saving of a list of validators
func (db BigDipperDb) SaveValidators(validators []staking.Validator) error {
	var insertParams []interface{}

	queryInsert := "INSERT INTO validator (consensus_address, consensus_pubkey) VALUES "
	for i, result := range validators {
		p1 := i * 2 // starting position for insert params

		queryInsert += fmt.Sprintf("($%d,$%d),", p1+1, p1+2)

		key, err := sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeConsPub, result.ConsPubKey)
		if err != nil {
			return err
		}

		insertParams = append(insertParams, result.ConsAddress().String(), key)
	}

	queryInsert = queryInsert[:len(queryInsert)-1] // remove trailing ","
	queryInsert += " ON CONFLICT DO NOTHING"
	_, err := db.Sql.Exec(queryInsert, insertParams...)
	return err
}
