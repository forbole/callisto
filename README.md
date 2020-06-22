# BDJuno
[![Build Status](https://travis-ci.com/forbole/bdjuno.svg?branch=master)](https://travis-ci.com/forbole/bdjuno)
[![Go Report Card](https://goreportcard.com/badge/github.com/forbole/bdjuno)](https://goreportcard.com/reporl-debugt/github.com/forbole/bdjuno)
[![Codecov](https://img.shields.io/codecov/c/github/forbole/bdjuno)](https://codecov.io/gh/forbole/bdjuno/branch/master)


BDJuno (shorthand for BigDipper Juno) is the [Juno](https://github.com/desmos-labs/juno) implementation for [BigDipper](https://github.com/forbole/big-dipper). 

It extends the custom Juno behavior by adding different handlers and custom operations to make it easier for BigDipper showing the data inside the UI. 

All the chains' data that are queried from the LCD and RPC endpoints are stored inside a [PostgreSQL](https://www.postgresql.org/) database on top of which [GraphQL](https://graphql.org/) APIs can then be created using [Hasura](https://hasura.io/). 
 
## Installation
To install the binary simply run `make install`.

**Note**: Requires [Go 1.13+](https://golang.org/dl/)

## Database
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
rpc_node = "<rpc-ip/host>:<rpc-port>"
client_node = "<client-ip/host>:<client-port>"

[database]
host = "<db-host>"
port = <db-port>
name = "<db-name>"
user = "<db-user>"
password = "<db-password>"
ssl_mode = "<ssl-mode>"
```

Example of a configuration to parse the chain state from a local full-node: 

```
rpc_node = "http://localhost:26657"
client_node = "http://localhost:1317"

[database]
host = "localhost"
port = 5432
name = "big-dipper"
user = "<username>"
password = "<password>"
```
