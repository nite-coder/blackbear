package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/nite-coder/blackbear/pkg/cast"
	"gopkg.in/yaml.v3"
)

var (
	cfg                     = New()
	ErrFileNotFound         = errors.New("config file was not found")
	ErrKeyNotFound          = errors.New("key was not found")
	ErrConfigTypeNotSupport = errors.New("config type is not support")
)

type Configuration interface {
	Load() error
	LoadContent(content string) error
	FileName() string
	SetFileName(fileName string)
	SetEnvPrefix(prefix string)
	AddPath(path string)
	String(key string, defaultValue ...string) (string, error)
	Int32(key string, defaultValue ...int32) (int32, error)
	UnmarshalKey(key string, val interface{}) error
}

type Config struct {
	content    []byte
	fileName   string
	configType string
	paths      []string
	envPrefix  string
	cache      map[string]interface{}
}

func New() Configuration {
	cfg := Config{
		content:    []byte{},
		fileName:   "app.yml",
		configType: "yaml",
		cache:      map[string]interface{}{},
	}

	return &cfg
}

// FileName return config file name.  The default config file name is "app.yml"
func (cfg *Config) FileName() string {
	return cfg.fileName
}

// SetFileName set a config file name.  The default config file name is "app.yml"
func (cfg *Config) SetFileName(fileName string) {
	cfg.fileName = fileName
}

// SetEnvPrefix set a prefix for env.
func (cfg *Config) SetEnvPrefix(prefix string) {
	cfg.envPrefix = prefix
}

// AddPath adds a path to look for config file.
func (cfg *Config) AddPath(path string) {
	cfg.paths = append(cfg.paths, path)
}

// String returns a string type value which has the key.  If the value can't convert to string type,
func (cfg *Config) String(key string, defaultValue ...string) (string, error) {
	val, err := cfg.find(key)

	if err == ErrKeyNotFound {
		if len(defaultValue) > 0 {
			return defaultValue[0], nil
		}

		return "", ErrKeyNotFound
	}

	return cast.ToString(val)
}

// UnmarshalKey binds a value which has the key.
func (cfg *Config) UnmarshalKey(key string, value interface{}) error {
	data, err := cfg.find(key)

	if err != nil {
		return err
	}

	err = mapstructure.Decode(data, value)

	if err != nil {
		return err
	}

	return nil
}

// Int32 returns a int32 type value which has the key.  If the value can't convert to string type,
func (cfg *Config) Int32(key string, defaultValue ...int32) (int32, error) {
	val, err := cfg.find(key)

	if err == ErrKeyNotFound {
		if len(defaultValue) > 0 {
			return defaultValue[0], nil
		}

		return 0, ErrKeyNotFound
	}

	return cast.ToInt32(val)
}

// Load initialize this package. It will load config into cache and get ready to work.  However,
// if the config file was not found, `ErrFileNotFound` will be returned
func (cfg *Config) Load() error {
	var err error

	if len(cfg.paths) == 0 {
		path, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		cfg.AddPath(filepath.Join(path, "config"))
		cfg.AddPath(path)

		// load config file from executed file's sub config folder
		path, err = os.Executable()
		if err != nil {
			panic(err)
		}

		cfg.AddPath(filepath.Join(path, "config"))
		cfg.AddPath(filepath.Dir(path))
	}

	for idx, path := range cfg.paths {
		// found config file
		if len(cfg.content) > 0 {
			break
		}

		if len(path) == 0 {
			continue
		}

		configFilePath := filepath.Join(path, cfg.fileName)
		cfg.content, err = ioutil.ReadFile(filepath.Clean(configFilePath))

		if err != nil {

			if errors.Is(err, os.ErrNotExist) {

				if (idx + 1) == len(cfg.paths) {
					return ErrFileNotFound
				}

				continue
			}

			return fmt.Errorf("config: read file error: %w", err)
		}
	}

	return cfg.start()
}

// LoadContent reads the content as config file
func (cfg *Config) LoadContent(content string) error {
	content = strings.TrimSpace(content)
	cfg.content = []byte(content)

	return cfg.start()
}

func (cfg *Config) start() error {
	switch cfg.configType {
	case "yaml", "yml":
		err := yaml.Unmarshal(cfg.content, &cfg.cache)
		if err != nil {
			return err
		}
	case "json":
	default:
		return ErrConfigTypeNotSupport
	}

	return nil
}

func getValueFromEnv(prefix, key string) (string, error) {
	key = strings.ReplaceAll(key, ".", "_")

	if len(prefix) > 0 {
		key = prefix + "_" + key
	}

	key = strings.ToUpper(key)
	val, present := os.LookupEnv(key)
	if present {
		return val, nil
	}
	return "", ErrKeyNotFound
}

func (cfg *Config) find(key string) (interface{}, error) {
	if len(key) == 0 {
		return cfg.cache, nil
	}

	val, err := getValueFromEnv(cfg.envPrefix, key)
	if err == nil {
		return val, nil
	}

	var lastOne, found bool
	keys := strings.Split(key, ".")
	var temp interface{}

	temp = cfg.cache

	for idx, key := range keys {
		if idx == len(keys)-1 {
			lastOne = true
		}

		if temp == nil {
			return nil, ErrKeyNotFound
		}

		myMap, ok := temp.(map[string]interface{})

		if ok {
			temp, found = myMap[key]

			if !found {
				return nil, ErrKeyNotFound
			}

			if lastOne {
				return temp, nil
			}

			continue

		}

		myArray, ok := temp.([]interface{})

		if ok {
			arIdx, err := cast.ToInt(key)

			if err != nil {
				return nil, ErrKeyNotFound
			}

			if arIdx >= len(myArray) {
				return nil, ErrKeyNotFound
			}

			temp = myArray[arIdx]

			if lastOne {
				return temp, nil
			}

			continue
		}

		myVal, ok := temp.(interface{})

		if ok && myVal != nil {
			return temp, nil
		}
	}

	return nil, ErrKeyNotFound
}
