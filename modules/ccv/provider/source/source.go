package source

type Source interface {
	GetValidatorProviderAddr(height int64, chainID, consumerAddress string) (string, error)
}
