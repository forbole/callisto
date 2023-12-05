package allowed

import (
	"git.ooo.ua/vipcoin/ovg-chain/x/allowed/types"

	db "github.com/forbole/bdjuno/v4/database/types"
)

// BLOCK AllowedAddresses

// toAddressesDomain - mapping func to a domain model.
func toAddressesDomain(a db.AllowedAddresses) types.Addresses {
	return types.Addresses{
		Id:      a.ID,
		Address: a.Address,
		Creator: a.Creator,
	}
}

// toAddressesDomainList - mapping func to a domain list.
func toAddressesDomainList(a []db.AllowedAddresses) []types.Addresses {
	res := make([]types.Addresses, 0, len(a))
	for _, address := range a {
		res = append(res, toAddressesDomain(address))
	}

	return res
}

// toAddressesDatabase - mapping func to a database model.
func toAddressesDatabase(a types.Addresses) db.AllowedAddresses {
	return db.AllowedAddresses{
		ID:      a.Id,
		Creator: a.Creator,
		Address: a.Address,
	}
}

// BLOCK AllowedCreateAddresses

// toCreateAddressesDatabase - mapping func to a database model.
func toCreateAddressesDatabase(hash string, m *types.MsgCreateAddresses) db.AllowedCreateAddresses {
	return db.AllowedCreateAddresses{
		TxHash:  hash,
		Creator: m.Creator,
		Address: m.Address,
	}
}

// toCreateAddressesDomain - mapping func to a domain model.
func toCreateAddressesDomain(a db.AllowedCreateAddresses) types.MsgCreateAddresses {
	return types.MsgCreateAddresses{
		Creator: a.Creator,
		Address: a.Address,
	}
}

// toCreateAddressesDomainList - mapping func to a domain list.
func toCreateAddressesDomainList(a []db.AllowedCreateAddresses) []types.MsgCreateAddresses {
	res := make([]types.MsgCreateAddresses, 0, len(a))
	for _, address := range a {
		res = append(res, toCreateAddressesDomain(address))
	}

	return res
}

// BLOCK AllowedDeleteByAddresses

// toDeleteByAddressesDatabase - mapping func to a database model.
func toDeleteByAddressesDatabase(hash string, m *types.MsgDeleteByAddresses) db.AllowedDeleteByAddresses {
	return db.AllowedDeleteByAddresses{
		TxHash:  hash,
		Creator: m.Creator,
		Address: m.Address,
	}
}

// toDeleteByAddressesDomain - mapping func to a domain model.
func toDeleteByAddressesDomain(a db.AllowedDeleteByAddresses) types.MsgDeleteByAddresses {
	return types.MsgDeleteByAddresses{
		Creator: a.Creator,
		Address: a.Address,
	}
}

// toDeleteByAddressesDomainList - mapping func to a domain list.
func toDeleteByAddressesDomainList(a []db.AllowedDeleteByAddresses) []types.MsgDeleteByAddresses {
	res := make([]types.MsgDeleteByAddresses, 0, len(a))
	for _, address := range a {
		res = append(res, toDeleteByAddressesDomain(address))
	}

	return res
}

// BLOCK AllowedDeleteByID

// toDeleteByIDDatabase - mapping func to a database model.
func toDeleteByIDDatabase(hash string, m *types.MsgDeleteByID) db.AllowedDeleteByID {
	return db.AllowedDeleteByID{
		ID:      m.Id,
		TxHash:  hash,
		Creator: m.Creator,
	}
}

// toDeleteByIDDomain - mapping func to a domain model.
func toDeleteByIDDomain(a db.AllowedDeleteByID) types.MsgDeleteByID {
	return types.MsgDeleteByID{
		Creator: a.Creator,
		Id:      a.ID,
	}
}

// toDeleteByIDDomainList - mapping func to a domain list.
func toDeleteByIDDomainList(a []db.AllowedDeleteByID) []types.MsgDeleteByID {
	res := make([]types.MsgDeleteByID, 0, len(a))
	for _, id := range a {
		res = append(res, toDeleteByIDDomain(id))
	}

	return res
}

// BLOCK AllowedUpdateAddresses

// toUpdateAddressesDatabase - mapping func to a database model.
func toUpdateAddressesDatabase(hash string, m *types.MsgUpdateAddresses) db.AllowedUpdateAddresses {
	return db.AllowedUpdateAddresses{
		ID:      m.Id,
		TxHash:  hash,
		Creator: m.Creator,
		Address: m.Address,
	}
}

// toUpdateAddressesDomain - mapping func to a domain model.
func toUpdateAddressesDomain(a db.AllowedUpdateAddresses) types.MsgUpdateAddresses {
	return types.MsgUpdateAddresses{
		Creator: a.Creator,
		Id:      a.ID,
		Address: a.Address,
	}
}

// toUpdateAddressesDomainList - mapping func to a domain list.
func toUpdateAddressesDomainList(a []db.AllowedUpdateAddresses) []types.MsgUpdateAddresses {
	res := make([]types.MsgUpdateAddresses, 0, len(a))
	for _, address := range a {
		res = append(res, toUpdateAddressesDomain(address))
	}

	return res
}
