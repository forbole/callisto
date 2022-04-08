# Requirements of BDjuno
The BDjuno is working as a backend for Big Dipper. The followings are the features currently supported in Big Dipper in the backend which BDjuno should adapt.

## On every block
### Done by Juno by default
- [x] Parsing all blocks
- [x] Parsing all transactions
- [x] Store validator set of the block

### Custom BDJuno implementations
- [x] Update missed block records
- [x] Read the latest consensus state
- [x] [x/auth] Store vesting accounts and vesting periods details
- [x] [x/distribution] Update community pool
- [x] [x/feegrant] Store feegrant allowance details
- [x] [x/gov] Get gov proposals, deposits and votes
- [x] [x/gov] Calculate the tally result
- [x] [x/mint] Update the inflation
- [x] [x/slashing] Get validators signing info
- [x] [x/staking] Update validator information 
- [x] [x/staking] Calculate validator voting power percentage 
- [x] [x/staking] Update the total staked tokens 
- [x] [x/staking] Update the double sign evidences
- [x] [x/supply] Update the total supply


### Achievable using GraphQL APIs
- [x] Calculate the average block time


### Achievable using Hasura Actions
Address/Delegator related data:
- [x] Get account balance
- [x] Get delegations
- [x] Get total delegations amount
- [x] Get delegation rewards
- [x] Get unbonding delegations
- [x] Get total unbonding delegations amount
- [x] Get redelegations
- [x] Get delegator withdraw address

Validator related data:
- [x] Get commission amount
- [x] Get validator delegations
- [x] Get validator redelegations
- [x] Get validator unbonding delegations


## On intervals
- [x] [x/bank] Get total supply (per 10 mins)
- [x] [x/distribution] Get community pool (per hour)
- [x] [x/mint] Get inflation (per day)
- [x] [x/pricefeed] Get token price and marketcap (per 2 minutes, per hour)
- [x] [x/staking] Calculate average delegation ratio (per hour, per day) *
- [x] [x/staking] Calculate voting power distribution (per hour) *

\* These should be doable using the `average` method inside GraphQL

## Not on Big Dipper now but we are considering to add

- [ ] All wallets activities
- [ ] Alert on events: 
   - [ ] Proposal creation
   - [ ] Slashing
   - [ ] Huge delegation
   - [ ] Validator low uptime
   - [ ] Huge undelegation
   - [ ] Proposal start voting 
   - [ ] Proposal voting ends
- [ ] Validators rating
   - [ ] Self-delegation
   - [ ] Uptime
   - [ ] Ever slashed
   - [ ] Gov participation
   - [ ] Community contributions
   - [ ] Number of delegators
