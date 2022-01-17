module github.com/forbole/bdjuno/v2

go 1.16

require (
	github.com/cosmos/cosmos-sdk v0.44.4
	github.com/cosmos/gaia/v6 v6.0.0-rc1
	github.com/cosmos/ibc-go v1.2.3
	github.com/desmos-labs/desmos/v2 v2.3.1
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

replace github.com/tendermint/tendermint => github.com/forbole/tendermint v0.34.13-0.20210820072129-a2a4af55563d

replace github.com/cosmos/cosmos-sdk => github.com/desmos-labs/cosmos-sdk v0.43.0-alpha1.0.20211102084520-683147efd235

replace github.com/99designs/keyring => github.com/cosmos/keyring v1.1.7-0.20210622111912-ef00f8ac3d76

replace github.com/cosmos/ledger-cosmos-go => github.com/desmos-labs/ledger-desmos-go v0.11.2-0.20210814121638-5d87e392e8a9
