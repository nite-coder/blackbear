package config

import (
	"github.com/mitchellh/mapstructure"
	"github.com/nite-coder/blackbear/pkg/cast"
)

var (
	cfg = new()
)

// Cfg return a singleton configuration instance.
func Cfg() Configuration {
	return cfg
}

// Load initialize this package. It will load config into cache and get ready to work.  However,
// if the config file was not found, `ErrFileNotFound` will be returned
func Load() error {
	return cfg.Load()
}

// LoadContent reads the val as content.
func LoadContent(content string) error {
	return cfg.LoadContent(content)
}

// ConfigName return config file name.  The default config file name is "app.yml"
func ConfigName() string {
	return cfg.ConfigName()
}

// SetConfigName set a config file name to default configuration instance.  The default config file name is "app.yml"
func SetConfigName(configName string) {
	cfg.SetConfigName(configName)
}

// SetEnvPrefix set a prefix for env.
func SetEnvPrefix(prefix string) {
	cfg.SetEnvPrefix(prefix)
}

// String returns a string type value which has the key.
func String(key string, defaultValue ...string) (string, error) {
	val, err := cfg.Get(key)

	if err == ErrKeyNotFound {
		if len(defaultValue) > 0 {
			return defaultValue[0], nil
		}

		return "", ErrKeyNotFound
	}

	return cast.ToString(val)
}

// Int32 returns a int32 type value which has the key.
func Int32(key string, defaultValue ...int32) (int32, error) {
	val, err := cfg.Get(key)

	if err == ErrKeyNotFound {
		if len(defaultValue) > 0 {
			return defaultValue[0], nil
		}

		return 0, ErrKeyNotFound
	}

	return cast.ToInt32(val)
}

// UnmarshalKey binds a value which has the key.
func UnmarshalKey(key string, value interface{}) error {
	data, err := cfg.Get(key)

	if err != nil {
		return err
	}

	err = mapstructure.Decode(data, value)

	if err != nil {
		return err
	}

	return nil
}
