package logging

import (
	"bufio"
	"context"
	"log/slog"
	"os"
	"path/filepath"
)

const LOGBUFFERSIZE = 64 * 1024

type Logger struct {
	file      *os.File
	writer    *bufio.Writer
	stdLogger *slog.Logger
}

func NewLogger(logPath string, verbose bool) (*Logger, error) {
	var file *os.File
	var bWriter *bufio.Writer
	var err error

	// 1. Handle the file opening
	if logPath != "none" {
		if logPath == "default" {
			cwd, _ := os.Getwd()
			logPath = filepath.Join(cwd, "slog.log")
		}
		file, err = os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}
	}

	// 2. Create the Handlers
	var handlers []slog.Handler

	if file != nil {
		// Use a buffered writer for file performance
		bWriter = bufio.NewWriterSize(file, LOGBUFFERSIZE)
		handlers = append(handlers, slog.NewJSONHandler(bWriter, nil))
	}

	if verbose {
		opts := &slog.HandlerOptions{Level: slog.LevelWarn} // Only show Warnings/Errors on terminal to avoid clutter
		handlers = append(handlers, slog.NewTextHandler(os.Stdout, opts))
	}

	// 3. Combine the handlers

	combinedHandler := NewMultiHandler(handlers...)

	return &Logger{
		file:      file,
		writer:    bWriter,
		stdLogger: slog.New(combinedHandler),
	}, nil
}

func (l *Logger) Info(msg string, args ...any) { //replaces .log
	l.stdLogger.Info(msg, args...)
	os.Stdout.Sync()
}

func (l *Logger) Close() {
	if l.writer != nil {
		// ensuring anything not written gets flushed from buffer to be written.
		l.writer.Flush()
	}
	if l.file != nil { //
		// tell the OS to push data from its own kernel buffers to the physical disk before closing
		l.file.Sync()
		l.file.Close()
	}
}

func (l *Logger) Error(msg string, args ...any) {
	l.stdLogger.Error(msg, args...)
}

// middleware for giving the handler multiwriter functionality
type MultiHandler struct {
	handlers []slog.Handler
}

func NewMultiHandler(handlers ...slog.Handler) *MultiHandler {
	return &MultiHandler{handlers: handlers}
}

func (m *MultiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, h := range m.handlers {
		if h.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (m *MultiHandler) Handle(ctx context.Context, r slog.Record) error {
	var errs []error
	for _, h := range m.handlers {
		if err := h.Handle(ctx, r); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return errs[0] // Or return a combined error
	}
	return nil
}

func (m *MultiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		newHandlers[i] = h.WithAttrs(attrs)
	}
	return NewMultiHandler(newHandlers...)
}

func (m *MultiHandler) WithGroup(name string) slog.Handler {
	newHandlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		newHandlers[i] = h.WithGroup(name)
	}
	return NewMultiHandler(newHandlers...)
}
