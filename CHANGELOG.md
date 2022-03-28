## Version v2.0.0


### Notes
This version introduces breaking changes to certain address-specific data that is no longer periodically parsed from the node and stored in the database. Instead the data is now obtained directly from the node when needed using Hasura Actions. Read more details about [migrating to v2.0.0](https://docs.bigdipper.live/cosmos-based/parser/migrations/v2.0.0)```

### New features
#### CLI
- Added `parse-genesis` command to parse the genesis file
- Added `fix` command:
  - auth: fix vesting-accounts details
  - blocks: fix missing blocks and transactions from the configured start height
  - gov: fix proposal with proposal ID specified  
  - staking: fix validators info at the latest height  

#### Hasura Actions
- Replaced some periodic queries with hasura actions 
- Here's a list of data acquired through Hasura Actions:
  - Of a certain address/delegator:
    - Account balance
    - Delegation rewards
    - Delegator withdraw address
    - Delegations
    - Total delegations amount
    - Unbonding delegations
    - Total unbonding delegations amount
    - Redelegations
  - Of a certain validator:
    - Commission amount
    - Delegations to this validator
    - Redelegations from this validator
    - Unbonding delegations
  - Note: graphQL queries on the frontend should be updated for the above info
- Added prometheus monitoring to hasura actions

#### Node Type Local
- Added note.type=local for parsing a static local node without gRPC query:
[config reference](https://docs.bigdipper.live/cosmos-based/parser/config/config#node)


#### Modules
- auth: `vesting accounts` and `vesting periods` are handled and stored in the database 


### Changes 

#### Database
- `bonded_tokens` and `not_bonded_tokens` db types in `staking_pool` table are changed to `TEXT` to avoid digits overflow
- `tombstone` status is accessible from `validator_status` table
