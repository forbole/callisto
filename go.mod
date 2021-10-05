module github.com/forbole/bdjuno

go 1.15

require (
	github.com/cosmos/cosmos-sdk v0.44.1
	github.com/desmos-labs/desmos/v2 v2.0.1
	github.com/desmos-labs/juno v0.0.0-20211005090705-514187767199
	github.com/go-co-op/gocron v0.3.3
	github.com/gogo/protobuf v1.3.3
	github.com/jmoiron/sqlx v1.2.1-0.20200324155115-ee514944af4b
	github.com/lib/pq v1.10.2
	github.com/pelletier/go-toml v1.9.3
	github.com/proullon/ramsql v0.0.0-20181213202341-817cee58a244
	github.com/rs/zerolog v1.23.0
	github.com/stretchr/testify v1.7.0
	github.com/tendermint/tendermint v0.34.13
	github.com/ziutek/mymysql v1.5.4 // indirect
	google.golang.org/grpc v1.40.0
)

replace google.golang.org/grpc => google.golang.org/grpc v1.33.2

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1

replace github.com/tendermint/tendermint => github.com/forbole/tendermint v0.34.13-0.20211005074050-33591287eca5
