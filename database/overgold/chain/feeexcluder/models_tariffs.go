package feeexcluder

import (
	"strconv"

	"git.ooo.ua/vipcoin/ovg-chain/x/feeexcluder/types"
	"github.com/shopspring/decimal"

	db "github.com/forbole/bdjuno/v4/database/types"
)

// BLOCK Fees

// toFeesDomain - mapping func to a domain model.
func toFeesDomain(f db.FeeExcluderFees) types.Fees {
	return types.Fees{
		AmountFrom:  strconv.FormatUint(f.AmountFrom, 10),
		Fee:         f.Fee.String(),
		RefReward:   f.RefReward.String(),
		StakeReward: f.StakeReward.String(),
		MinAmount:   f.MinAmount,
		NoRefReward: f.NoRefReward,
		Creator:     f.Creator,
		Id:          f.MsgID,
	}
}

// toFeesDomainList - mapping func to a domain list.
func toFeesDomainList(a []db.FeeExcluderFees) []*types.Fees {
	res := make([]*types.Fees, 0, len(a))
	for _, address := range a {
		d := toFeesDomain(address)
		res = append(res, &d)
	}

	return res
}

// toFeesDatabase - mapping func to a database model.
func toFeesDatabase(id uint64, f *types.Fees) (db.FeeExcluderFees, error) {
	amountFrom, err := strconv.ParseUint(f.AmountFrom, 10, 64)
	if err != nil {
		return db.FeeExcluderFees{}, err
	}

	fee, err := decimal.NewFromString(f.Fee)
	if err != nil {
		return db.FeeExcluderFees{}, err
	}

	refReward, err := decimal.NewFromString(f.RefReward)
	if err != nil {
		return db.FeeExcluderFees{}, err
	}

	stakeReward, err := decimal.NewFromString(f.StakeReward)
	if err != nil {
		return db.FeeExcluderFees{}, err
	}

	return db.FeeExcluderFees{
		NoRefReward: f.NoRefReward,
		ID:          id,
		MsgID:       f.Id,
		MinAmount:   f.MinAmount,
		AmountFrom:  amountFrom,
		Fee:         fee,
		RefReward:   refReward,
		StakeReward: stakeReward,
		Creator:     f.Creator,
	}, nil
}

// BLOCK Tariff

// toTariffDomain - mapping func to a domain model.
func toTariffDomain(f db.FeeExcluderTariff, fees []*types.Fees) *types.Tariff {
	return &types.Tariff{
		Id:            f.MsgID,
		Amount:        strconv.FormatUint(f.Amount, 10),
		Denom:         f.Denom,
		MinRefBalance: strconv.FormatUint(f.MinRefBalance, 10),
		Fees:          fees,
	}
}

// toTariffDomainList - mapping func to a domain list.
func toTariffDomainList(t []db.FeeExcluderTariff, fees []*types.Fees) []*types.Tariff {
	res := make([]*types.Tariff, 0, len(t))
	for _, tariff := range t {
		res = append(res, toTariffDomain(tariff, fees))
	}

	return res
}

// toTariffDatabase - mapping func to a database model.
func toTariffDatabase(id uint64, f *types.Tariff) (db.FeeExcluderTariff, error) {
	amount, err := strconv.ParseUint(f.Amount, 10, 64)
	if err != nil {
		return db.FeeExcluderTariff{}, err
	}

	minRefBalance, err := strconv.ParseUint(f.MinRefBalance, 10, 64)
	if err != nil {
		return db.FeeExcluderTariff{}, err
	}

	return db.FeeExcluderTariff{
		ID:            id,
		MsgID:         f.Id,
		Amount:        amount,
		MinRefBalance: minRefBalance,
		Denom:         f.Denom,
	}, nil
}

// BLOCK Tariffs

// toTariffsDomain - mapping func to a domain model.
func toTariffsDomain(t db.FeeExcluderTariffs, tariffs []*types.Tariff) types.Tariffs {
	return types.Tariffs{
		Denom:   t.Denom,
		Tariffs: tariffs,
		Creator: t.Creator,
	}
}

