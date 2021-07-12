package types

import( 
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtype "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
)

// Account represents a chain account
type Account struct {
	Address string
}

// NewAccount builds a new Account instance
func NewAccount(address string) Account {
	return Account{
		Address: address,
	}
}

type PeriodicVestingAccount struct{
	Address          string
	PubKey           string
	AccountNumber    uint64  
	Sequence         uint64
	OriginalVesting  sdk.Coins
	DelegatedFree    sdk.Coins
	DelegatedVesting sdk.Coins
	EndTime          int64
	StartTime        int64
	VestingPeriods   authtype.Periods       
}
func NewPeriodicVestingAccount (address string,pubKey string,accountNumber uint64,sequence uint64,
	originalVesting sdk.Coins,delegatedFree sdk.Coins,delegatedVesting sdk.Coins,
	endTime int64,startTime int64,vestingPeriods authtype.Periods,
	)PeriodicVestingAccount{
		return PeriodicVestingAccount{
	Address          :address         ,
	PubKey           :pubKey          ,
	AccountNumber    :accountNumber   ,
	Sequence         :sequence        ,
	OriginalVesting  :originalVesting ,
	DelegatedFree    :delegatedFree   ,
	DelegatedVesting :delegatedVesting,
	EndTime          :endTime         ,
	StartTime        :startTime       ,
	VestingPeriods   :vestingPeriods ,
		} 
}