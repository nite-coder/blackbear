package log

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/fatih/color"
)

// color: https://en.wikipedia.org/wiki/ANSI_escape_code#Colors
var (
	timeColor = color.New(90) // gray

	debugLevelColor   = color.New(color.Bold, color.FgGreen)
	infoLevelColor    = color.New(color.Bold, color.FgBlue)
	warnLevelColor    = color.New(color.Bold, color.FgYellow)
	errorLevelColor   = color.New(color.Bold, color.FgHiRed)
	defaultLevelColor = color.New(color.Bold, color.FgHiRed)

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

// TextHandler is an instance of the text handler
type TextHandler struct {
	mu     sync.Mutex
	writer io.Writer

	opts *HandlerOptions
}

// New create a new Console instance
func NewTextHandler(w io.Writer, opts *HandlerOptions) *TextHandler {
	if w == nil {
		w = io.Discard
	}

	h := TextHandler{
		writer: w,
		opts:   opts,
	}

	color.NoColor = true
	if !opts.DisableColor {
		color.NoColor = false
	}

	return &h
}

// Enabled reports whether the handler handles records at the given level.
// The handler ignores records whose level is lower.
func (h *TextHandler) Enabled(_ context.Context, level Level) bool {
	return level >= h.opts.Level
}

// Handle formats its argument Record as a text object on a single line.
func (h *TextHandler) Handle(_ context.Context, e *Entry) error {

	level := e.Level.String()
	levelColor := levelToColor(level)

	// fmt is not goroutine safe
	// https://stackoverflow.com/questions/14694088/is-it-safe-for-more-than-one-goroutine-to-print-to-stdout
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.opts.DisableTime {
		time := time.Now().Format("03:04:05.000")
		_, _ = fmt.Fprintf(h.writer, "%s %s %s", timeColor.Sprint(time), levelColor.Sprintf("%-6s", level), e.Message)
	} else {
		_, _ = fmt.Fprintf(h.writer, "%s %s", levelColor.Sprintf("%-6s", level), e.Message)
	}

	for _, field := range e.Fields() {
		k := field.Key
		_, _ = fmt.Fprintf(h.writer, " %s=%v", keyColor.Sprint(k), valueColor.Sprintf("%v", field.Value))
	}

	fmt.Fprintln(h.writer)

	return nil
}
