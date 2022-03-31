package remote

import (
	"github.com/forbole/juno/v2/node/remote"

	"github.com/forbole/bdjuno/v2/modules/vipcoin/chain/accounts/source"

	accountstypes "git.ooo.ua/vipcoin/chain/x/accounts/types"
)

var (
	_ source.Source = &Source{}
)

type Source struct {
	*remote.Source
	accounts accountstypes.QueryClient
}

// NewSource builds a new Source instance
func NewSource(source *remote.Source, accClient accountstypes.QueryClient) *Source {
	return &Source{
		Source:   source,
		accounts: accClient,
	}
}

// GetAccounts implements bankkeeper.Source
func (s Source) GetAccounts(addresses []string, height int64) ([]*accountstypes.Account, error) {
	return []*accountstypes.Account{}, nil
}
