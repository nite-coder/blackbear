package file

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
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
	rwMutex           sync.RWMutex
	content           []byte
	contentHash       string
	configPath        string
	configName        string
	configType        ConfigType
	paths             []string
	cache             map[string]interface{}
	lastFileUpdatedAt time.Time
	eventChan         chan bool
}

func New() *FileProvider {
	return &FileProvider{
		content:    []byte{},
		configName: "app.yml",
		configType: ConfigTypeYAML,
		cache:      map[string]interface{}{},
		eventChan:  make(chan bool, 1),
	}
}

func (p *FileProvider) Content() string {
	return string(p.content)
}

// NotifyChange return a channel and notify you when file be changed.
func (p *FileProvider) NotifyChange() chan bool {
	return p.eventChan
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
	p.configPath, err = p.getConfigPath()
	if err != nil {
		return err
	}

	p.content, err = os.ReadFile(filepath.Clean(p.configPath))
	if err != nil {
		return fmt.Errorf("config: read file error: %w", err)
	}

	if len(p.contentHash) > 0 {
		hasher := sha256.New()
		hasher.Write(p.content)
		newHash := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

		if p.contentHash == newHash {
			return errors.New("content is the same")
		}

	}

	return p.decode()
}

// LoadContent reads the content as config file
func (p *FileProvider) LoadContent(content string) error {
	content = strings.TrimSpace(content)
	p.content = []byte(content)

	return p.decode()
}

func (p *FileProvider) decode() error {
	p.rwMutex.Lock()
	defer p.rwMutex.Unlock()

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

	p.lastFileUpdatedAt = time.Time{}

	hasher := sha256.New()
	hasher.Write(p.content)
	p.contentHash = base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return nil
}

func (p *FileProvider) getConfigPath() (string, error) {
	if len(p.paths) == 0 {
		path, err := os.Getwd()
		if err != nil {
			return "", err
		}

		p.AddPath(filepath.Join(path, "config"))
		p.AddPath(path)

		// load config file from executed file's sub config folder
		path, err = os.Executable()
		if err != nil {
			return "", err
		}
		path = filepath.Dir(path)

		p.AddPath(filepath.Join(path, "config"))
		p.AddPath(path)
	}

	configPath := ""
	for idx, path := range p.paths {
		if len(path) == 0 {
			continue
		}

		configPath = filepath.Clean(filepath.Join(path, p.configName))
		_, err := os.Stat(configPath)

		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				if (idx + 1) == len(p.paths) {
					return "", fmt.Errorf("%w.  Directories are %v, ConfigName is %s", config.ErrFileNotFound, p.paths, p.configName)
				}

				continue
			}

			return "", err
		}

		break
	}

	return configPath, nil
}

func (p *FileProvider) Get(key string) (interface{}, error) {
	p.rwMutex.RLock()
	defer p.rwMutex.RUnlock()

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

		if temp != nil {
			return temp, nil
		}
	}

	return nil, config.ErrKeyNotFound
}

func (p *FileProvider) WatchConfig() error {
	configPath, err := p.getConfigPath()
	if err != nil {
		return err
	}

	sw := sync.WaitGroup{}
	sw.Add(1)

	go func() {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Fatal(err)
		}
		defer watcher.Close()

		err = watcher.Add(configPath)
		if err != nil {
			log.Fatal(err)
		}

		sw.Done()
		isUpdate := false
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("watcher event:", event)
				// 1. The Write, Create, and Rename events will be fired when saving the data, depending on the text editor you are using. For example, vi uses Rename
				// 2. When a large amount of data is being saved, the Write event will be triggered multiple times. Hence, we utilize a ticker and the 'isUpdate' parameter here.
				if event.Op&fsnotify.Write == fsnotify.Write ||
					event.Op&fsnotify.Create == fsnotify.Create ||
					event.Op&fsnotify.Rename == fsnotify.Rename {
					isUpdate = true
					log.Println("modified file:", event.Name)
				} else if event.Op&fsnotify.Remove == fsnotify.Remove {
					// Some editors will remove the path from the watch list when the event is triggered, so we need to re-add it
					if fileExist(configPath) {
						log.Printf("reread  file:%s ", configPath)
						err := watcher.Add(configPath)
						if err != nil {
							log.Println("watcher error:", err)
						}
					}
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("watcher error:", err)

			case <-ticker.C:
				if !isUpdate {
					continue
				}
				isUpdate = false

				err := p.Load()
				if err != nil {
					log.Println("watcher error:", err)
					continue
				}

				p.eventChan <- true
			}
		}
	}()

	sw.Wait()
	return nil
}

func fileExist(file string) bool {
	_, err := os.Stat(file)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		if os.IsNotExist(err) {
			return false
		}
		return false
	}
	return true
}
