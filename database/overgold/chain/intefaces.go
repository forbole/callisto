package chain

import (
	"git.ooo.ua/vipcoin/lib/filter"
	allowed "git.ooo.ua/vipcoin/ovg-chain/x/allowed/types"
)

type (
	// Allowed - describes an interface for working with database models.
	Allowed interface {
		DeleteAddressesByAddress(addresses ...string) error
		DeleteAddressesByID(ids ...uint64) error
		GetAllAddresses(filter filter.Filter) ([]allowed.Addresses, error)
		InsertToAddresses(addresses ...allowed.Addresses) error
		UpdateAddresses(assets ...allowed.Addresses) error

		GetAllCreateAddresses(filter filter.Filter) ([]allowed.MsgCreateAddresses, error)
		InsertToCreateAddresses(hash string, msgs ...*allowed.MsgCreateAddresses) error

		GetAllDeleteByAddresses(filter filter.Filter) ([]allowed.MsgDeleteByAddresses, error)
		InsertToDeleteByAddresses(hash string, msgs ...*allowed.MsgDeleteByAddresses) error

		GetAllDeleteByID(filter filter.Filter) ([]allowed.MsgDeleteByID, error)
		InsertToDeleteByID(hash string, msgs ...*allowed.MsgDeleteByID) error

		GetAllUpdateAddresses(filter filter.Filter) ([]allowed.MsgUpdateAddresses, error)
		InsertToUpdateAddresses(hash string, msgs ...*allowed.MsgUpdateAddresses) error
	}

	// Core - describes an interface for working with database models.
	Core interface {
		// TODO
	}

	// FeeExcluder - describes an interface for working with database models.
	FeeExcluder interface {
		// TODO
	}

	// LastBlock - describes an interface for working with database models.
	LastBlock interface {
		// TODO
	}

	// Referral - describes an interface for working with database models.
	Referral interface {
		// TODO
	}

	// Stake - describes an interface for working with database models.
	Stake interface {
		// TODO
	}
)
