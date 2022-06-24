package banking

import (
	"errors"
	"strconv"

	banking "git.ooo.ua/vipcoin/chain/x/banking/types"
)

type baseTransfer struct {
	Id        string `json:"id,omitempty"`
	Asset     string `json:"asset,omitempty"`
	Amount    string `json:"amount,omitempty"`
	Kind      string `json:"kind,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
	TxHash    string `json:"tx_hash,omitempty"`
}

// toVipcoinBaseTransfer - convert to vipcoin base transfer
func (b *baseTransfer) toVipcoinBaseTransfer() (banking.BaseTransfer, error) {
	id, err := strconv.ParseUint(b.Id, 10, 64)
	if err != nil {
		return banking.BaseTransfer{}, err
	}

	amount, err := strconv.ParseUint(b.Amount, 10, 64)
	if err != nil {
		return banking.BaseTransfer{}, err
	}

	kind, ok := banking.TransferKind_value[b.Kind]
	if !ok {
		return banking.BaseTransfer{}, errors.New("unknown kind")
	}

	timestamp, err := strconv.ParseInt(b.Timestamp, 10, 64)
	if err != nil {
		return banking.BaseTransfer{}, err
	}

	result := banking.BaseTransfer{
		Id:        id,
		Asset:     b.Asset,
		Amount:    amount,
		Kind:      banking.TransferKind(kind),
		Timestamp: timestamp,
		TxHash:    b.TxHash,
	}

	return result, nil
}
