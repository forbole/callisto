package chain

import (
	"git.ooo.ua/vipcoin/lib/filter"
	allowed "git.ooo.ua/vipcoin/ovg-chain/x/allowed/types"
	core "git.ooo.ua/vipcoin/ovg-chain/x/core/types"
	fe "git.ooo.ua/vipcoin/ovg-chain/x/feeexcluder/types"
	referral "git.ooo.ua/vipcoin/ovg-chain/x/referral/types"
	stake "git.ooo.ua/vipcoin/ovg-chain/x/stake/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/jmoiron/sqlx"

	"github.com/forbole/bdjuno/v4/database/types"
)

// custom ovg types
type (
	// Allowed - describes an interface for working with database models.
	Allowed interface {
		DeleteAddressesByAddress(addresses ...string) error
		DeleteAddressesByID(ids ...uint64) error
		GetAllAddresses(filter filter.Filter) ([]allowed.Addresses, error)
		InsertToAddresses(addresses ...allowed.Addresses) error
		UpdateAddresses(addresses ...allowed.Addresses) error

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
		FeeExcluderM2MTables
		FeeExcluderLinkedTables

		DeleteMsgCreateAddress(id uint64) error
		GetAllMsgCreateAddress(filter filter.Filter) ([]fe.MsgCreateAddress, error)
		InsertToMsgCreateAddress(hash string, address fe.MsgCreateAddress) error
		UpdateMsgCreateAddress(hash string, id uint64, address fe.MsgCreateAddress) error

		DeleteMsgUpdateAddress(id uint64) error
		GetAllMsgUpdateAddress(filter filter.Filter) ([]fe.MsgUpdateAddress, error)
		InsertToMsgUpdateAddress(hash string, addresses fe.MsgUpdateAddress) error
		UpdateMsgUpdateAddress(hash string, addresses ...fe.MsgUpdateAddress) error

		DeleteMsgDeleteAddress(id uint64) error
		GetAllMsgDeleteAddress(filter filter.Filter) ([]fe.MsgDeleteAddress, error)
		InsertToMsgDeleteAddress(hash string, addresses ...fe.MsgDeleteAddress) error
		UpdateMsgDeleteAddress(hash string, addresses ...fe.MsgDeleteAddress) error

		GetAllMsgCreateTariffs(f filter.Filter) ([]fe.MsgCreateTariffs, error)
		InsertToMsgCreateTariffs(hash string, ct fe.MsgCreateTariffs) error

		DeleteMsgUpdateTariffs(id uint64) error
		GetAllMsgUpdateTariffs(f filter.Filter) ([]fe.MsgUpdateTariffs, error)
		InsertToMsgUpdateTariffs(hash string, ut fe.MsgUpdateTariffs) error
		UpdateMsgUpdateTariffs(hash string, id uint64, ut fe.MsgUpdateTariffs) error

		DeleteMsgDeleteTariffs(id uint64) error
		GetAllMsgDeleteTariffs(f filter.Filter) ([]fe.MsgDeleteTariffs, error)
		InsertToMsgDeleteTariffs(hash string, dt fe.MsgDeleteTariffs) error
		UpdateMsgDeleteTariffs(hash string, id uint64, ut fe.MsgDeleteTariffs) error

		DeleteGenesisState(id uint64) error
		GetAllGenesisState(filter filter.Filter) ([]fe.GenesisState, error)
		InsertToGenesisState(gsList fe.GenesisState) error
	}

	FeeExcluderLinkedTables interface {
		DeleteAddress(tx *sqlx.Tx, id uint64) error
		GetAllAddress(filter filter.Filter) ([]fe.Address, error)
		InsertToAddress(tx *sqlx.Tx, addresses fe.Address) (uint64, error)
		UpdateAddress(tx *sqlx.Tx, id uint64, address fe.Address) error

		DeleteFees(tx *sqlx.Tx, id uint64) error
		GetAllFees(filter filter.Filter) ([]*fe.Fees, error)
		InsertToFees(tx *sqlx.Tx, fees *fe.Fees) (uint64, error)
		UpdateFees(tx *sqlx.Tx, id uint64, fees *fe.Fees) error

		DeleteStats(tx *sqlx.Tx, id string) error
		GetAllStats(filter filter.Filter) ([]fe.Stats, error)
		InsertToStats(tx *sqlx.Tx, stats fe.Stats) (string, error)
		UpdateStats(tx *sqlx.Tx, stats fe.Stats) error

		DeleteDailyStats(tx *sqlx.Tx, id uint64) error
		GetAllDailyStats(f filter.Filter) ([]fe.DailyStats, error)
		InsertToDailyStats(tx *sqlx.Tx, dailyStats fe.DailyStats) (uint64, error)
		UpdateDailyStats(tx *sqlx.Tx, id uint64, ut fe.DailyStats) error

		DeleteTariff(tx *sqlx.Tx, id uint64) error
		GetAllTariff(f filter.Filter) ([]*fe.Tariff, error)
		InsertToTariff(tx *sqlx.Tx, tariff *fe.Tariff) (uint64, error)
		UpdateTariff(tx *sqlx.Tx, id uint64, tariff *fe.Tariff) error

		DeleteTariffs(tx *sqlx.Tx, id uint64) error
		GetAllTariffs(f filter.Filter) ([]fe.Tariffs, error)
		GetTariffsDB(req filter.Filter) (types.FeeExcluderTariffs, error)
		GetAllTariffsDB(f filter.Filter) ([]types.FeeExcluderTariffs, error)
		InsertToTariffs(tx *sqlx.Tx, tariffs fe.Tariffs) (uint64, error)
		UpdateTariffs(tx *sqlx.Tx, id uint64, tariffs fe.Tariffs) error
	}

	FeeExcluderM2MTables interface {
		DeleteM2MTariffFeesByTariff(tx *sqlx.Tx, tariffID uint64) error
		GetAllM2MTariffFees(filter filter.Filter) ([]types.FeeExcluderM2MTariffFees, error)
		InsertToM2MTariffFees(tx *sqlx.Tx, ids ...types.FeeExcluderM2MTariffFees) error

		DeleteM2MTariffTariffsByTariffs(tx *sqlx.Tx, id uint64) error
		GetAllM2MTariffTariffs(filter filter.Filter) ([]types.FeeExcluderM2MTariffTariffs, error)
		InsertToM2MTariffTariffs(tx *sqlx.Tx, ids ...types.FeeExcluderM2MTariffTariffs) error

		DeleteM2MGenesisStateAddressByGenesisState(tx *sqlx.Tx, id uint64) error
		GetAllM2MGenesisStateAddress(filter filter.Filter) ([]types.FeeExcluderM2MGenesisStateAddress, error)
		InsertToM2MGenesisStateAddress(tx *sqlx.Tx, ids ...types.FeeExcluderM2MGenesisStateAddress) error

		DeleteM2MGenesisStateDailyStatsByGenesisState(tx *sqlx.Tx, id uint64) error
		GetAllM2MGenesisStateDailyStats(filter filter.Filter) ([]types.FeeExcluderM2MGenesisStateDailyStats, error)
		InsertToM2MGenesisStateDailyStats(tx *sqlx.Tx, ids ...types.FeeExcluderM2MGenesisStateDailyStats) error

		DeleteM2MGenesisStateStatsByGenesisState(tx *sqlx.Tx, id uint64) error
		GetAllM2MGenesisStateStats(filter filter.Filter) ([]types.FeeExcluderM2MGenesisStateStats, error)
		InsertToM2MGenesisStateStats(tx *sqlx.Tx, ids ...types.FeeExcluderM2MGenesisStateStats) error

		DeleteM2MGenesisStateTariffsByGenesisState(tx *sqlx.Tx, id uint64) error
		GetAllM2MGenesisStateTariffs(filter filter.Filter) ([]types.FeeExcluderM2MGenesisStateTariffs, error)
		InsertToM2MGenesisStateTariffs(tx *sqlx.Tx, ids ...types.FeeExcluderM2MGenesisStateTariffs) error
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
