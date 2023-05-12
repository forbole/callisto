package remote

// import "github.com/tendermint/tendermint/proto/tendermint/crypto"

type PK struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type Validators struct {
	Address          string `json:"address"`
	PubKey           PK     `json:"pub_key"`
	VotingPower      string `json:"voting_power"`
	ProposerPriority string `json:"proposer_priority"`
}

type ResultValidators struct {
	BlockHeight string       `json:"block_height"`
	Validators  []Validators `json:"validators"`
	Count       string       `json:"count"`
	Total       string       `json:"total"`
}

type ValidatorQueryResult struct {
	JSONRPC string           `json:"jsonrpc"`
	ID      int64            `json:"id"`
	Result  ResultValidators `json:"result"`
}
