package database

import (
	"fmt"

	dbtypes "github.com/forbole/bdjuno/v2/database/types"
	dbutils "github.com/forbole/bdjuno/v2/database/utils"
	markettypes "github.com/ovrclk/akash/x/market/types/v1beta2"
)

func (db *Db) SaveGenesisLeases(leases []markettypes.Lease, height int64) error {

	return nil
}

func (db *Db) SaveLeases(responses []markettypes.QueryLeaseResponse, height int64) error {
	if len(responses) == 0 {
		return nil
	}

	paramsNumber := 17
	slices := dbutils.SplitLeases(responses, paramsNumber)
	for _, s := range slices {
		err := db.saveLeases(s, paramsNumber, height)
		if err != nil {
			return fmt.Errorf("error while saving leases: %s", err)
		}
	}

	return nil
}

func (db *Db) saveLeases(slices []markettypes.QueryLeaseResponse, paramsNumber int, height int64) error {
	stmt := `INSERT INTO lease ( 
		owner_address, dseq, gseq, oseq, provider_address, 
		lease_state, price, created_at, closed_on, 
		account_id, payment_id, owner_address, payment_state, rate, balance, withdrawn, height 
		) VALUES `

	var params []interface{}
	for i, s := range slices {

		ai := i * paramsNumber

		stmt += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d),",
			ai+1, ai+2, ai+3, ai+4, ai+5, ai+6, ai+7, ai+8, ai+9,
			ai+10, ai+11, ai+12, ai+13, ai+14, ai+15, ai+16, ai+17)

		leaseID := s.Lease.LeaseID
		lease := s.Lease
		escrow := s.EscrowPayment

		params = append(params,
			// Lease ID
			leaseID.Owner, leaseID.DSeq, leaseID.GSeq, leaseID.OSeq, leaseID.Provider,

			// Lease
			lease.State, dbtypes.NewDbDecCoin(lease.Price), lease.CreatedAt, lease.ClosedOn,

			// Escrow Payment
			dbtypes.NewDbLeaseAccountID(escrow.AccountID), escrow.PaymentID, escrow.State, dbtypes.NewDbDecCoin(escrow.Rate),
			dbtypes.NewDbDecCoin(escrow.Balance), dbtypes.NewDbCoin(escrow.Withdrawn),

			height,
		)
	}

	stmt = stmt[:len(stmt)-1] // Remove trailing ","
	stmt += `
ON CONFLICT (owner_address, dseq, gseq, oseq, provider_address) DO UPDATE SET 
	lease_state = excluded.lease_state,
	price = excluded.price,
	created_at = excluded.created_at,
	closed_on = excluded.closed_on,
	account_id = excluded.account_id,
	payment_id = excluded.payment_id,
	owner_address = excluded.owner_address,
	payment_state = excluded.payment_state,
	rate = excluded.rate,
	balance = excluded.balance,
	withdrawn = excluded.withdrawn,
	height = excluded.height
WHERE lease.height <= excluded.height`

	_, err := db.Sql.Exec(stmt, params...)
	if err != nil {
		return fmt.Errorf("error while storing leases: %s", err)
	}

	return nil
}
