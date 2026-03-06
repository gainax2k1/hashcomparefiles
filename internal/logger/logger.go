package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

type Logger struct {
	file      *os.File
	stdLogger *log.Logger
	verbose   bool
}

func NewLogger(logPath string, verbose bool) (*Logger, error) {
	var writer io.Writer
	var file *os.File
	var err error
	var cwd string

	if logPath == "none" {
		writer = os.Stdout
	} else {
		if logPath == "default" {
			cwd, err = os.Getwd()
			if err != nil {
				return nil, fmt.Errorf("error getting current working directory: %w", err)
			}
			logPath = filepath.Join(cwd, "log.log")
		}
		// append log if exists, else create new log
		file, err = os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Could not open log file %s: %v. Logging to console instead.\n", logPath, err)
			writer = os.Stdout
		} else {
			writer = io.MultiWriter(os.Stdout, file)
		}
	}

	return &Logger{
		file:      file,
		stdLogger: log.New(writer, "", log.LstdFlags),
		verbose:   verbose,
	}, nil
}

func (l *Logger) Log(format string, args ...interface{}) {
	l.stdLogger.Printf(format, args...)

}

func (l *Logger) Close() {
	if l.file != nil {
		l.file.Close()
	}
}

func (l *Logger) Error(format string, args ...interface{}) {
	l.stdLogger.Printf("ERROR: "+format, args...)
}
