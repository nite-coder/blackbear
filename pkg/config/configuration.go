package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v3"
)

var (
	cfg             = New()
	defaultFileName = "app.yml"
	configType      = "yaml"
	ErrFileNotFound = errors.New("config file was not found")
)

type Configuration interface {
	Load() error
	AddPath(path string)
	String(key string, defaultValue ...string) (string, error)
	Set(key string, val string) error
}

type Config struct {
	paths []string
	mu    sync.RWMutex
	cache map[string]interface{}
}

func New() Configuration {
	cfg := &Config{
		mu:    sync.RWMutex{},
		cache: map[string]interface{}{},
	}

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

	return cfg
}

// SetFileName set a new config file name.  The default config file name is "app.yml"
func (cfg *Config) SetFileName(fileName string) {
	defaultFileName = fileName
}

// AddPath adds a path to look for config file.
func (cfg *Config) AddPath(path string) {
	cfg.paths = append(cfg.paths, path)
}

func (cfg *Config) String(key string, defaultValue ...string) (string, error) {
	cfg.mu.RLock()
	defer cfg.mu.RUnlock()

	val, ok := cfg.cache[key].(string)
	if !ok {
		return "", errors.New("the value is not string type")
	}
	return val, nil
}

func (cfg *Config) Set(key string, val string) error {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	cfg.cache[key] = val
	return nil
}

// Load initialize this package. It will load config into cache and get ready to work.  However,
// if the config file was not found, `ErrFileNotFound` will be returned
func (cfg *Config) Load() error {
	var file []byte
	var err error
	for idx, path := range cfg.paths {
		// found config file
		if len(file) > 0 {
			break
		}
		if len(path) == 0 {
			continue
		}

		configFilePath := filepath.Join(path, defaultFileName)
		file, err = ioutil.ReadFile(filepath.Clean(configFilePath))
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

	items := map[string]interface{}{}
	err = yaml.Unmarshal(file, &items)
	if err != nil {
		return err
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

// Cfg return a singleton configuration instance.
func Cfg() Configuration {
	return cfg
}
