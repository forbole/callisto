package consumer

type ProviderModule interface {
	GetValidatorProviderAddr(height int64, chainID, consumerAddress string) (string, error)
}
