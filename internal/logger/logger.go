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
}

func NewLogger(logPath string, toScreen bool) (*Logger, error) {
	var writer io.Writer
	var file *os.File

	if logPath == "none" {
		writer = os.Stdout
	} else {
		var err error
		file, err = os.Create(logPath)
		if err != nil {
			return nil, fmt.Errorf("creating log file: %w", err)
		}

		if toScreen {
			writer = io.MultiWriter(os.Stdout, file)
		} else {
			writer = file
		}
	}

	return &Logger{
		file:      file,
		stdLogger: log.New(writer, "", log.LstdFlags),
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
