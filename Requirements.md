# Requirements of BDjuno
The BDjuno is working as a backend for Big Dipper. The followings are the features currently supported in Big Dipper in the backend which BDjuno should adapt.

## On every block
### Done by Juno by default
- [x] Parsing all blocks
- [x] Parsing all transactions
- [x] Store validator set of the block

### Custom BDJuno implementations
- [x] Update miss block records
- [x] Read the latest consensus state
- [x] [x/staking] Update validator information 
- [x] [x/staking] Calculate validator voting power percentage 
- [x] [x/staking] Calculate validator self delegation ratio
- [x] [x/staking] Update the total staked tokens 
- [x] [x/supply] Update the total supply
- [x] [x/mint] Update the inflation
- [x] [x/distribution] Update community pool
- [x] [x/gov] Get gov proposals
- [x] [x/gov] Calculate the tally result

### Achievable using GraphQL APIs
- [x] Calculate the average block time

## On intervals
- [x] Get token price and marketcap (per 30 seconds)
- [x] [x/staking] Calculate average delegation ratio (per hour, per day) *
- [x] [x/staking] Calculate voting power distribution (per hour) *
- [x] [x/staking] Record all delegations (per day) *
- [x] [x/staking] Record all undelegatios (per day) *
- [x] [x/staking] Record all redelegations (per day) *

\* These should be duable using the `average` method inside GraphQL

## Not on Big Dipper now but we are considering to add
- [x] Validators signing-info (slashing)
- [ ] All wallets activities
- [ ] Alert on events: 
   - [ ] Proposal creation
   - [ ] Slashing
   - [ ] Huge delegation
   - [ ] Validator low uptime
   - [ ] Huge undelegation
   - [ ] Proposal start voting 
   - [ ] Proposal voting ends
- [x] Validators information update history
- [ ] Validators rating
   - [ ] Self-delegation
   - [ ] Uptime
   - [ ] Ever slashed
   - [ ] Gov participation
   - [ ] Community contributions
   - [ ] Number of delegators
