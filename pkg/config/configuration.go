package config

import (
	"errors"
	"sync"

	"github.com/mitchellh/mapstructure"
	"github.com/nite-coder/blackbear/pkg/cast"
)

var (
	cfg                     = new()
	ErrFileNotFound         = errors.New("config: config file was not found")
	ErrKeyNotFound          = errors.New("config: key was not found")
	ErrProviderNotFound     = errors.New("config: no provider is added to config.  Provider need to be added first")
	ErrConfigTypeNotSupport = errors.New("config: type is not support")
)

type Configuration struct {
	providers []Provider
	rwMutex   sync.RWMutex
}

func new() *Configuration {
	return &Configuration{
		providers: []Provider{},
	}
}

func AddProvider(provider Provider) {
	cfg.rwMutex.Lock()
	defer cfg.rwMutex.Unlock()

	cfg.providers = append(cfg.providers, provider)
}

func RemoveAllPrividers() {
	cfg.rwMutex.Lock()
	defer cfg.rwMutex.Unlock()

	cfg.providers = []Provider{}
}

// String returns a string type value which has the key.
func String(key string, defaultValue ...string) (string, error) {
	if len(cfg.providers) == 0 {
		return "", ErrProviderNotFound
	}

	for _, p := range cfg.providers {
		val, err := p.Get(key)

		if err != nil {
			continue
		}

		return cast.ToString(val)
	}

	if len(defaultValue) > 0 {
		return defaultValue[0], nil
	}

	return "", ErrKeyNotFound
}

// Int32 returns a int32 type value which has the key.
func Int32(key string, defaultValue ...int32) (int32, error) {
	if len(cfg.providers) == 0 {
		return 0, ErrProviderNotFound
	}

	for _, p := range cfg.providers {
		val, err := p.Get(key)

		if err != nil {
			continue
		}

		return cast.ToInt32(val)
	}

	if len(defaultValue) > 0 {
		return defaultValue[0], nil
	}

	return 0, ErrKeyNotFound
}

// UnmarshalKey binds a value which has the key.
func UnmarshalKey(key string, value interface{}) error {
	if len(cfg.providers) == 0 {
		return ErrProviderNotFound
	}

	for _, p := range cfg.providers {
		val, err := p.Get(key)

		if err != nil {
			continue
		}

		err = mapstructure.Decode(val, value)
		if err != nil {
			return err
		}
		return nil
	}

	return ErrKeyNotFound
}
