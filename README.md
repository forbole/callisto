# BDJuno
[![GitHub Workflow Status](https://img.shields.io/github/workflow/status/forbole/bdjuno/Tests)](https://github.com/forbole/bdjuno/actions?query=workflow%3ATests)
[![Go Report Card](https://goreportcard.com/badge/github.com/forbole/bdjuno)](https://goreportcard.com/report/github.com/forbole/bdjuno)
![Codecov branch](https://img.shields.io/codecov/c/github/forbole/bdjuno/cosmos/v0.40.x)

BDJuno (shorthand for BigDipper Juno) is the [Juno](https://github.com/forbole/juno) implementation
for [BigDipper](https://github.com/forbole/big-dipper).

It extends the custom Juno behavior by adding different handlers and custom operations to make it easier for BigDipper
showing the data inside the UI.

All the chains' data that are queried from the RPC and gRPC endpoints are stored inside
a [PostgreSQL](https://www.postgresql.org/) database on top of which [GraphQL](https://graphql.org/) APIs can then be
created using [Hasura](https://hasura.io/).

## Usage
To know how to setup and run BDJuno, please refer to
the [docs website](https://docs.bigdipper.live/cosmos-based/parser/overview/).

## Local env via Docker
1. Setup a local postgres DB
2. Execute examples/big_dipper_2_init_script_combined.sql on a new DB
3. Run a local cudos-node on localhost
4. Rename bdjuno_sample_configs and the files inside it to bdjuno/config.yaml and bdjuno/genesis.json
5. Fix the configs to match your env ( https://docs.bigdipper.live/cosmos-based/parser/config/config )
6. Rename .env-dev.sample to .env and configure it as per your needs (default is enough, just replace DB connection info)
7. Execute ```docker-compose up -d```
8. cd /path/to/BDJuno/hasura
9. Run ```curl -L https://github.com/hasura/graphql-engine/raw/stable/cli/get.sh | bash```
10. Run ```hasura metadata apply --endpoint http://localhost:8080 --admin-secret <hasura_password>``` 
11. BDJuno should start parsing info from your local node and feed it in the specified postgres DB

## Testing
If you want to test the code, you can do so by running

```shell
$ make test-unit
```

**Note**: Requires [Docker](https://docker.com).

This will:
1. Create a Docker container running a PostgreSQL database.
2. Run all the tests using that database as support.


