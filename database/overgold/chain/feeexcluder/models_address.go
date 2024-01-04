package feeexcluder

import (
	"git.ooo.ua/vipcoin/ovg-chain/x/feeexcluder/types"

	"github.com/forbole/bdjuno/v4/database/overgold/chain"
	db "github.com/forbole/bdjuno/v4/database/types"
)

// BLOCK Address

// toAddressDomain - mapping func to a domain model.
func toAddressDomain(a db.FeeExcluderAddress) types.Address {
	return types.Address{
		Id:      uint64(a.MsgID.Int64),
		Address: a.Address,
		Creator: a.Creator,
	}
}

// toAddressDomainList - mapping func to a domain list.
func toAddressDomainList(a []db.FeeExcluderAddress) []types.Address {
	res := make([]types.Address, 0, len(a))
	for _, address := range a {
		res = append(res, toAddressDomain(address))
	}

	return res
}

// toAddressDatabase - mapping func to a database model.
func toAddressDatabase(id uint64, a types.Address) db.FeeExcluderAddress {
	return db.FeeExcluderAddress{
		ID:      id,
		MsgID:   chain.ToNullInt64(int64(a.Id)),
		Address: a.Address,
		Creator: a.Creator,
	}
}

// BLOCK MsgCreateAddress

// toAddressDomain - mapping func to a domain model.
func toMsgCreateAddressDomain(a db.FeeExcluderCreateAddress) types.MsgCreateAddress {
	return types.MsgCreateAddress{
		Address: a.Address,
		Creator: a.Creator,
	}
}

// toAddressDomainList - mapping func to a domain list.
func toMsgCreateAddressDomainList(a []db.FeeExcluderCreateAddress) []types.MsgCreateAddress {
	res := make([]types.MsgCreateAddress, 0, len(a))
	for _, address := range a {
		res = append(res, toMsgCreateAddressDomain(address))
	}

	return res
}

// toAddressDatabase - mapping func to a database model.
func toMsgCreateAddressDatabase(hash string, id uint64, a types.MsgCreateAddress) db.FeeExcluderCreateAddress {
	return db.FeeExcluderCreateAddress{
		ID:      id,
		TxHash:  hash,
		Creator: a.Creator,
		Address: a.Address,
	}
}

// BLOCK MsgUpdateAddress

// toMsgUpdateAddressDomain - mapping func to a domain model.
func toMsgUpdateAddressDomain(a db.FeeExcluderUpdateAddress) types.MsgUpdateAddress {
	return types.MsgUpdateAddress{
		Creator: a.Creator,
		Id:      a.ID,
		Address: a.Address,
	}
}

// toAddressDomainList - mapping func to a domain list.
func toMsgUpdateAddressDomainList(a []db.FeeExcluderUpdateAddress) []types.MsgUpdateAddress {
	res := make([]types.MsgUpdateAddress, 0, len(a))
	for _, address := range a {
		res = append(res, toMsgUpdateAddressDomain(address))
	}

	return res
}

// toAddressDatabase - mapping func to a database model.
func toMsgUpdateAddressDatabase(hash string, a types.MsgUpdateAddress) db.FeeExcluderUpdateAddress {
	return db.FeeExcluderUpdateAddress{
		ID:      a.Id,
		TxHash:  hash,
		Creator: a.Creator,
		Address: a.Address,
	}
}

// BLOCK MsgDeleteAddress

// toMsgDeleteAddressDomain - mapping func to a domain model.
func toMsgDeleteAddressDomain(a db.FeeExcluderDeleteAddress) types.MsgDeleteAddress {
	return types.MsgDeleteAddress{
		Creator: a.Creator,
		Id:      a.ID,
	}
}

// toAddressDomainList - mapping func to a domain list.
func toMsgDeleteAddressDomainList(a []db.FeeExcluderDeleteAddress) []types.MsgDeleteAddress {
	res := make([]types.MsgDeleteAddress, 0, len(a))
	for _, address := range a {
		res = append(res, toMsgDeleteAddressDomain(address))
	}

	return res
}

// toAddressDatabase - mapping func to a database model.
func toMsgDeleteAddressDatabase(hash string, a types.MsgDeleteAddress) db.FeeExcluderDeleteAddress {
	return db.FeeExcluderDeleteAddress{
		ID:      a.Id,
		TxHash:  hash,
		Creator: a.Creator,
	}
}
