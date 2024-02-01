package v5

import (
	"fmt"
	"os"
	"path"

	"github.com/forbole/juno/v5/types/config"
	"gopkg.in/yaml.v3"
)

// GetConfig returns the configuration reading it from the config.yaml file present inside the home directory
func GetConfig() (config.Config, error) {
	file := path.Join(config.HomePath, "config.yaml")

	// Make sure the path exists
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return config.Config{}, fmt.Errorf("config file does not exist")
	}

	bz, err := os.ReadFile(file)
	if err != nil {
		return config.Config{}, fmt.Errorf("error while reading config file: %s", err)
	}

	var cfg config.Config
	err = yaml.Unmarshal(bz, &cfg)
	return cfg, err
}
