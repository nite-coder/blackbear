package config

import (
	"errors"
	"sync"
	"time"

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

// Int returns a int type value which has the key.
func Int(key string, defaultValue ...int) (int, error) {
	for _, p := range cfg.providers {
		val, err := p.Get(key)

		if err != nil {
			continue
		}

		return cast.ToInt(val)
	}

	if len(defaultValue) > 0 {
		return defaultValue[0], nil
	}

	return 0, ErrKeyNotFound
}

// Int32 returns a int32 type value which has the key.
func Int32(key string, defaultValue ...int32) (int32, error) {
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

// Int64 returns a int64 type value which has the key.
func Int64(key string, defaultValue ...int64) (int64, error) {
	for _, p := range cfg.providers {
		val, err := p.Get(key)

		if err != nil {
			continue
		}

		return cast.ToInt64(val)
	}

	if len(defaultValue) > 0 {
		return defaultValue[0], nil
	}

	return 0, ErrKeyNotFound
}

// Float32 returns a float32 type value which has the key.
func Float32(key string, defaultValue ...float32) (float32, error) {
	for _, p := range cfg.providers {
		val, err := p.Get(key)

		if err != nil {
			continue
		}

		return cast.ToFloat32(val)
	}

	if len(defaultValue) > 0 {
		return defaultValue[0], nil
	}

	return 0, ErrKeyNotFound
}

// Float64 returns a float64 type value which has the key.
func Float64(key string, defaultValue ...float64) (float64, error) {
	for _, p := range cfg.providers {
		val, err := p.Get(key)

		if err != nil {
			continue
		}

		return cast.ToFloat64(val)
	}

	if len(defaultValue) > 0 {
		return defaultValue[0], nil
	}

	return 0, ErrKeyNotFound
}

// Bool returns a boolean type value which has the key.
func Bool(key string, defaultValue ...bool) (bool, error) {
	for _, p := range cfg.providers {
		val, err := p.Get(key)

		if err != nil {
			continue
		}

		return cast.ToBool(val)
	}

	if len(defaultValue) > 0 {
		return defaultValue[0], nil
	}

	return false, ErrKeyNotFound
}

// Bool returns a boolean type value which has the key.
func Duration(key string, defaultValue ...time.Duration) (time.Duration, error) {
	for _, p := range cfg.providers {
		val, err := p.Get(key)

		if err != nil {
			continue
		}

		return cast.ToDuration(val)
	}

	var d time.Duration

	if len(defaultValue) > 0 {
		return defaultValue[0], nil
	}

	return d, ErrKeyNotFound
}

// Scan binds a value which has the key.
func Scan(key string, value interface{}) error {
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
