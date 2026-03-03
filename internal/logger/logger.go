package logger

import (
	"fmt"
	"io"
	"log"
	"os"
)

type Logger struct {
	file      *os.File
	stdLogger *log.Logger
	verbose   bool
}

func NewLogger(logPath string, toScreen bool, verbose bool) (*Logger, error) {
	var writer io.Writer
	var file *os.File

	if logPath == "none" {

		writer = os.Stdout
	} else if logPath == "default" {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("error getting current working directory: %w", err)
		}
		logPath = fmt.Sprintf("%s/log.log", cwd)

		file, err = os.Create(logPath)
		if err != nil {
			return nil, fmt.Errorf("error creating log file: %w", err)
		}

		writer = io.MultiWriter(os.Stdout, file)

	} else {

		file, err := os.Create(logPath)
		if err != nil {
			return nil, fmt.Errorf("error creating log file: %w", err)
		}

		writer = io.MultiWriter(os.Stdout, file)
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
