package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	yamlContent = `
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

	val, err := String("app.id")
	assert.NoError(t, err)
	assert.Equal(t, "blackbear", val)

	_, err = String("id")
	assert.ErrorIs(t, ErrKeyNotFound, err)

	noVal := "no id value"
	val, err = String("id", noVal)
	assert.NoError(t, err)
	assert.Equal(t, noVal, val)

	_ = Set("app.id", "HelloApp")
	val, err = cfg.String("app.id")
	assert.NoError(t, err)
	assert.Equal(t, "HelloApp", val)

	int32Result, err := Int32("web.port")
	assert.NoError(t, err)
	assert.Equal(t, int32(10080), int32Result)
}

func TestUnmarshalKey(t *testing.T) {
	err := LoadContent(yamlContent)
	require.NoError(t, err)
}
