package file

import (
	"testing"

	"github.com/nite-coder/blackbear/pkg/config"
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
	fileProvder := New()
	err := fileProvder.LoadContent(yamlContent)
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
}
