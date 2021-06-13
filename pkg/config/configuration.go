package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

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
	AddPath(path string)
	String(key string, defaultValue ...string) (string, error)
	Int32(key string, defaultValue ...int32) (int32, error)
	UnmarshalKey(key string, val interface{}) error
	Set(key string, val string) error
}

type Config struct {
	content    []byte
	fileName   string
	configType string
	paths      []string
	mu         sync.RWMutex
	cache      map[string]interface{}
}

func New() Configuration {
	cfg := Config{
		content:    []byte{},
		fileName:   "app.yml",
		configType: "yaml",
		mu:         sync.RWMutex{},
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

// AddPath adds a path to look for config file.
func (cfg *Config) AddPath(path string) {
	cfg.paths = append(cfg.paths, path)
}

// String returns a string type value which has the key.  If the value can't convert to string type,
func (cfg *Config) String(key string, defaultValue ...string) (string, error) {
	cfg.mu.RLock()
	defer cfg.mu.RUnlock()

	val, found := cfg.cache[key]
	if !found {
		if len(defaultValue) > 0 {
			return defaultValue[0], nil
		}

		return "", ErrKeyNotFound
	}

	return cast.ToString(val)
}

func (cfg *Config) UnmarshalKey(key string, val interface{}) error {
	return nil
}

// Int32 returns a int32 type value which has the key.  If the value can't convert to string type,
func (cfg *Config) Int32(key string, defaultValue ...int32) (int32, error) {
	cfg.mu.RLock()
	defer cfg.mu.RUnlock()

	val, found := cfg.cache[key]
	if !found {
		if len(defaultValue) > 0 {
			return defaultValue[0], nil
		}

		return 0, ErrKeyNotFound
	}

	return cast.ToInt32(val)
}

// Set set a new value with key into config.  If the key doesn't exist, a new key will be created and no error be returned.
func (cfg *Config) Set(key string, val string) error {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	cfg.cache[key] = val

	return nil
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
	items := map[string]interface{}{}

	switch cfg.configType {
	case "yaml", "yml":
		err := yaml.Unmarshal(cfg.content, &items)
		if err != nil {
			return err
		}
	case "json":
	default:
		return ErrConfigTypeNotSupport
	}

	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	for k, v := range items {
		cfg.buildCache(k, v)
	}

	return nil
}

func (cfg *Config) buildCache(key string, val interface{}) {
	if val == nil {
		return
	}

	myArray, ok := val.([]interface{})
	if ok {
		for keyA, valA := range myArray {
			newKey := fmt.Sprintf("%s[%d]", key, keyA)
			cfg.buildCache(newKey, valA)
		}

		return
	}

	myMap, ok := val.(map[string]interface{})
	if ok {
		for keyA, valA := range myMap {
			newKey := fmt.Sprintf("%s.%s", key, keyA)
			cfg.buildCache(newKey, valA)
		}

		return
	}

	myVal, ok := val.(interface{})
	if ok && myVal != nil {
		cfg.cache[key] = myVal
	}
}
