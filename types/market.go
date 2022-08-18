package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	markettypes "github.com/ovrclk/akash/x/market/types/v1beta2"
)

// ------- x/market params -------

// MarketParams represents the x/market parameters
type MarketParams struct {
	markettypes.Params
	Height int64
}

// NewMarketParams allows to build a new MarketParams instance
func NewMarketParams(params markettypes.Params, height int64) *MarketParams {
	return &MarketParams{
		Params: params,
		Height: height,
	}
}

// ------- x/market lease -------
type MarketLease struct {
	Owner     string
	DSeq      uint64
	GSeq      uint32
	OSeq      uint32
	Provider  string
	State     int32
	Price     sdk.DecCoin
	CreatedAt int64
	ClosedOn  int64
	Height    int64
}

// NewMarketLease allows to build a new MarketLease instance
func NewMarketLease(res markettypes.QueryLeaseResponse, height int64) *MarketLease {
	return &MarketLease{
		Owner:     res.Lease.LeaseID.Owner,
		DSeq:      res.Lease.LeaseID.DSeq,
		GSeq:      res.Lease.LeaseID.GSeq,
		OSeq:      res.Lease.LeaseID.OSeq,
		Provider:  res.Lease.LeaseID.Provider,
		State:     int32(res.Lease.State),
		Price:     res.Lease.Price,
		CreatedAt: res.Lease.CreatedAt,
		ClosedOn:  res.Lease.ClosedOn,
		Height:    height,
	}
}
