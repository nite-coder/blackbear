package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
)

func TestConfig(t *testing.T) {
	err := LoadContent(yamlContent)
	require.NoError(t, err)

	val, err := String("env")
	assert.NoError(t, err)
	assert.Equal(t, "test", val)

	val, err = String("app.id")
	assert.NoError(t, err)
	assert.Equal(t, "blackbear", val)

	_, err = String("hello")
	assert.ErrorIs(t, ErrKeyNotFound, err)

	defaultValue := "default value"
	val, err = String("logs.Hello", defaultValue)
	assert.NoError(t, err)
	assert.Equal(t, defaultValue, val)

	int32Result, err := Int32("web.port")
	assert.NoError(t, err)
	assert.Equal(t, int32(10080), int32Result)

	val, err = String("datasource.2")
	assert.NoError(t, err)
	assert.Equal(t, "name3", val)

	_, err = String("datasource.4")
	assert.ErrorIs(t, ErrKeyNotFound, err)
}

func TestEnv(t *testing.T) {
	err := os.Setenv("ENV", "from_env")
	if err != nil {
		panic(err)
	}

	err = LoadContent(yamlContent)
	require.NoError(t, err)

	val, err := String("env")
	require.NoError(t, err)
	assert.Equal(t, "from_env", val)

	err = os.Unsetenv("ENV")
	if err != nil {
		panic(err)
	}

	t.Run("prefix test", func(t *testing.T) {
		err = os.Setenv("HELLO_WEB_MODE", "debug")
		if err != nil {
			panic(err)
		}

		SetEnvPrefix("HELLO")

		val, err = String("web.MODE")
		require.NoError(t, err)
		assert.Equal(t, "debug", val)

		err = os.Unsetenv("HELLO_WEB_MODE")
		if err != nil {
			panic(err)
		}
	})
}

type LogItem struct {
	Name     string
	Type     string
	MinLevel string
}

type Web struct {
	Port int
	Ping bool
}

type Root struct {
	Env string
	Web Web
}

func TestUnmarshalKey(t *testing.T) {
	err := LoadContent(yamlContent)
	require.NoError(t, err)

	logSetting := []LogItem{}

	err = UnmarshalKey("logs", &logSetting)
	require.NoError(t, err)
	assert.Equal(t, "clog", logSetting[0].Name)

	data := []string{}
	err = UnmarshalKey("datasource", &data)
	require.NoError(t, err)
	assert.Equal(t, 3, len(data))
	assert.Equal(t, "name2", data[1])

	env := ""
	err = UnmarshalKey("env", &env)
	require.NoError(t, err)
	assert.Equal(t, "test", env)

	web := Web{}
	err = UnmarshalKey("web", &web)
	require.NoError(t, err)
	assert.Equal(t, 10080, web.Port)
	assert.Equal(t, true, web.Ping)

	root := Root{}
	err = UnmarshalKey("", &root)
	require.NoError(t, err)
	assert.Equal(t, 10080, root.Web.Port)
	assert.Equal(t, true, root.Web.Ping)
	assert.Equal(t, "test", root.Env)
}
