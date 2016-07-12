package logger

import "errors"

var (
	// ErrInvalidLevel represents an invalid logging level error.
	ErrInvalidLevel = errors.New("invalid logging level")
)
