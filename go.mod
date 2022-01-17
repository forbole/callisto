module github.com/forbole/bdjuno/v2

go 1.16

require (
	github.com/MonikaCat/ag0/v6 v6.0.0-20211026142553-1cc79476b438
	github.com/cosmos/cosmos-sdk v0.44.3
	github.com/cosmos/gaia/v6 v6.0.0-rc1
	github.com/cosmos/ibc-go v1.2.0
	github.com/forbole/juno/v2 v2.0.0-20220117075513-d927d34156a9
	github.com/go-co-op/gocron v1.11.0
	github.com/gogo/protobuf v1.3.3
	github.com/jmoiron/sqlx v1.2.1-0.20200324155115-ee514944af4b
	github.com/lib/pq v1.10.4
	github.com/pelletier/go-toml v1.9.4
	github.com/proullon/ramsql v0.0.0-20181213202341-817cee58a244
	github.com/rs/zerolog v1.26.1
	github.com/spf13/cobra v1.3.0
	github.com/stretchr/testify v1.7.0
	github.com/tendermint/tendermint v0.34.14
	google.golang.org/grpc v1.42.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1

replace google.golang.org/grpc => google.golang.org/grpc v1.33.2

replace github.com/tendermint/tendermint => github.com/agoric-labs/tendermint v0.34.13-alpha.agoric.7

replace github.com/cosmos/cosmos-sdk => github.com/agoric-labs/cosmos-sdk v0.44.2-alpha.agoric.1

replace github.com/keybase/go-keychain => github.com/99designs/go-keychain v0.0.0-20191008050251-8e49817e8af4

replace github.com/99designs/keyring => github.com/cosmos/keyring v1.1.7-0.20210622111912-ef00f8ac3d76
