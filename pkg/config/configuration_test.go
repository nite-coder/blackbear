package config

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/nite-coder/blackbear/internal/iofile"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	fmt.Println("run")
	// copy config file to executed file's directory
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	srcPath := filepath.Join(path, "../../test/config/app_test.yml")

	path, err = os.Executable()
	if err != nil {
		panic(err)
	}
	dstPath := filepath.Join(filepath.Dir(path), "app.yml")
	iofile.CopyFile(srcPath, dstPath)

	fmt.Printf("dstpath: %s", dstPath)
	m.Run()
	fmt.Println("end")
}

func TestConfig(t *testing.T) {
	cfg := Cfg()

	err := cfg.Load()
	assert.NoError(t, err)

	val, _ := cfg.String("app.id")
	assert.Equal(t, "blackbear", val)

	cfg.Set("app.id", "HelloApp")
	val, _ = cfg.String("app.id")
	assert.Equal(t, "HelloApp", val)
}
