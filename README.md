# BDJuno
[![GitHub Workflow Status](https://img.shields.io/github/workflow/status/forbole/bdjuno/Tests)](https://github.com/forbole/bdjuno/actions?query=workflow%3ATests)
[![Go Report Card](https://goreportcard.com/badge/github.com/forbole/bdjuno)](https://goreportcard.com/report/github.com/forbole/bdjuno)
[![Codecov](https://img.shields.io/codecov/c/github/forbole/bdjuno)](https://codecov.io/gh/forbole/bdjuno/branch/cosmos-v0.39.x)

BDJuno (shorthand for BigDipper Juno) is the [Juno](https://github.com/desmos-labs/juno) implementation for [BigDipper](https://github.com/forbole/big-dipper). 

It extends the custom Juno behavior by adding different handlers and custom operations to make it easier for BigDipper
showing the data inside the UI.

All the chains' data that are queried from the LCD and RPC endpoints are stored inside
a [PostgreSQL](https://www.postgresql.org/) database on top of which [GraphQL](https://graphql.org/) APIs can then be
created using [Hasura](https://hasura.io/).

## Installation
### Install the binaries
To install the binary simply run `make install`.

**Note**: Requires [Go 1.13+](https://golang.org/dl/).

### Database
Before running the parser, you need to:

1. Create a [PostgreSQL](https://www.postgresql.org/) database.
2. Run the SQL queries you find inside the `schema` folder inside such database to create all the necessary tables.

## Running the parser
To parse the chain state, you need to use the following command:

```shell
bdjuno parse <path/to/config.toml>

# Example
# bdjuno parse config.toml 
```

The configuration must be a TOML file containing the following fields:

```toml
[cosmos]
prefix = "desmos"
modules = []

[rpc]
address = "<rpc-ip/host>:<rpc-port>"

[grpc]
address = "<grpc-ip/host>:<grpc-port>"
insecure = true

[api]
address = "<client-ip/host>:<client-port>"

[cosmos]
prefix = "desmos"
modules = []

[database]
type = "postgresql"

[database.config]
host = "<db-host>"
port = 5432
name = "<db-name>"
user = "<db-user>"
password = "<db-password>"
```

Example of a configuration to parse the chain state from a local Cosmos full-node:

<details>

<summary>Open here</summary>

```toml
[cosmos]
prefix = "desmos"
modules = []

[rpc]
address = "http://localhost:26657"

[grpc]
address = "localhost:9090"
insecure = true

[api]
address = "http://localhost:1317"

[cosmos]
prefix = "desmos"
modules = [
    "auth",
    "bank",
    "consensus",
    "distribution",
    "gov",
    "mint",
    "modules",
    "pricefeed",
    "slashing",
    "staking"
]

[database]
type = "postgresql"

[database.config]
host = "localhost"
port = 5432
name = "bdjuno"
user = "user"
password = "password"
```

</details>

## Testing
If you want to test the code, you can do so by running

```sh
make test-unit
```

**Note**: Requires [Docker](https://docker.com).

This will:
1. Create a Docker container running a PostgreSQL database.
2. Run all the tests using that database as support.


