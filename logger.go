// Package logger is a simple wrapper around log.Logger with usual logging levels "error", "warning", "notice",
// "info" and "debug".
package logger

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/mgutz/ansi"
)

// Logger represents a logger instance.
type Logger struct {
	logger  *log.Logger
	level   int
	context string
	out     *os.File
}

const (
	_ = iota
	levelError
	levelWarning
	levelNotice
	levelInfo
	levelDebug

	defaultLevel = "info"
)

var (
	levelLabels map[int]string
	levelMap    = map[string]int{
		"error":   levelError,
		"warning": levelWarning,
		"notice":  levelNotice,
		"info":    levelInfo,
		"debug":   levelDebug,
	}
)

// NewLogger returns a new Logger instance initialized with the logging system output path and level. If logPath is
// either empty or "-", logging will be output to os.Stderr. Log messages with severity higher than level will be
// discarded.
func NewLogger(logPath, level string) (*Logger, error) {
	var (
		err error
		ok  bool
	)

	logger := &Logger{level: levelMap[defaultLevel]}

	if logger.level, ok = levelMap[level]; !ok {
		return nil, ErrInvalidLevel
	}

	if logPath != "" && logPath != "-" {
		// Set logging output to a file
		ansi.DisableColors(true)

		// Create parent folders if needed
		dirPath, _ := path.Split(logPath)

		if err = os.MkdirAll(dirPath, 0755); err != nil {
			return nil, err
		}

		if logger.out, err = os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err != nil {
			return nil, err
		}
	} else {
		// Set logging output to stderr
		logger.out = os.Stderr
	}

	levelLabels = map[int]string{
		levelError:   ansi.Color("ERROR", "red"),
		levelWarning: ansi.Color("WARNING", "yellow"),
		levelNotice:  ansi.Color("NOTICE", "magenta"),
		levelInfo:    ansi.Color("INFO", "blue"),
		levelDebug:   ansi.Color("DEBUG", "cyan"),
	}

	logger.logger = log.New(logger.out, "", log.LstdFlags|log.Lmicroseconds)

	return logger, nil
}

// Logger returns the underlying log.Logger instance.
func (l *Logger) Logger() *log.Logger {
	return l.logger
}

// Context clones the Logger instance and sets the context to the provided string.
func (l *Logger) Context(context string) *Logger {
	logger := *l
	logger.context = context
	return &logger
}

// Error prints an error message in the logging system.
func (l *Logger) Error(format string, v ...interface{}) *Logger {
	return l.print(levelError, format, v...)
}

// Warning prints a warning message in the logging system.
func (l *Logger) Warning(format string, v ...interface{}) *Logger {
	return l.print(levelWarning, format, v...)
}

// Notice prints a notice message in the logging system.
func (l *Logger) Notice(format string, v ...interface{}) *Logger {
	return l.print(levelNotice, format, v...)
}

// Info prints an information message in the logging system.
func (l *Logger) Info(format string, v ...interface{}) *Logger {
	return l.print(levelInfo, format, v...)
}

// Debug prints a debug message in the logging system.
func (l *Logger) Debug(format string, v ...interface{}) *Logger {
	return l.print(levelDebug, format, v...)
}

// Close closes the logger output file.
func (l *Logger) Close() error {
	if l.out != nil {
		return l.out.Close()
	}

	return nil
}

func (l *Logger) print(level int, format string, v ...interface{}) *Logger {
	if level > l.level {
		return l
	}

	if l.context != "" {
		l.logger.Printf(
			"%s: %s",
			fmt.Sprintf("%s: %s", levelLabels[level], l.context),
			fmt.Sprintf(format, v...),
		)
	} else {
		l.logger.Printf("%s: %s", levelLabels[level], fmt.Sprintf(format, v...))
	}

	return l
}
