package chain

import (
	"git.ooo.ua/vipcoin/lib/filter"
	allowed "git.ooo.ua/vipcoin/ovg-chain/x/allowed/types"
	core "git.ooo.ua/vipcoin/ovg-chain/x/core/types"
	referral "git.ooo.ua/vipcoin/ovg-chain/x/referral/types"
	stake "git.ooo.ua/vipcoin/ovg-chain/x/stake/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// custom ovg types
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
		GetAllMsgIssue(filter filter.Filter) ([]core.MsgIssue, error)
		InsertMsgIssue(hash string, msgs ...core.MsgIssue) error

		GetAllMsgWithdraw(filter filter.Filter) ([]core.MsgWithdraw, error)
		InsertMsgWithdraw(hash string, msgs ...core.MsgWithdraw) error

		GetAllMsgSend(filter filter.Filter) ([]core.MsgSend, error)
		InsertMsgSend(hash string, msgs ...core.MsgSend) error
	}

	// FeeExcluder - describes an interface for working with database models.
	FeeExcluder interface {
		// TODO
	}

	// Referral - describes an interface for working with database models.
	Referral interface {
		GetAllMsgSetReferrer(filter filter.Filter) ([]referral.MsgSetReferrer, error)
		InsertMsgSetReferrer(hash string, msgs ...referral.MsgSetReferrer) error
	}

	// Stake - describes an interface for working with database models.
	Stake interface {
		GetAllMsgSell(filter filter.Filter) ([]stake.MsgSellRequest, error)
		InsertMsgSell(hash string, msgs ...stake.MsgSellRequest) error

		GetAllMsgBuy(filter filter.Filter) ([]stake.MsgBuyRequest, error)
		InsertMsgBuy(hash string, msgs ...stake.MsgBuyRequest) error

		GetAllMsgSellCancel(filter filter.Filter) ([]stake.MsgMsgCancelSell, error)
		InsertMsgSellCancel(hash string, msgs ...stake.MsgMsgCancelSell) error

		GetAllMsgClaimReward(filter filter.Filter) ([]stake.MsgClaimReward, error)
		InsertMsgClaimReward(hash string, msgs ...stake.MsgClaimReward) error

		GetAllMsgDistributeRewards(filter filter.Filter) ([]stake.MsgDistributeRewards, error)
		InsertMsgDistributeRewards(hash string, msgs ...stake.MsgDistributeRewards) error
	}
)

// custom sdk types
type (
	// Bank - describes an interface for working with database models.
	Bank interface {
		GetAllMsgMultiSend(filter filter.Filter) ([]bank.MsgMultiSend, error)
		InsertMsgMultiSend(hash string, msgs ...bank.MsgMultiSend) error

		GetAllMsgSend(filter filter.Filter) ([]bank.MsgSend, error)
		InsertMsgSend(hash string, msgs ...bank.MsgSend) error
	}

	// LastBlock - describes an interface for working with database models.
	LastBlock interface {
		Get() (uint64, error)
		Update(id uint64) error
	}
)
