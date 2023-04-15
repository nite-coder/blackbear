package config

import (
	"errors"
	"fmt"
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

type ChangedEvent func() error

type Configuration struct {
	providers      []Provider
	rwMutex        sync.RWMutex
	eventChan      chan bool
	OnChangedEvent ChangedEvent
}

func new() *Configuration {
	c := &Configuration{
		providers: []Provider{},
		eventChan: make(chan bool, 1),
	}

	go func() {
		for range c.eventChan {
			_ = c.OnChangedEvent()
		}
	}()

	return c
}

func OnChangedEvent(event ChangedEvent) {
	cfg.OnChangedEvent = event
}

func AddProvider(provider Provider) {
	cfg.rwMutex.Lock()
	defer cfg.rwMutex.Unlock()

	cfg.providers = append(cfg.providers, provider)

	eventChan := cfg.eventChan
	go func() {
		evt := <-eventChan
		cfg.eventChan <- evt
	}()
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

	return "", fmt.Errorf("%w, key: %s", ErrKeyNotFound, key)
}

// StringMap returns a map[string]interface{} type value which has the key.
func StringMap(key string, defaultValue ...map[string]interface{}) (map[string]interface{}, error) {
	cfg.rwMutex.RLock()
	defer cfg.rwMutex.RUnlock()

	for _, p := range cfg.providers {
		val, err := p.Get(key)

		if err != nil {
			continue
		}

		return cast.ToStringMap(val)
	}

	if len(defaultValue) > 0 {
		return defaultValue[0], nil
	}

	return nil, fmt.Errorf("%w, key: %s", ErrKeyNotFound, key)
}

// StringMapString returns a map[string]string type value which has the key.
func StringMapString(key string, defaultValue ...map[string]string) (map[string]string, error) {
	cfg.rwMutex.RLock()
	defer cfg.rwMutex.RUnlock()

	for _, p := range cfg.providers {
		val, err := p.Get(key)

		if err != nil {
			continue
		}

		return cast.ToStringMapString(val)
	}

	if len(defaultValue) > 0 {
		return defaultValue[0], nil
	}

	return nil, fmt.Errorf("%w, key: %s", ErrKeyNotFound, key)
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

	return 0, fmt.Errorf("%w, key: %s", ErrKeyNotFound, key)
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

	return 0, fmt.Errorf("%w, key: %s", ErrKeyNotFound, key)
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

	return 0, fmt.Errorf("%w, key: %s", ErrKeyNotFound, key)
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

	return 0, fmt.Errorf("%w, key: %s", ErrKeyNotFound, key)
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

	return 0, fmt.Errorf("%w, key: %s", ErrKeyNotFound, key)
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

	return false, fmt.Errorf("%w, key: %s", ErrKeyNotFound, key)
}

// StringSlice returns the value associated with the key as a slice of strings.
func StringSlice(key string, defaultValue ...[]string) ([]string, error) {
	for _, p := range cfg.providers {
		val, err := p.Get(key)

		if err != nil {
			continue
		}

		return cast.ToStringSlice(val)
	}

	if len(defaultValue) > 0 {
		return defaultValue[0], nil
	}

	return []string{}, fmt.Errorf("%w, key: %s", ErrKeyNotFound, key)
}

// IntSlice returns the value associated with the key as a slice of int.
func IntSlice(key string, defaultValue ...[]int) ([]int, error) {
	for _, p := range cfg.providers {
		val, err := p.Get(key)

		if err != nil {
			continue
		}

		return cast.ToIntSlice(val)
	}

	if len(defaultValue) > 0 {
		return defaultValue[0], nil
	}

	return []int{}, fmt.Errorf("%w, key: %s", ErrKeyNotFound, key)
}

// Float64Slice returns the value associated with the key as a slice of float.
func Float64Slice(key string, defaultValue ...[]float64) ([]float64, error) {
	for _, p := range cfg.providers {
		val, err := p.Get(key)

		if err != nil {
			continue
		}

		return cast.ToFloat64Slice(val)
	}

	if len(defaultValue) > 0 {
		return defaultValue[0], nil
	}

	return []float64{}, fmt.Errorf("%w, key: %s", ErrKeyNotFound, key)
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

	return d, fmt.Errorf("%w, key: %s", ErrKeyNotFound, key)
}

// Scan binds a value which has the key.
func Scan(key string, value interface{}) error {
	for _, p := range cfg.providers {
		val, err := p.Get(key)

		if err != nil {
			continue
		}

		decoder, _ := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
			DecodeHook:       mapstructure.StringToTimeDurationHookFunc(),
			WeaklyTypedInput: true,
			Result:           value,
		})

		err = decoder.Decode(val)
		if err != nil {
			return err
		}
		return nil
	}

	return fmt.Errorf("%w, key: %s", ErrKeyNotFound, key)
}
