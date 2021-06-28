package file

import (
	"testing"

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
