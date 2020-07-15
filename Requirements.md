# Requirements of BDjuno
The BDjuno is working as a backend for Big Dipper. The followings are the features currently supported in Big Dipper in the backend which BDjuno should adapt.

## On every block
### Done by Juno by default
- [x] Parsing all blocks
- [x] Parsing all transactions
- [x] Store validator set of the block

### Custom BDJuno implementations
- [x] Update miss block records
- [ ] Read the latest consensus state
- [x] [x/staking] Update validator information 
- [ ] [x/staking] Calculate validator voting power percentage 
- [x] [x/staking] Calculate validator self delegation ratio
- [x] [x/staking] Update the total staked tokens 
- [ ] [x/supply] Update the total supply
- [x] [x/mint] Update the inflation
- [ ] [x/distribution] Update community pool

### Achievable using GraphQL APIs
- [x] Calculate the average block time

## On intervals
- [ ] Calculate average block time (per minute, per hour, per day)
- [ ] Get token price and marketcap (per 30 seconds)
- [ ] [x/staking] Calculate average delegation ratio (per hour, per day)
- [ ] [x/staking] Calculate voting power distribution (per hour)
- [ ] [x/staking] Record all delegations (per day)
- [ ] [x/staking] Record all undelegatios (per day)
- [ ] [x/staking] Record all redelegations (per day)
- [ ] [x/gov] Get gov proposals (per 30 seconds)
- [ ] [x/gox] Calculate the tally result (per 30 seconds)

## Not on Big Dipper now but we are considering to add
- [ ] Validators signing-info (slashing)
- [ ] All wallets activities
- [ ] Alert on events: 
   - [ ] Proposal creation
   - [ ] Slashing
   - [ ] Huge delegation
   - [ ] Validator low uptime
   - [ ] Huge undelegation
   - [ ] Proposal start voting 
   - [ ] Proposal voting ends
- [ ] Validators information update history
- [ ] Validators rating
   - [ ] Self-delegation
   - [ ] Uptime
   - [ ] Ever slashed
   - [ ] Gov participation
   - [ ] Community contributions
   - [ ] Number of delegators
