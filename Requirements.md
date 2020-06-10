# Requirements of BDjuno
The BDjuno is working as a backend for Big Dipper. The followings are the features currently supported in Big Dipper in the backend which BDjuno should adapt.

## On every block
### Done by Juno by default
- [x] Parsing all blocks
- [x] Parsing all transactions

### Custom BDJuno implementations
- [ ] Update validator information (staking) 
- [ ] Calculate validator voting power percentage
- [ ] Calculate validator self delegation ratio
- [ ] Store validator set of the block
- [ ] Update miss block records
- [ ] Read the latest consensus state
- [ ] Update the total supply (supply)
- [ ] Update the total staked tokens (staking)
- [ ] Update the inflation (mint)
- [ ] Update community pool (distribution)

### Achievable using GraphQL APIs
- [x] Calculate the average block time

## On intervals
- [ ] Calculate average block time (per minute, per hour, per day)
- [ ] Calculate average delagtion ratio (per hour, per day)
- [ ] Calculate voting power distribution (per hour)
- [ ] Record all delegations (per day)
- [ ] Record all undelegatios (per day)
- [ ] Record all redelegations (per day)
- [ ] Get token price and marketcap (per 30 seconds)
- [ ] Get gov proposals (per 30 seconds, gov)
- [ ] Calculate the tally result (per 30 seconds, gov)

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
