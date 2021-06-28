package env

import (
	"os"
	"strings"

	"github.com/nite-coder/blackbear/pkg/config"
)

type EnvProvider struct {
	envPrefix string
}

func New() *EnvProvider {
	return &EnvProvider{}
}

func (p *EnvProvider) SetEnvPrefix(prefix string) {
	p.envPrefix = prefix
}

func (p *EnvProvider) Get(key string) (interface{}, error) {
	key = strings.ReplaceAll(key, ".", "_")

	if len(p.envPrefix) > 0 {
		key = p.envPrefix + "_" + key
	}

	key = strings.ToUpper(key)
	val, present := os.LookupEnv(key)

	if present {
		return val, nil
	}

	return "", config.ErrKeyNotFound
}
