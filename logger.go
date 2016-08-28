// Package logger is a simple wrapper around log.Logger with usual logging levels "error", "warning", "notice", "info"
// and "debug".
package logger

import "sync"

const defaultLevel = "info"

const (
	_ = iota
	// LevelError represents the error logging level.
	LevelError
	// LevelWarning represents the warning logging level.
	LevelWarning
	// LevelNotice represents the notice logging level.
	LevelNotice
	// LevelInfo represents the info logging level.
	LevelInfo
	// LevelDebug represents the debug logging level.
	LevelDebug
)

var levelMap = map[string]int{
	"error":   LevelError,
	"warning": LevelWarning,
	"notice":  LevelNotice,
	"info":    LevelInfo,
	"debug":   LevelDebug,
}

// Logger represents a logger instance.
type Logger struct {
	backends []backend
	context  string

	wg sync.WaitGroup

	sync.Mutex
}

// NewLogger returns a new Logger instance initialized with the given configuration.
func NewLogger(configs ...interface{}) (*Logger, error) {
	// Initialize logger backends
	logger := &Logger{
		backends: []backend{},
		wg:       sync.WaitGroup{},
	}

	for _, config := range configs {
		var (
			backend backend
			err     error
		)

		switch config.(type) {
		case FileConfig:
			backend, err = newFileBackend(config.(FileConfig), logger)

		case SyslogConfig:
			backend, err = newSyslogBackend(config.(SyslogConfig), logger)

		default:
			err = ErrUnsupportedBackend
		}

		if err != nil {
			return nil, err
		}

		logger.backends = append(logger.backends, backend)
	}

	return logger, nil
}

// Context clones the Logger instance and sets the context to the provided string.
func (l *Logger) Context(context string) *Logger {
	return &Logger{
		backends: l.backends,
		context:  context,
		wg:       sync.WaitGroup{},
	}
}

// Error prints an error message in the logging system.
func (l *Logger) Error(format string, v ...interface{}) *Logger {
	l.write(LevelError, format, v...)
	return l
}

// Warning prints a warning message in the logging system.
func (l *Logger) Warning(format string, v ...interface{}) *Logger {
	l.write(LevelWarning, format, v...)
	return l
}

// Notice prints a notice message in the logging system.
func (l *Logger) Notice(format string, v ...interface{}) *Logger {
	l.write(LevelNotice, format, v...)
	return l
}

// Info prints an information message in the logging system.
func (l *Logger) Info(format string, v ...interface{}) *Logger {
	l.write(LevelInfo, format, v...)
	return l
}

// Debug prints a debug message in the logging system.
func (l *Logger) Debug(format string, v ...interface{}) *Logger {
	l.write(LevelDebug, format, v...)
	return l
}

// Close closes the logger output file.
func (l *Logger) Close() {
	for _, b := range l.backends {
		b.Close()
	}
}

func (l *Logger) write(level int, format string, v ...interface{}) {
	l.Lock()
	defer l.Unlock()

	l.wg.Add(len(l.backends))

	for _, b := range l.backends {
		go func(b backend) {
			b.Write(level, l.context, format, v...)
			l.wg.Done()
		}(b)
	}

	l.wg.Wait()
}
