package benchmarks

import (
	"context"
	"errors"
	"time"

	"github.com/nite-coder/blackbear/pkg/log"
)

const TestMessage = "Test logging, but use a somewhat realistic message length."

var (
	TestTime     = time.Date(2022, time.May, 1, 0, 0, 0, 0, time.UTC)
	TestString   = "7e3b3b2aaeff56a7108fe11e154200dd/7819479873059528190"
	TestInt      = 32768
	TestDuration = 23 * time.Second
	TestError    = errors.New("fail")
)

type disabledHandler struct{}

func (disabledHandler) Enabled(context.Context, log.Level) bool  { return false }
func (disabledHandler) Handle(context.Context, *log.Entry) error { panic("should not be called") }
