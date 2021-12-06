module github.com/forbole/bdjuno/v2

go 1.16

require (
	github.com/cosmos/cosmos-sdk v0.44.3
	github.com/cosmos/gaia/v6 v6.0.0-rc1
	github.com/forbole/juno/v2 v2.0.0-20211122103136-7926db0202f2
	github.com/go-co-op/gocron v0.3.3
	github.com/gogo/protobuf v1.3.3
	github.com/jmoiron/sqlx v1.2.0
	github.com/lib/pq v1.10.2
	github.com/pelletier/go-toml v1.9.4
	github.com/proullon/ramsql v0.0.0-20211120092837-c8d0a408b939
	github.com/rs/zerolog v1.23.0
	github.com/spf13/cobra v1.2.1
	github.com/stretchr/testify v1.7.0
	github.com/tendermint/tendermint v0.34.14
	google.golang.org/genproto v0.0.0-20210920155426-26f343e4c215 // indirect
	google.golang.org/grpc v1.40.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)

require (
	github.com/fsnotify/fsnotify v1.5.1 // indirect
	github.com/mitchellh/mapstructure v1.4.2 // indirect
	github.com/spf13/cast v1.4.1 // indirect
	golang.org/x/crypto v0.0.0-20210817164053-32db794688a5 // indirect
	gopkg.in/ini.v1 v1.63.2 // indirect
)

replace google.golang.org/grpc => google.golang.org/grpc v1.33.2

replace github.com/tendermint/tendermint => github.com/forbole/tendermint v0.34.13-0.20210820072129-a2a4af55563d

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.4

replace github.com/cosmos/cosmos-sdk => github.com/regen-network/cosmos-sdk v0.44.2-regen-1

replace github.com/regen-network/regen-ledger/x/data => github.com/regen-network/regen-ledger/x/data v0.0.0-20210602121340-fa967f821a6e

replace github.com/regen-network/regen-ledger/types => github.com/regen-network/regen-ledger/types v1.0.0
