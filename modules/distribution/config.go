package distribution

// DistributionConfig contains the configuration about distribution frequency
type DistributionConfig struct {
	DistributionFrequency int64 `toml:"distribution_frequency"`
}

// GetDistributionFrequency returns distribution frequency int64 value
func (b *DistributionConfig) GetDistributionFrequency() int64 {
	return b.DistributionFrequency
}

