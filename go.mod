module github.com/forbole/bdjuno/v2

go 1.16

require (
	github.com/cosmos/cosmos-sdk v0.42.9
	github.com/forbole/juno/v2 v2.0.0-20211018063058-c8a1eb5d5312
	github.com/go-co-op/gocron v0.3.3
	github.com/gogo/protobuf v1.3.3
	github.com/jmoiron/sqlx v1.2.1-0.20200324155115-ee514944af4b
	github.com/lib/pq v1.10.3
	github.com/pelletier/go-toml v1.9.4
	github.com/proullon/ramsql v0.0.0-20181213202341-817cee58a244
	github.com/rs/zerolog v1.25.0
	github.com/spf13/cobra v1.2.1
	github.com/stretchr/testify v1.7.0
	github.com/tendermint/tendermint v0.34.12
	google.golang.org/grpc v1.38.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1

replace google.golang.org/grpc => google.golang.org/grpc v1.33.2

replace github.com/tendermint/tendermint => github.com/forbole/tendermint v0.34.13-0.20210820072129-a2a4af55563d
