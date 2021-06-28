package file

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/nite-coder/blackbear/pkg/cast"
	"github.com/nite-coder/blackbear/pkg/config"
	"gopkg.in/yaml.v3"
)

type ConfigType string

const (
	ConfigTypeYAML ConfigType = "yaml"
	ConfigTypeJSON ConfigType = "json"
)

type FileProvider struct {
	content    []byte
	configName string
	configType ConfigType
	paths      []string
	cache      map[string]interface{}
}

func New() *FileProvider {
	return &FileProvider{
		content:    []byte{},
		configName: "app.yml",
		configType: ConfigTypeYAML,
		cache:      map[string]interface{}{},
	}
}

// ConfigName return config file name.  The default config file name is "app.yml"
func (p *FileProvider) ConfigName() string {
	return p.configName
}

// SetConfigName set a config file name.  The default config file name is "app.yml"
func (p *FileProvider) SetConfigName(configName string) {
	if len(p.configName) == 0 {
		return
	}
	p.configName = configName
}

// ConfigType set the config encoding type.  Default is "YAML".  `YAML`, `JSON` are supported
func (p *FileProvider) ConfigType() ConfigType {
	return p.configType
}

// SetConfigType set the config encoding type.  Default is "YAML".  `YAML`, `JSON` are supported
func (p *FileProvider) SetConfigType(configType ConfigType) {
	p.configType = configType
}

// AddPath adds a path to look for config file.  Please don't include filename. Directory only
func (p *FileProvider) AddPath(path string) {
	p.paths = append(p.paths, path)
}

// Load initialize this package. It will load config into cache and get ready to work.  However,
// if the config file was not found, `ErrFileNotFound` will be returned
func (p *FileProvider) Load() error {
	var err error

	if len(p.paths) == 0 {
		path, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		p.AddPath(filepath.Join(path, "config"))
		p.AddPath(path)

		// load config file from executed file's sub config folder
		path, err = os.Executable()
		if err != nil {
			panic(err)
		}

		p.AddPath(filepath.Join(path, "config"))
		p.AddPath(filepath.Dir(path))
	}

	for idx, path := range p.paths {
		// found config file
		if len(p.content) > 0 {
			break
		}

		if len(path) == 0 {
			continue
		}

		configFilePath := filepath.Join(path, p.configName)
		p.content, err = ioutil.ReadFile(filepath.Clean(configFilePath))

		if err != nil {

			if errors.Is(err, os.ErrNotExist) {

				if (idx + 1) == len(p.paths) {
					return config.ErrFileNotFound
				}

				continue
			}

			return fmt.Errorf("config: read file error: %w", err)
		}
	}

	return p.start()
}

// LoadContent reads the content as config file
func (p *FileProvider) LoadContent(content string) error {
	content = strings.TrimSpace(content)
	p.content = []byte(content)

	return p.start()
}

func (p *FileProvider) start() error {
	switch p.configType {
	case ConfigTypeYAML:
		err := yaml.Unmarshal(p.content, &p.cache)
		if err != nil {
			return fmt.Errorf("config: yaml unmarshal failed. err: %w", err)
		}
	case ConfigTypeJSON:
		err := json.Unmarshal(p.content, &p.cache)
		if err != nil {
			return fmt.Errorf("config: json unmarshal failed. err: %w", err)
		}
	default:
		return config.ErrConfigTypeNotSupport
	}

	return nil
}

func (p *FileProvider) Get(key string) (interface{}, error) {
	if len(key) == 0 {
		return p.cache, nil
	}

	var lastOne, found bool
	keys := strings.Split(key, ".")
	var temp interface{}

	temp = p.cache

	for idx, key := range keys {
		if idx == len(keys)-1 {
			lastOne = true
		}

		if temp == nil {
			return nil, config.ErrKeyNotFound
		}

		myMap, ok := temp.(map[string]interface{})

		if ok {
			temp, found = myMap[key]

			if !found {
				return nil, config.ErrKeyNotFound
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
				return nil, config.ErrKeyNotFound
			}

			if arIdx >= len(myArray) {
				return nil, config.ErrKeyNotFound
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

	return nil, config.ErrKeyNotFound
}
