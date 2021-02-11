module github.com/forbole/bdjuno

go 1.15

require (
	github.com/cosmos/cosmos-sdk v0.40.1
	github.com/desmos-labs/juno v0.0.0-20210211081720-48281291b206
	github.com/go-co-op/gocron v0.3.3
	github.com/jmoiron/sqlx v1.2.1-0.20200324155115-ee514944af4b
	github.com/lib/pq v1.9.0
	github.com/proullon/ramsql v0.0.0-20181213202341-817cee58a244
	github.com/rs/zerolog v1.20.0
	github.com/stretchr/testify v1.7.0
	github.com/tendermint/tendermint v0.34.3
	github.com/ziutek/mymysql v1.5.4 // indirect
	google.golang.org/grpc v1.35.0
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776 // indirect
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.4
