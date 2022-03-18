package database

import (
	"fmt"

	dbtypes "github.com/forbole/bdjuno/v2/database/types"
	"github.com/lib/pq"
	escrowtypes "github.com/ovrclk/akash/x/escrow/types/v1beta2"
	markettypes "github.com/ovrclk/akash/x/market/types/v1beta2"
)

func (db *Db) SaveGenesisLeases(leases []markettypes.Lease, height int64) error {
	for _, lease := range leases {
		leaseID, err := db.saveLeaseID(lease)
		if err != nil {
			return fmt.Errorf("error while storing lease ID: %s", err)
		}

		err = db.saveLease(leaseID, lease, height)
		if err != nil {
			return fmt.Errorf("error while storing lease: %s", err)
		}
	}

	return nil
}

func (db *Db) SaveLeases(responses []markettypes.QueryLeaseResponse, height int64) error {
	for _, res := range responses {
		leaseID, err := db.saveLeaseID(res.Lease)
		if err != nil {
			return fmt.Errorf("error while storing lease ID: %s", err)
		}

		err = db.saveLease(leaseID, res.Lease, height)
		if err != nil {
			return fmt.Errorf("error while storing lease: %s", err)
		}

		err = db.saveEscrowPayment(leaseID, res.EscrowPayment, height)
		if err != nil {
			return fmt.Errorf("error while storing escrow payment: %s", err)
		}
	}

	return nil
}

func (db *Db) SaveMarketParams(p markettypes.Params, height int64) error {
	stmt := `
	INSERT INTO market_params (bid_min_deposit, order_max_bids, height) 
	VALUES ($1, $2, $3) 
	ON CONFLICT (one_row_id) DO UPDATE 
	SET bid_min_deposit = excluded.bid_min_deposit, 
	order_max_bids = excluded.order_max_bids, 
		height = excluded.height 
	WHERE market_params.height <= excluded.height`

	_, err := db.Sql.Exec(stmt,
		pq.Array(dbtypes.NewDbCoin(p.BidMinDeposit)),
		p.OrderMaxBids,
		height,
	)
	if err != nil {
		return fmt.Errorf("error while storing market params: %s", err)
	}
	return nil
}

func (db *Db) saveLeaseID(lease markettypes.Lease) (int64, error) {
	stmt := `
	INSERT INTO lease_id (owner_address, dseq, gseq, oseq, provider_address) 
	VALUES ($1, $2, $3, $4, $5) 
	ON CONFLICT DO NOTHING 
	RETURNING id`

	var leaseID int64
	err := db.Sql.QueryRow(stmt,
		lease.LeaseID.Owner, lease.LeaseID.DSeq, lease.LeaseID.GSeq, lease.LeaseID.OSeq, lease.LeaseID.Provider,
	).Scan(&leaseID)
	if err != nil {
		return leaseID, fmt.Errorf("error while storing lease ID: %s", err)
	}

	return leaseID, nil
}

func (db *Db) saveLease(leaseID int64, l markettypes.Lease, height int64) error {
	stmt := `
	INSERT INTO lease (lease_id, lease_state, price, created_at, closed_on, height) 
	VALUES ($1, $2, $3, $4, $5, $6) 
	ON CONFLICT (lease_id) DO UPDATE 
	SET lease_state = excluded.lease_state, 
		price = excluded.price, 
		created_at = excluded.created_at,  
		closed_on = excluded.closed_on, 
		height = excluded.height 
	WHERE lease.height <= excluded.height`

	_, err := db.Sql.Exec(stmt,
		leaseID,
		l.State,
		pq.Array(dbtypes.NewDbDecCoin(l.Price)),
		l.CreatedAt,
		l.ClosedOn,
		height,
	)
	if err != nil {
		return fmt.Errorf("error while storing lease: %s", err)
	}

	return nil
}

func (db *Db) saveEscrowPayment(leaseID int64, e escrowtypes.FractionalPayment, height int64) error {
	stmt := `
	INSERT INTO escrow_payment (lease_id, account_id, payment_id, owner_address, payment_state, rate, balance, withdrawn, height) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
	ON CONFLICT (lease_id) DO UPDATE 
	SET account_id = excluded.account_id, 
		payment_id = excluded.payment_id, 
		owner_address = excluded.owner_address, 
		payment_state = excluded.payment_state, 
		rate = excluded.rate 
		balance = excluded.balance 
		withdrawn = excluded.withdrawn 
		height = excluded.height 
	WHERE escrow_payment.height <= excluded.height`

	accountID := dbtypes.NewDbLeaseAccountID(e.AccountID)
	accountIDValue, err := accountID.Value()
	if err != nil {
		return fmt.Errorf("error while converting account ID to DbLeaseAccountID value: %s", err)
	}

	_, err = db.Sql.Exec(stmt,
		leaseID,
		accountIDValue,
		e.PaymentID,
		e.Owner,
		e.State,
		pq.Array(dbtypes.NewDbDecCoin(e.Rate)),
		pq.Array(dbtypes.NewDbDecCoin(e.Balance)),
		pq.Array(dbtypes.NewDbCoin(e.Withdrawn)),
		height,
	)
	if err != nil {
		return fmt.Errorf("error while storing escrow payment: %s", err)
	}

	return nil
}
