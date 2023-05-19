package v3

import (
	v3 "github.com/forbole/juno/v4/cmd/migrate/v3"
)

type Config struct {
	v3.Config `yaml:"-,inline"`
}
