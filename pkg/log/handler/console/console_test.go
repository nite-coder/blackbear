package console

import (
	"github.com/nite-coder/blackbear/pkg/log"
)

func Example_log() {
	logger := log.New()
	h := New()
	logger.AddHandler(h, log.AllLevels...)
	log.SetLogger(logger)

	log.Str("app_id", "blackbear").Debug("hello world")
	// Output:
	// DEBUG    hello world                                        app_id=blackbear
}
