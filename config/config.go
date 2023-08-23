package config

import (
	"git.ooo.ua/vipcoin/lib/vault"
	junoconf "github.com/forbole/juno/v3/types/config"
	"github.com/spf13/cobra"
)

// CheckVaultConfig - checking the ability to load the configuration from vault
func CheckVaultConfig(serviceName string, cmdAndConf *cobra.Command) *cobra.Command {
	config, err := loadFromVault(serviceName)
	if err != nil {
		return cmdAndConf
	}

	cmdAndConf.PreRunE = func(_ *cobra.Command, _ []string) error {
		// Set the global configuration
		junoconf.Cfg = config
		return nil
	}

	return cmdAndConf
}

func loadFromVault(serviceName string) (junoconf.Config, error) {
	vaultClient, err := vault.NewClient(vault.NamespaceCubbyhole)
	if err != nil {
		return junoconf.Config{}, err
	}

	vaultData, err := vaultClient.Pull(serviceName)
	if err != nil {
		return junoconf.Config{}, err
	}

	return junoconf.DefaultConfigParser(vaultData)
}