// toTariffsDomainList - mapping func to a domain list.
func toTariffsDomainList(t []db.FeeExcluderTariffs, tariffs []*types.Tariff) []types.Tariffs {
	res := make([]types.Tariffs, 0, len(t))
	for _, tariff := range t {
		res = append(res, toTariffsDomain(tariff, tariffs))
	}

	return res
}

// toTariffsDatabase - mapping func to a database model.
func toTariffsDatabase(id uint64, t types.Tariffs) db.FeeExcluderTariffs {
	return db.FeeExcluderTariffs{
		ID:      id,
		Denom:   t.Denom,
		Creator: t.Creator,
	}
}

// BLOCK MsgCreateTariffs

// toTariffsDomain - mapping func to a domain model.
func toMsgCreateTariffsDomain(t db.FeeExcluderCreateTariffs, tariff *types.Tariff) types.MsgCreateTariffs {
	return types.MsgCreateTariffs{
		Creator: t.Creator,
		Denom:   t.Denom,
		Tariff:  tariff,
	}
}

// toTariffsDatabase - mapping func to a database model.
func toMsgCreateTariffsDatabase(hash string, id, tariffID uint64, t types.MsgCreateTariffs) db.FeeExcluderCreateTariffs {
	return db.FeeExcluderCreateTariffs{
		ID:       id,
		TariffID: tariffID,
		TxHash:   hash,
		Creator:  t.Creator,
		Denom:    t.Denom,
	}
}

// BLOCK MsgUpdateTariffs

// toMsgUpdateTariffsDomain - mapping func to a domain model.
func toMsgUpdateTariffsDomain(t db.FeeExcluderUpdateTariffs, tariff *types.Tariff) types.MsgUpdateTariffs {
	return types.MsgUpdateTariffs{
		Denom:   t.Denom,
		Tariff:  tariff,
		Creator: t.Creator,
	}
}

// toTariffsDatabase - mapping func to a database model.
func toMsgUpdateTariffsDatabase(hash string, id, tariffID uint64, t types.MsgUpdateTariffs) db.FeeExcluderUpdateTariffs {
	return db.FeeExcluderUpdateTariffs{
		ID:       id,
		TariffID: tariffID,
		TxHash:   hash,
		Creator:  t.Creator,
		Denom:    t.Denom,
	}
}

// BLOCK MsgDeleteTariffs

// toMsgDeleteTariffsDomain - mapping func to a domain model.
func toMsgDeleteTariffsDomain(t db.FeeExcluderDeleteTariffs) types.MsgDeleteTariffs {
	return types.MsgDeleteTariffs{
		Creator:  t.Creator,
		Denom:    t.Denom,
		TariffID: strconv.FormatUint(t.TariffID, 10),
		FeeID:    strconv.FormatUint(t.FeesID, 10),
	}
}

// toTariffsDomainList - mapping func to a domain list.
func toMsgDeleteTariffsDomainList(t []db.FeeExcluderDeleteTariffs) []types.MsgDeleteTariffs {
	res := make([]types.MsgDeleteTariffs, 0, len(t))
	for _, tariff := range t {
		res = append(res, toMsgDeleteTariffsDomain(tariff))
	}

	return res
}

// toTariffsDatabase - mapping func to a database model.
func toMsgDeleteTariffsDatabase(hash string, id uint64, t types.MsgDeleteTariffs) (db.FeeExcluderDeleteTariffs, error) {
	tariffID, err := strconv.ParseUint(t.TariffID, 10, 64)
	if err != nil {
		return db.FeeExcluderDeleteTariffs{}, err
	}

	feeID, err := strconv.ParseUint(t.FeeID, 10, 64)
	if err != nil {
		return db.FeeExcluderDeleteTariffs{}, err
	}

	return db.FeeExcluderDeleteTariffs{
		ID:       id,
		TariffID: tariffID,
		FeesID:   feeID,
		TxHash:   hash,
		Creator:  t.Creator,
		Denom:    t.Denom,
	}, nil
}
