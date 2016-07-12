package logger

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/mgutz/ansi"
)

const (
	_ = iota
	levelError
	levelWarning
	levelNotice
	levelInfo
	levelDebug

	defaultLevel string = "info"
)

var (
	logger   *log.Logger
	logLevel int

	levelLabels map[int]string
	levelSuffix string

	levelMap = map[string]int{
		"error":   levelError,
		"warning": levelWarning,
		"notice":  levelNotice,
		"info":    levelInfo,
		"debug":   levelDebug,
	}
)

// Init initializes the logging system output path and level.
func Init(logPath, level string) error {
	var (
		logOut *os.File
		ok     bool
		err    error
	)

	// Set default log level if none provided
	if level == "" {
		level = defaultLevel
	}

	logLevel, ok = levelMap[level]
	if !ok {
		return ErrInvalidLevel
	}

	if logPath != "" && logPath != "-" {
		ansi.DisableColors(true)
		levelSuffix = ":"

		// Create parent folders if needed
		dirPath, _ := path.Split(logPath)

		err = os.MkdirAll(dirPath, 0755)
		if err != nil {
			return err
		}

		// Open logging output file
		logOut, err = os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("unable to open log file: %s", err)
		}
		defer logOut.Close()
	} else {
		// Set logging output to 'stderr' and activate color
		logOut = os.Stderr
	}

	logger = log.New(logOut, "", log.LstdFlags|log.Lmicroseconds)

	// Set level labels
	levelLabels = map[int]string{
		levelError:   ansi.Color("ERROR", "red"),
		levelWarning: ansi.Color("WARNING", "yellow"),
		levelNotice:  ansi.Color("NOTICE", "magenta"),
		levelInfo:    ansi.Color("INFO", "blue"),
		levelDebug:   ansi.Color("DEBUG", "cyan"),
	}

	return nil
}

// Error prints an error message in the logging system.
func Error(context, format string, v ...interface{}) {
	printLog(levelError, context, format, v...)
}

// Warning prints a warning message in the logging system.
func Warning(context, format string, v ...interface{}) {
	printLog(levelWarning, context, format, v...)
}

// Notice prints a notice message in the logging system.
func Notice(context, format string, v ...interface{}) {
	printLog(levelNotice, context, format, v...)
}

// Info prints an information message in the logging system.
func Info(context, format string, v ...interface{}) {
	printLog(levelInfo, context, format, v...)
}

// Debug prints a debug message in the logging system.
func Debug(context, format string, v ...interface{}) {
	printLog(levelDebug, context, format, v...)
}

// printLog prints a message in the logging system.
func printLog(level int, context, format string, v ...interface{}) {
	if level > logLevel {
		return
	}

	logger.Printf("%s%s %s: %s", levelLabels[level], levelSuffix, context, fmt.Sprintf(format, v...))
}
