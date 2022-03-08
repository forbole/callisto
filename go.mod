module github.com/forbole/bdjuno/v2

go 1.16

require (
	github.com/cosmos/cosmos-sdk v0.44.5
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/forbole/juno/v2 v2.0.0-20220223115732-dbb226a91ce9
	github.com/go-co-op/gocron v1.11.0
	github.com/gogo/protobuf v1.3.3
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/hashicorp/golang-lru v0.5.5-0.20210104140557-80c98217689d // indirect
	github.com/hashicorp/hcl v1.0.1-0.20191016231534-914dc3f8dd7c // indirect
	github.com/jmhodges/levigo v1.0.1-0.20191019112844-b572e7f4cdac // indirect
	github.com/jmoiron/sqlx v1.2.1-0.20200324155115-ee514944af4b
	github.com/lib/pq v1.10.4
	github.com/libp2p/go-buffer-pool v0.0.3-0.20190619091711-d94255cb3dfc // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.2-0.20181231171920-c182affec369 // indirect
	github.com/pelletier/go-toml v1.9.4
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/proullon/ramsql v0.0.0-20181213202341-817cee58a244
	github.com/rs/cors v1.7.1-0.20191011001009-dcbccb712443 // indirect
	github.com/rs/zerolog v1.26.1
	github.com/spf13/cobra v1.3.0
	github.com/stretchr/objx v0.2.1-0.20190415111823-35313a95ee26 // indirect
	github.com/stretchr/testify v1.7.0
	github.com/subosito/gotenv v1.2.1-0.20190917103637-de67a6614a4d // indirect
	github.com/tendermint/tendermint v0.34.14
	golang.org/x/term v0.0.0-20210220032956-6a3ed077a48d // indirect
	google.golang.org/grpc v1.42.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1

replace google.golang.org/grpc => google.golang.org/grpc v1.33.2

replace github.com/tendermint/tendermint => github.com/forbole/tendermint v0.34.13-0.20210820072129-a2a4af55563d

replace github.com/99designs/keyring => github.com/cosmos/keyring v1.1.7-0.20210622111912-ef00f8ac3d76

replace github.com/cosmos/ledger-cosmos-go => github.com/desmos-labs/ledger-desmos-go v0.11.2-0.20210814121638-5d87e392e8a9

replace github.com/cosmos/cosmos-sdk => github.com/ovrclk/cosmos-sdk v0.44.5-patches
