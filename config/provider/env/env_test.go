package env

import (
	"os"
	"testing"

	"github.com/nite-coder/blackbear/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnv(t *testing.T) {
	envProvider := New()
	config.AddProvider(envProvider)

	err := os.Setenv("ENV", "test")
	if err != nil {
		panic(err)
	}

	val, err := envProvider.Get("env")
	require.NoError(t, err)
	assert.Equal(t, "test", val)

	err = os.Unsetenv("ENV")
	if err != nil {
		panic(err)
	}

	// not found
	_, err = envProvider.Get("env1")
	require.ErrorIs(t, config.ErrKeyNotFound, err)

	// prefix
	envProvider.SetEnvPrefix("BLACKBEAR")

	err = os.Setenv("BLACKBEAR_MODE", "debug")
	if err != nil {
		panic(err)
	}

	val, err = envProvider.Get("MODE")
	require.NoError(t, err)
	assert.Equal(t, "debug", val)

	err = os.Unsetenv("BLACKBEAR_MODE")
	if err != nil {
		panic(err)
	}
}
