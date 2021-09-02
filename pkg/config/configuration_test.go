package config_test

import (
	"os"
	"testing"
	"time"

	"github.com/nite-coder/blackbear/pkg/config"
	"github.com/nite-coder/blackbear/pkg/config/provider/env"
	"github.com/nite-coder/blackbear/pkg/config/provider/file"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	yamlContent = `
env: "test"
app:
  id: blackbear
  timeout_sec: 60

money: 123.42

logs:
  - name: clog
    type: console
    min_level: debug
  - name: graylog
    type: gelf
    min_level: info

datasource:
  - name1
  - name2
  - name3

web:
  port: 10080
  ping: true
`
)

type LogItem struct {
	Name     string
	Type     string
	MinLevel string `yaml:"min_level" mapstructure:"min_level"`
}

type Web struct {
	Port int
	Ping bool
}

type Root struct {
	Env        string
	Web        Web
	Datasource []string  `yaml:"datasource"`
	Logs       []LogItem `yaml:"logs"`
}

func TestNoProvider(t *testing.T) {
	config.RemoveAllPrividers()

	_, err := config.String("hello")
	require.ErrorIs(t, config.ErrKeyNotFound, err)

	count, _ := config.Int("app.wokrer_count", 5)
	assert.Equal(t, 5, count)
}

func TestAddProviders(t *testing.T) {
	config.RemoveAllPrividers()

	envProvider := env.New()
	config.AddProvider(envProvider)

	fileProvder := file.New()
	err := fileProvder.LoadContent(yamlContent)
	require.NoError(t, err)
	config.AddProvider(fileProvder)

	err = os.Setenv("ENV", "first")
	if err != nil {
		panic(err)
	}

	env, _ := config.String("env")
	assert.Equal(t, "first", env)

	appID, _ := config.String("app.id")
	assert.Equal(t, "blackbear", appID)

	err = os.Unsetenv("ENV")
	if err != nil {
		panic(err)
	}
}

func TestConverterType(t *testing.T) {
	config.RemoveAllPrividers()

	fileProvder := file.New()
	err := fileProvder.LoadContent(yamlContent)
	require.NoError(t, err)
	config.AddProvider(fileProvder)

	val, err := config.String("app.id")
	require.NoError(t, err)
	assert.Equal(t, "blackbear", val)

	intVal, err := config.Int("web.port")
	assert.NoError(t, err)
	assert.Equal(t, 10080, intVal)

	int32Result, err := config.Int32("web.port")
	assert.NoError(t, err)
	assert.Equal(t, int32(10080), int32Result)

	int64Val, err := config.Int64("web.port")
	assert.NoError(t, err)
	assert.Equal(t, int64(10080), int64Val)

	boolVal, err := config.Bool("web.ping")
	assert.NoError(t, err)
	assert.Equal(t, true, boolVal)

	float32Val, err := config.Float32("money")
	assert.NoError(t, err)
	assert.Equal(t, float32(123.42), float32Val)

	float64Val, err := config.Float64("money")
	assert.NoError(t, err)
	assert.Equal(t, 123.42, float64Val)

	defaultValue := "default value"
	val, err = config.String("logs.Hello", defaultValue)
	assert.NoError(t, err)
	assert.Equal(t, defaultValue, val)

	timeoutSec, _ := config.Duration("app.timeout_sec", 180)
	assert.Equal(t, time.Duration(60), timeoutSec)
}

func TestScan(t *testing.T) {
	config.RemoveAllPrividers()

	fileProvder := file.New()
	err := fileProvder.LoadContent(yamlContent)
	require.NoError(t, err)
	config.AddProvider(fileProvder)

	logSetting := []LogItem{}

	err = config.Scan("logs", &logSetting)
	require.NoError(t, err)
	assert.Equal(t, "clog", logSetting[0].Name)

	data := []string{}
	err = config.Scan("datasource", &data)
	require.NoError(t, err)
	assert.Equal(t, 3, len(data))
	assert.Equal(t, "name2", data[1])

	env := ""
	err = config.Scan("env", &env)
	require.NoError(t, err)
	assert.Equal(t, "test", env)

	web := Web{}
	err = config.Scan("web", &web)
	require.NoError(t, err)
	assert.Equal(t, 10080, web.Port)
	assert.Equal(t, true, web.Ping)

	root := Root{}
	err = config.Scan("", &root)
	require.NoError(t, err)
	assert.Equal(t, 10080, root.Web.Port)
	assert.Equal(t, true, root.Web.Ping)
	assert.Equal(t, "test", root.Env)
	assert.Equal(t, "clog", root.Logs[0].Name)
	assert.Equal(t, "console", root.Logs[0].Type)
	assert.Equal(t, "debug", root.Logs[0].MinLevel)
	assert.Equal(t, "name1", root.Datasource[0])
}
