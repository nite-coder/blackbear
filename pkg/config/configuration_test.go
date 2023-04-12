package config_test

import (
	"errors"
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
  timeout: 60s

money: 123.42

logs:
  - name: clog
    type: console
    min_level: debug
    timeout: 5s
  - name: graylog
    type: gelf
    min_level: info
    timeout: 5s

datasource:
  - name1
  - name2
  - name3

currency: ["btc", "usdt", "usd"]
users: [1,2,3]

web:
  port: 10080
  ping: true

book:
  book1: john
  book2: angela
`
)

type LogItem struct {
	Name     string
	Type     string
	MinLevel string        `yaml:"min_level" mapstructure:"min_level"`
	Timeout  time.Duration `yaml:"timeout" mapstructure:"timeout"`
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
	require.True(t, errors.Is(err, config.ErrKeyNotFound))

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

	timeout, _ := config.Duration("app.timeout", 180*time.Second)
	assert.Equal(t, time.Duration(60*time.Second), timeout)

	// string to map
	bookMap, err := config.StringMap("book")
	assert.NoError(t, err)
	assert.Equal(t, "angela", bookMap["book2"])

	bookMapString, err := config.StringMapString("book")
	assert.NoError(t, err)
	assert.Equal(t, "angela", bookMapString["book2"])

	stringSlice, err := config.StringSlice("currency")
	assert.NoError(t, err)
	assert.Len(t, stringSlice, 3)
	assert.Equal(t, stringSlice[2], "usd")

	defaultStringSlice := []string{"usdt"}
	stringSlice, err = config.StringSlice("currency1", defaultStringSlice)
	assert.NoError(t, err)
	assert.Len(t, stringSlice, 1)
	assert.Equal(t, stringSlice[0], "usdt")

	intSlice, err := config.IntSlice("users")
	assert.NoError(t, err)
	assert.Len(t, intSlice, int(3))
	assert.Equal(t, int(3), intSlice[2])

	defaultIntSlice := []int{100}
	intSlice, err = config.IntSlice("users1", defaultIntSlice)
	assert.NoError(t, err)
	assert.Len(t, intSlice, 1)
	assert.Equal(t, intSlice[0], 100)
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
	assert.Equal(t, float64(5), logSetting[0].Timeout.Seconds())

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
