package json

import (
	"github.com/nite-coder/blackbear/pkg/log"
)

func Example_log() {
	logger := log.New()
	logger.DisableTimeField = true
	h := New()
	logger.AddHandler(h, log.AllLevels...)
	log.SetLogger(logger)

	log.Str("app_id", "blackbear").Debug("hello world")
	// Output:
	// {"app_id":"blackbear","level":"DEBUG","msg":"hello world"}
}
