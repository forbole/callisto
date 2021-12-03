module github.com/forbole/bdjuno/v2

go 1.16

require (
	github.com/cosmos/cosmos-sdk v0.44.3
	github.com/cosmos/gaia/v6 v6.0.0-rc1
	github.com/forbole/juno/v2 v2.0.0-20211020184842-e358a33007ff
	github.com/go-co-op/gocron v1.10.0
	github.com/gogo/protobuf v1.3.3
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/jmoiron/sqlx v1.2.1-0.20200324155115-ee514944af4b
	github.com/lib/pq v1.10.4
	github.com/pelletier/go-toml v1.9.4
	github.com/proullon/ramsql v0.0.0-20181213202341-817cee58a244
	github.com/regen-network/regen-ledger/v2 v2.1.0
	github.com/rs/zerolog v1.26.0
	github.com/spf13/cobra v1.2.1
	github.com/stretchr/testify v1.7.0
	github.com/tendermint/tendermint v0.34.14
	google.golang.org/grpc v1.40.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1

replace google.golang.org/grpc => google.golang.org/grpc v1.33.2

replace github.com/tendermint/tendermint => github.com/forbole/tendermint v0.34.13-0.20210820072129-a2a4af55563d

replace github.com/cosmos/cosmos-sdk => github.com/regen-network/cosmos-sdk v0.44.2-regen-1

replace github.com/regen-network/regen-ledger/types => github.com/regen-network/regen-ledger/types v1.0.0

replace github.com/regen-network/regen-ledger/orm => github.com/regen-network/regen-ledger/orm v1.0.0-beta1

replace github.com/regen-network/regen-ledger/x/data => github.com/regen-network/regen-ledger/x/data v0.0.0-20211202162855-98bb2b15c314

replace github.com/regen-network/regen-ledger/x/ecocredit => github.com/regen-network/regen-ledger/x/ecocredit v1.0.0

replace github.com/regen-network/regen-ledger/x/group => github.com/regen-network/regen-ledger/x/group v1.0.0-beta1
