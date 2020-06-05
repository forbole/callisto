# Requirements of BDjuno
The BDjuno is working as a backend for Big Dipper. The followings are the features currently supported in Big Dipper in the backend which BDjuno should adapt.

## On every block
- [x] Parsing all blocks
- [x] Parsing all transactions
- [ ] Update Validator information (staking) and calculate voting power percentage and self delegation ratio
- [ ] Store validator set of the block
- [ ] Update miss block records
- [ ] Calculate the average block time
- [ ] Read the latest consensus state
- [ ] Update the total supply (supply)
- [ ] Update the total staked tokens (staking)
- [ ] Update the inflation (mint)
- [ ] Update community pool (distribution)

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

## Not on Big Dipper now but considering to add
- [ ] Validators signing-info (slashing)
- [ ] All wallets activities
- [ ] Alert on events (proposal creation, slashing, hugh delegation, validator low uptime, huge undelegation, proposal start voting, proposal voting ends)
- [ ] Validators information update history
- [ ] Valdiators rating (self-delegation, uptime, ever slashed, gov participation, community contributions, number of delegators)
