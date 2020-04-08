# Requirements of BDjuno

The BDjuno is working as a backend for Big Dipper. The followings are the features currently supported in Big Dipper in the backend which BDjuno should adapt.

## On every block

1. Parsing all blocks
2. Parsing all transactions
3. Update Validator information (staking) and calculate voting power percentage and self delegation ratio
4. Store validator set of the block
5. Update miss block records
6. Calculate the average block time
7. Read the latest consensus state
8. Update the total supply (supply)
9. Update the total staked tokens (staking)
10. Update the inflation (mint)
11. Update community pool (distribution)

## On intervals

1. Calculate average block time (per minute, per hour, per day)
2. Calculate average delagtion ratio (per hour, per day)
3. Calculate voting power distribution (per hour)
4. Record all delegations (per day)
5. Record all undelegatios (per day)
6. Record all redelegations (per day)
7. Get token price and marketcap (per 30 seconds)
8. Get gov proposals (per 30 seconds, gov)
9. Calculate the tally result (per 30 seconds, gov)

## Not on Big Dipper now but considering to add

1. Validators signing-info (slashing)
2. All wallets activities
3. Alert on events (proposal creation, slashing, hugh delegation, validator low uptime, huge undelegation, proposal start voting, proposal voting ends)
4. Validators information update history
5. Valdiators rating (self-delegation, uptime, ever slashed, gov participation, community contributions, number of delegators)
