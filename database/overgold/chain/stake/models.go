package stake

import (
	"strconv"

	"git.ooo.ua/vipcoin/lib/errs"
	chainDomain "git.ooo.ua/vipcoin/ovg-chain/x/domain"
	"git.ooo.ua/vipcoin/ovg-chain/x/stake/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	db "github.com/forbole/bdjuno/v4/database/types"
)

// toMsgSellDomain - mapping func to a domain model.
func toMsgSellDomain(m db.StakeMsgSell) types.MsgSellRequest {
	return types.MsgSellRequest{
		Creator: m.Creator,
		Amount:  strconv.FormatUint(m.Amount, 10),
	}
}

// toMsgSellDomainList - mapping func to a domain list.
func toMsgSellDomainList(m []db.StakeMsgSell) []types.MsgSellRequest {
	res := make([]types.MsgSellRequest, 0, len(m))
	for _, msg := range m {
		res = append(res, toMsgSellDomain(msg))
	}

	return res
}

// toMsgSellDatabase - mapping func to a database model.
func toMsgSellDatabase(hash string, m types.MsgSellRequest) (db.StakeMsgSell, error) {
	amount, err := strconv.ParseUint(m.Amount, 10, 64)
	if err != nil {
		return db.StakeMsgSell{}, errs.Internal{Cause: err.Error()}
	}

	return db.StakeMsgSell{
		TxHash:  hash,
		Creator: m.Creator,
		Amount:  amount,
	}, nil
}

// toMsgSellCancel - mapping func to a domain model.
func toMsgSellCancel(m db.StakeMsgSellCancel) types.MsgMsgCancelSell {
	return types.MsgMsgCancelSell{
		Creator: m.Creator,
		Amount:  sdk.NewCoin(chainDomain.DenomSTOVG, sdk.NewIntFromUint64(m.Amount)),
	}
}

// toMsgSellCancelDomainList - mapping func to a domain list.
func toMsgSellCancelDomainList(m []db.StakeMsgSellCancel) []types.MsgMsgCancelSell {
	res := make([]types.MsgMsgCancelSell, 0, len(m))
	for _, msg := range m {
		res = append(res, toMsgSellCancel(msg))
	}

	return res
}

// toMsgSellCancelDatabase - mapping func to a database model.
func toMsgSellCancelDatabase(hash string, m types.MsgMsgCancelSell) (db.StakeMsgSellCancel, error) {
	return db.StakeMsgSellCancel{
		TxHash:  hash,
		Creator: m.Creator,
		Amount:  m.Amount.Amount.Uint64(),
	}, nil
}

// toMsgBuy - mapping func to a domain model.
func toMsgBuy(m db.StakeMsgBuy) types.MsgBuyRequest {
	return types.MsgBuyRequest{
		Creator: m.Creator,
		Amount:  strconv.FormatUint(m.Amount, 10),
	}
}

// toMsgBuyDomainList - mapping func to a domain list.
func toMsgBuyDomainList(m []db.StakeMsgBuy) []types.MsgBuyRequest {
	res := make([]types.MsgBuyRequest, 0, len(m))
	for _, msg := range m {
		res = append(res, toMsgBuy(msg))
	}

	return res
}

// toMsgBuyDatabase - mapping func to a database model.
func toMsgBuyDatabase(hash string, m types.MsgBuyRequest) (db.StakeMsgBuy, error) {
	amount, err := strconv.ParseUint(m.Amount, 10, 64)
	if err != nil {
		return db.StakeMsgBuy{}, errs.Internal{Cause: err.Error()}
	}
	return db.StakeMsgBuy{
		TxHash:  hash,
		Creator: m.Creator,
		Amount:  amount,
	}, nil
}

// toMsgBuy - mapping func to a domain model.
func toMsgDistribute(m db.StakeMsgDistribute) types.MsgDistributeRewards {
	return types.MsgDistributeRewards{
		Creator: m.Creator,
	}
}

// toMsgDistributeDomainList - mapping func to a domain list.
func toMsgDistributeDomainList(m []db.StakeMsgDistribute) []types.MsgDistributeRewards {
	res := make([]types.MsgDistributeRewards, 0, len(m))
	for _, msg := range m {
		res = append(res, toMsgDistribute(msg))
	}

	return res
}

// toMsgDistributeDatabase - mapping func to a database model.
func toMsgDistributeDatabase(hash string, m types.MsgDistributeRewards) db.StakeMsgDistribute {
	return db.StakeMsgDistribute{
		TxHash:  hash,
		Creator: m.Creator,
	}
}

// toMsgClaimReward - mapping func to a domain model.
func toMsgClaimReward(m db.StakeMsgClaim) types.MsgClaimReward {
	return types.MsgClaimReward{
		Creator: m.Creator,
		Amount:  sdk.NewCoin(chainDomain.DenomSTOVG, sdk.NewIntFromUint64(m.Amount)),
	}
}

// toMsgClaimRewardDomainList - mapping func to a domain list.
func toMsgClaimRewardDomainList(m []db.StakeMsgClaim) []types.MsgClaimReward {
	res := make([]types.MsgClaimReward, 0, len(m))
	for _, msg := range m {
		res = append(res, toMsgClaimReward(msg))
	}

	return res
}

// toMsgClaimRewardDatabase - mapping func to a database model.
func toMsgClaimRewardDatabase(hash string, m types.MsgClaimReward) db.StakeMsgClaim {
	return db.StakeMsgClaim{
		TxHash:  hash,
		Creator: m.Creator,
		Amount:  m.Amount.Amount.Uint64(),
	}
}
