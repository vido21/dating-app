package error

import (
	"fmt"
	"runtime"
)

// Color codes
const (
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorReset  = "\033[0m"
)

// AppError is a custom error type with a trace field.
type AppError struct {
	File          string
	Line          int
	OriginalError error
}

// Error implements the error interface for AppError.
func (e *AppError) Error() string {
	return fmt.Sprintf("%s%s:%d %s%s%s\n",
		ColorRed, e.File, e.Line, ColorReset,
		ColorYellow, e.OriginalError)
}

// NewAppError creates a new AppError with an original error and a trace.
func NewAppError(originalError error) *AppError {
	return &AppError{
		OriginalError: originalError,
	}
}

// AddTrace adds a trace to an existing error, returning a new AppError.
func AddTrace(err error) error {
	if err == nil {
		return nil
	}

	// Get the PC of the caller.
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		return err
	}

	// Get the function for the PC.
	caller := runtime.FuncForPC(pc)
	if caller == nil {
		return err
	}

	return &AppError{
		File:          file,
		Line:          line,
		OriginalError: err,
	}
}

func colorize(text, color string) string {
	return fmt.Sprintf("%s%s%s", color, text, ColorReset)
}
