package console

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"sync"

	"github.com/fatih/color"
	"github.com/nite-coder/blackbear/pkg/log"
)

var (
	debugLevelColor   = color.New(color.Bold).Add(color.FgGreen)
	infoLevelColor    = color.New(color.Bold).Add(color.FgBlue)
	warnLevelColor    = color.New(color.Bold).Add(color.FgYellow)
	errorLevelColor   = color.New(color.Bold).Add(color.FgHiRed)
	defaultLevelColor = color.New(color.Bold).Add(color.FgHiRed)

	keyColor   = color.New(color.FgCyan)
	valueColor = color.New(color.FgWhite)
)

func levelToColor(level string) *color.Color {
	switch level {
	case "DEBUG":
		return debugLevelColor
	case "INFO":
		return infoLevelColor
	case "WARN":
		return warnLevelColor
	case "ERROR", "PANIC", "FATAL":
		return errorLevelColor
	default:
		return defaultLevelColor
	}
}

// Console is an instance of the console logger
type Console struct {
	mutex  sync.Mutex
	writer io.Writer

	DisableColor bool
}

type ConsoleOptions struct {
	DisableColor bool
}

// New create a new Console instance
func New(opts ConsoleOptions) log.Handler {

	h := Console{
		writer: os.Stdout,
	}

	color.NoColor = true

	if !opts.DisableColor {
		color.NoColor = false
	}

	return &h
}

// BeforeWriting handles the log entry
func (h *Console) BeforeWriting(e *log.Entry) error {
	e.Str("level", e.Level.String())
	e.Time("time", e.CreatedAt)

	return nil
}

// Write handles the log entry
func (h *Console) Write(bytes []byte) error {
	kv := map[string]interface{}{}
	err := json.Unmarshal(bytes, &kv)

	if err != nil {
		return err
	}

	level := fmt.Sprintf("%v", kv["level"])
	msg := kv["msg"]
	levelColor := levelToColor(level)
	time := kv["time"]

	// sort map by key
	keys := make([]string, 0, len(kv))

	for k := range kv {
		if k == "level" || k == "msg" {
			continue
		}

		keys = append(keys, k)
	}

	sort.Strings(keys)
	// fmt is not goroutine safe
	// https://stackoverflow.com/questions/14694088/is-it-safe-for-more-than-one-goroutine-to-print-to-stdout
	h.mutex.Lock()
	defer h.mutex.Unlock()

	_, _ = fmt.Fprintf(h.writer, "%s %s %s", time, levelColor.Sprintf("%-6s", level), msg)

	for _, k := range keys {
		if k == "time" {
			continue
		}
		_, _ = fmt.Fprintf(h.writer, " %s=%v", keyColor.Sprint(k), valueColor.Sprintf("%v", kv[k]))
	}

	fmt.Fprintln(h.writer)

	return nil
}
