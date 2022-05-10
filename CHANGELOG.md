## Version v3.0.1
### Dependencies
- Updated Juno to `v3.1.1`

## Version v3.0.0
### Notes
This version introduces breaking changes to `transaction` and `message` PostgreSQL tables. It implements PostgreSQL table partitioning to fix slow data retrieval from database that stores large amount of transactions and messages. Read more details about [migrating to v3.0.0](https://docs.bigdipper.live/cosmos-based/parser/migrations/v2.0.0)

### New features 
#### CLI
- ([\#356](https://github.com/forbole/bdjuno/pull/356)) Implemented `migrate` command to perform easy migration to higher BDJuno versions
- ([\#356](https://github.com/forbole/bdjuno/pull/356)) Updated `parse-genesis` command to parse genesis file without accessing the node

#### Database
- ([\#356](https://github.com/forbole/bdjuno/pull/356)) Added PostgreSQL table partition to `transaction` and `message` table
- ([\#356](https://github.com/forbole/bdjuno/pull/356)) Created new `messages_by_address` function

### Changes 
#### Hasura
- ([\#377](https://github.com/forbole/bdjuno/pull/377)) Updated Hasura metadata
- ([\#381](https://github.com/forbole/bdjuno/pull/381)) Hasura actions are now a module 

### Dependencies
- ([\#356](https://github.com/forbole/bdjuno/pull/356)) Updated Juno to `v3.0.0`

## Version v2.0.0
### Notes
This version introduces breaking changes to certain address-specific data that is no longer periodically parsed from the node and stored in the database. Instead, the data is now obtained directly from the node when needed using Hasura Actions. Read more details about [migrating to v2.0.0](https://docs.bigdipper.live/cosmos-based/parser/migrations/v2.0.0)

### New features
#### CLI
- ([\#257](https://github.com/forbole/bdjuno/pull/257)) Added `parse-genesis` command to parse the genesis file
- ([\#228](https://github.com/forbole/bdjuno/pull/228)) ([\#248](https://github.com/forbole/bdjuno/pull/248)) Added `fix` command:
  - `auth`: fix vesting accounts details
  - `blocks`: fix missing blocks and transactions from given start height
  - `gov`: fix proposal with given proposal ID  
  - `staking`: fix validators info at the latest height  

#### Hasura Actions
- ([\#329](https://github.com/forbole/bdjuno/pull/329)) Implemented Hasura Actions service to replace periodic queries. If you are using GraphQL queries on your application, you should updated the old queries to use the below new actions instead.  
  Here's a list of data acquired through Hasura Actions:
    - Of a certain address/delegator:
      - Account balance (`action_account_balance`)
      - Delegation rewards (`action_delegation_reward`)
      - Delegator withdraw address (`action_delegator_withdraw_address`)
      - Delegations (`action_delegation`)
      - Total delegations amount (`action_delegation_total`)
      - Unbonding delegations (`action_unbonding_delegation`)
      - Total unbonding delegations amount (`action_unbonding_delegation_total`)
      - Redelegations (`action_redelegation`)
    - Of a certain validator:
      - Commission amount (`action_validator_commission_amount`)
      - Delegations to this validator (`action_validator_delegations`)
      - Redelegations from this validator (`action_validator_redelegations_from`)
      - Unbonding delegations (`action_validator_unbonding_delegations`)
- ([\#352](https://github.com/forbole/bdjuno/pull/352)) Added prometheus monitoring to hasura actions

#### Local node support
- Added the support for `node.type = "local"` for parsing a static local node without the usage gRPC queries: [config reference](https://docs.bigdipper.live/cosmos-based/parser/config/config#node).

#### Modules
- ([\#232](https://github.com/forbole/bdjuno/pull/232)) Updated the `x/auth` module support to handle and store `vesting accounts` and `vesting periods` inside the database. 
- ([\#276](https://github.com/forbole/bdjuno/pull/276)) Added the support for the `x/feegrant` module (v0.44.x)

### Changes 

#### CLI
- ([\#351](https://github.com/forbole/bdjuno/pull/351)) Fixed version display for `bdjuno version` cmd 

#### Database
- ([\#300](https://github.com/forbole/bdjuno/pull/300)) Changed `bonded_tokens` and `not_bonded_tokens` type inside `staking_pool` table  to `TEXT` to avoid value overflow
- ([\#275](https://github.com/forbole/bdjuno/pull/275)) Added `tombstoned` column inside `validator_status` table
- ([\#232](https://github.com/forbole/bdjuno/pull/232)) Added `vesting_account` and `vesting_period` table
- ([\#276](https://github.com/forbole/bdjuno/pull/276)) Added `fee_grant_allowance` table (v0.44.x)

#### Modules
- ([\#353](https://github.com/forbole/bdjuno/pull/353)) Removed the support for the `history` module