package file

import (
	"errors"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"github.com/nite-coder/blackbear/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// convert: https://onlineyamltools.com/convert-yaml-to-json

var (
	yamlContent = `
env: "test"
app:
  id: blackbear

logs:
  - name: clog
    type: console
    min_level: debug
  - name: graylog
    type: gelf
    min_level: debug

datasource:
  - name1
  - name2
  - name3

web:
  port: 10080
  ping: true
`

	newYamlContent = `
env: "test1"
app:
  id: blackbear1
`

	jsonContent = `
{
	"env": "test",
	"app": {
		"id": "blackbear"
	},
	"logs": [
		{
		"name": "clog",
		"type": "console",
		"min_level": "debug"
		},
		{
		"name": "graylog",
		"type": "gelf",
		"min_level": "debug"
		}
	],
	"datasource": [
		"name1",
		"name2",
		"name3"
	],
	"web": {
		"port": 10080,
		"ping": true
	}
}
`
)

func TestConfig(t *testing.T) {
	tests := []struct {
		format  ConfigType
		content string
	}{
		{
			ConfigTypeYAML,
			yamlContent,
		},
		{
			ConfigTypeJSON,
			jsonContent,
		},
	}

	for _, test := range tests {
		fileProvder := New()
		fileProvder.SetConfigType(test.format)
		err := fileProvder.LoadContent(test.content)
		require.NoError(t, err)

		val, err := fileProvder.Get("env")
		require.NoError(t, err)
		require.Equal(t, "test", val)

		val, err = fileProvder.Get("app.id")
		require.NoError(t, err)
		require.Equal(t, "blackbear", val)

		_, err = fileProvder.Get("hello")
		require.ErrorIs(t, config.ErrKeyNotFound, err)

		val, err = fileProvder.Get("datasource.2")
		assert.NoError(t, err)
		assert.Equal(t, "name3", val)

		_, err = fileProvder.Get("datasource.4")
		assert.ErrorIs(t, config.ErrKeyNotFound, err)

		val, err = fileProvder.Get("logs.1.name")
		assert.NoError(t, err)
		assert.Equal(t, "graylog", val)

	}
}

func TestLoadPathOrder(t *testing.T) {
	fileProvder := New()
	err := fileProvder.Load()
	assert.True(t, errors.Is(err, config.ErrFileNotFound))
}

func TestWatchConfig(t *testing.T) {
	// create temp config
	tmpDir := os.TempDir()
	tmpFile, err := os.CreateTemp(tmpDir, "config.*.yaml")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	// fmt.Println("file: " + tmpFile.Name())

	text := []byte(yamlContent)
	_, err = tmpFile.Write(text)
	require.NoError(t, err)
	tmpFile.Close()

	fileProvder := New()
	fileProvder.AddPath(tmpDir)
	fileProvder.SetConfigName(filepath.Base(tmpFile.Name()))
	err = fileProvder.Load()
	require.NoError(t, err)

	err = fileProvder.WatchConfig()
	require.NoError(t, err)

	config.AddProvider(fileProvder)

	var count int64
	config.OnChangedEvent(func() error {
		atomic.AddInt64(&count, int64(1))
		val, err := config.String("env")
		require.NoError(t, err)
		require.Equal(t, "test1", val)
		return nil
	})

	val, err := config.String("env")
	require.NoError(t, err)
	require.Equal(t, "test", val)

	// replace to new config
	f, err := os.OpenFile(tmpFile.Name(), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	require.NoError(t, err)

	text = []byte(newYamlContent)
	_, err = f.Write(text)
	require.NoError(t, err)
	f.Close()

	time.Sleep(3 * time.Second) // wait for onChangedEvent fired

	assert.Equal(t, int64(1), atomic.LoadInt64(&count))
}
